// Package parser implements Dockerfile parsing for kubefy.
package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ParseDockerfile reads and parses a Dockerfile, returning the extracted configuration.
func ParseDockerfile(path string) (*DockerfileConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening dockerfile: %w", err)
	}
	defer file.Close()

	config := NewDockerfileConfig()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if err := parseLine(config, line); err != nil {
			return nil, err
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading dockerfile: %w", err)
	}

	if config.Image == "" {
		return nil, fmt.Errorf("no FROM directive found in Dockerfile")
	}

	return config, nil
}

func parseLine(config *DockerfileConfig, line string) error {
	parts := strings.SplitN(line, " ", 2)
	if len(parts) < 2 {
		return nil
	}

	directive := strings.ToUpper(parts[0])
	args := strings.TrimSpace(parts[1])

	switch directive {
	case "FROM":
		parseFrom(config, args)
	case "EXPOSE":
		parseExpose(config, args)
	case "ENV":
		parseEnv(config, args)
	case "CMD":
		parseCmd(config, args)
	case "ENTRYPOINT":
		parseEntrypoint(config, args)
	case "WORKDIR":
		config.WorkDir = args
	}

	return nil
}

func parseFrom(config *DockerfileConfig, args string) {
	// Handle "FROM image:tag AS builder" format
	parts := strings.Fields(args)
	if len(parts) == 0 {
		return
	}

	imageSpec := parts[0]

	// Handle scratch (no tag)
	if imageSpec == "scratch" {
		config.Image = "scratch"
		config.Tag = ""
		return
	}

	// Split image:tag
	if idx := strings.LastIndex(imageSpec, ":"); idx != -1 {
		config.Image = imageSpec[:idx]
		config.Tag = imageSpec[idx+1:]
	} else {
		config.Image = imageSpec
		config.Tag = "latest"
	}
}

func parseExpose(config *DockerfileConfig, args string) {
	for _, portSpec := range strings.Fields(args) {
		// Remove protocol suffix like /tcp or /udp
		portStr := strings.Split(portSpec, "/")[0]
		if port, err := strconv.Atoi(portStr); err == nil {
			config.Ports = append(config.Ports, port)
		}
	}
}

func parseEnv(config *DockerfileConfig, args string) {
	// Support KEY=value format
	if idx := strings.Index(args, "="); idx != -1 {
		key := args[:idx]
		value := args[idx+1:]
		// Remove quotes if present
		value = strings.Trim(value, `"'`)
		config.Env = append(config.Env, EnvPair{Key: key, Value: value})
	}
}

func parseCmd(config *DockerfileConfig, args string) {
	config.Command = parseCommand(args)
}

func parseEntrypoint(config *DockerfileConfig, args string) {
	config.Entrypoint = parseCommand(args)
}

func parseCommand(args string) []string {
	args = strings.TrimSpace(args)

	// JSON array format: ["cmd", "arg"]
	if strings.HasPrefix(args, "[") {
		var result []string
		if err := json.Unmarshal([]byte(args), &result); err == nil {
			return result
		}
	}

	// Shell format: cmd arg1 arg2
	return []string{"/bin/sh", "-c", args}
}
