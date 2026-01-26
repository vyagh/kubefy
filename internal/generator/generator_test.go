package generator

import (
	"strings"
	"testing"

	"github.com/vyagh/kubefy/internal/parser"
)

func TestCreateDeployment(t *testing.T) {
	config := &parser.DockerfileConfig{
		Image:   "node",
		Tag:     "18-alpine",
		Ports:   []int{3000},
		Env:     []parser.EnvPair{{Key: "NODE_ENV", Value: "production"}},
		Command: []string{"npm", "start"},
		WorkDir: "/app",
	}

	gen := New(config, Options{
		AppName:   "myapp",
		Namespace: "default",
		Replicas:  2,
	})

	yaml, err := gen.CreateDeployment()
	if err != nil {
		t.Fatalf("CreateDeployment failed: %v", err)
	}

	yamlStr := string(yaml)

	expected := []string{
		"apiVersion: apps/v1",
		"kind: Deployment",
		"name: myapp",
		"replicas: 2",
		"image: node:18-alpine",
		"containerPort: 3000",
	}

	for _, s := range expected {
		if !strings.Contains(yamlStr, s) {
			t.Errorf("expected %q in YAML", s)
		}
	}
}

func TestCreateService(t *testing.T) {
	config := &parser.DockerfileConfig{
		Image: "node",
		Tag:   "18",
		Ports: []int{3000},
	}

	gen := New(config, Options{
		AppName:     "myapp",
		ServiceType: "LoadBalancer",
	})

	yaml, err := gen.CreateService()
	if err != nil {
		t.Fatalf("CreateService failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "type: LoadBalancer") {
		t.Error("expected LoadBalancer service type")
	}
}

func TestCreateService_NoPorts(t *testing.T) {
	config := &parser.DockerfileConfig{
		Image: "ubuntu",
		Tag:   "22.04",
		Ports: []int{},
	}

	gen := New(config, Options{AppName: "test"})
	yaml, _ := gen.CreateService()

	if yaml != nil {
		t.Error("expected nil service when no ports")
	}
}
