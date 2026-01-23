// Package parser implements Dockerfile parsing for kubefy.
package parser

// EnvPair represents a key-value environment variable.
type EnvPair struct {
	Key   string
	Value string
}

// DockerfileConfig holds the parsed configuration from a Dockerfile.
type DockerfileConfig struct {
	Image      string    // Base image name from FROM
	Tag        string    // Image tag from FROM (defaults to "latest")
	Ports      []int     // Exposed ports from EXPOSE
	Env        []EnvPair // Environment variables from ENV (order preserved)
	Command    []string  // Command from CMD
	Entrypoint []string  // Entrypoint from ENTRYPOINT
	WorkDir    string    // Working directory from WORKDIR
}

// NewDockerfileConfig creates a DockerfileConfig with default values.
func NewDockerfileConfig() *DockerfileConfig {
	return &DockerfileConfig{
		Tag:   "latest",
		Ports: make([]int, 0),
		Env:   make([]EnvPair, 0),
	}
}
