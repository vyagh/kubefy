package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseDockerfile_Simple(t *testing.T) {
	content := `FROM node:18-alpine
WORKDIR /app
ENV NODE_ENV=production
EXPOSE 3000
CMD ["npm", "start"]`

	config := parseTestDockerfile(t, content)

	assertEqual(t, "Image", config.Image, "node")
	assertEqual(t, "Tag", config.Tag, "18-alpine")
	assertEqual(t, "WorkDir", config.WorkDir, "/app")

	if len(config.Ports) != 1 || config.Ports[0] != 3000 {
		t.Errorf("expected Ports [3000], got %v", config.Ports)
	}

	if len(config.Command) != 2 || config.Command[0] != "npm" {
		t.Errorf("expected Command [npm, start], got %v", config.Command)
	}
}

func TestParseDockerfile_MultiplePorts(t *testing.T) {
	content := `FROM nginx:alpine
EXPOSE 80 443`

	config := parseTestDockerfile(t, content)

	if len(config.Ports) != 2 || config.Ports[0] != 80 || config.Ports[1] != 443 {
		t.Errorf("expected ports [80, 443], got %v", config.Ports)
	}
}

func TestParseDockerfile_Entrypoint(t *testing.T) {
	content := `FROM python:3.11-slim
ENTRYPOINT ["python"]
CMD ["app.py"]`

	config := parseTestDockerfile(t, content)

	if len(config.Entrypoint) != 1 || config.Entrypoint[0] != "python" {
		t.Errorf("expected Entrypoint [python], got %v", config.Entrypoint)
	}
}

func TestParseDockerfile_NoFrom(t *testing.T) {
	content := `EXPOSE 3000
CMD ["npm", "start"]`

	tmpFile := createTempDockerfile(t, content)
	_, err := ParseDockerfile(tmpFile)

	if err == nil {
		t.Error("expected error for Dockerfile without FROM")
	}
}

func TestParseDockerfile_FromVariants(t *testing.T) {
	tests := []struct {
		content string
		image   string
		tag     string
	}{
		{"FROM ubuntu", "ubuntu", "latest"},
		{"FROM ubuntu:22.04", "ubuntu", "22.04"},
		{"FROM scratch", "scratch", ""},
	}

	for _, tc := range tests {
		config := parseTestDockerfile(t, tc.content)
		assertEqual(t, "Image", config.Image, tc.image)
		assertEqual(t, "Tag", config.Tag, tc.tag)
	}
}

// Helper functions

func parseTestDockerfile(t *testing.T, content string) *DockerfileConfig {
	t.Helper()
	tmpFile := createTempDockerfile(t, content)
	config, err := ParseDockerfile(tmpFile)
	if err != nil {
		t.Fatalf("ParseDockerfile failed: %v", err)
	}
	return config
}

func createTempDockerfile(t *testing.T, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "Dockerfile")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create temp Dockerfile: %v", err)
	}
	return tmpFile
}

func assertEqual(t *testing.T, name, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %q, want %q", name, got, want)
	}
}
