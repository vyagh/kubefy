// Package generator creates Kubernetes manifests from parsed Dockerfile config.
package generator

import (
	"fmt"

	"github.com/vyagh/kubefy/internal/parser"
	"gopkg.in/yaml.v3"
)

// Options holds configuration for the generator.
type Options struct {
	AppName     string
	Namespace   string
	Replicas    int
	ServiceType string
}

// Generator holds the state for generating manifests.
type Generator struct {
	config  *parser.DockerfileConfig
	options Options
}

// New creates a new Generator.
func New(config *parser.DockerfileConfig, options Options) *Generator {
	return &Generator{
		config:  config,
		options: options,
	}
}

// CreateDeployment generates the Deployment YAML.
func (g *Generator) CreateDeployment() ([]byte, error) {
	matchLabels := map[string]string{"app": g.options.AppName}

	var ports []PortMapping
	for _, p := range g.config.Ports {
		ports = append(ports, PortMapping{ContainerPort: p})
	}

	var envs []EnvVar
	for _, e := range g.config.Env {
		envs = append(envs, EnvVar{Name: e.Key, Value: e.Value})
	}

	deployment := Deployment{
		APIVersion: "apps/v1",
		Kind:       "Deployment",
		Metadata: Metadata{
			Name:      g.options.AppName,
			Namespace: g.options.Namespace,
			Labels:    matchLabels,
		},
		Spec: DeploymentSpec{
			Replicas: g.options.Replicas,
			Selector: LabelSelector{MatchLabels: matchLabels},
			Template: PodTemplateSpec{
				Metadata: PodMetadata{Labels: matchLabels},
				Spec: PodSpec{
					Containers: []Container{
						{
							Name:       g.options.AppName,
							Image:      fmt.Sprintf("%s:%s", g.config.Image, g.config.Tag),
							Ports:      ports,
							Env:        envs,
							Command:    g.config.Entrypoint,
							Args:       g.config.Command,
							WorkingDir: g.config.WorkDir,
						},
					},
				},
			},
		},
	}

	return yaml.Marshal(deployment)
}

// CreateService generates the Service YAML.
func (g *Generator) CreateService() ([]byte, error) {
	if len(g.config.Ports) == 0 {
		return nil, nil
	}

	var ports []ServicePort
	for _, p := range g.config.Ports {
		ports = append(ports, ServicePort{
			Port:       p,
			TargetPort: p,
		})
	}

	service := Service{
		APIVersion: "v1",
		Kind:       "Service",
		Metadata: Metadata{
			Name:      g.options.AppName,
			Namespace: g.options.Namespace,
		},
		Spec: ServiceSpec{
			Selector: map[string]string{"app": g.options.AppName},
			Ports:    ports,
			Type:     g.options.ServiceType,
		},
	}

	return yaml.Marshal(service)
}
