// Package generator creates Kubernetes manifests from parsed Dockerfile config.
package generator

import (
	"fmt"

	"github.com/vyagh/kubefy/internal/parser"
	"gopkg.in/yaml.v3"
)

// Options configures the manifest generation.
type Options struct {
	AppName     string
	Namespace   string
	Replicas    int
	ServiceType string
}

// Generator creates Kubernetes manifests from Dockerfile configuration.
type Generator struct {
	config  *parser.DockerfileConfig
	options Options
}

// New creates a new Generator with the given config and options.
func New(config *parser.DockerfileConfig, opts Options) *Generator {
	return &Generator{
		config:  config,
		options: opts,
	}
}

// CreateDeployment generates a Kubernetes Deployment manifest.
func (g *Generator) CreateDeployment() ([]byte, error) {
	labels := map[string]string{
		"app": g.options.AppName,
	}

	// Handle images without tags (like scratch)
	imageRef := g.config.Image
	if g.config.Tag != "" {
		imageRef = fmt.Sprintf("%s:%s", g.config.Image, g.config.Tag)
	}

	container := Container{
		Name:  g.options.AppName,
		Image: imageRef,
	}

	for _, port := range g.config.Ports {
		container.Ports = append(container.Ports, PortMapping{
			ContainerPort: port,
		})
	}

	for _, envPair := range g.config.Env {
		container.Env = append(container.Env, EnvVar{
			Name:  envPair.Key,
			Value: envPair.Value,
		})
	}

	if len(g.config.Entrypoint) > 0 {
		container.Command = g.config.Entrypoint
	}
	if len(g.config.Command) > 0 {
		container.Args = g.config.Command
	}

	if g.config.WorkDir != "" {
		container.WorkingDir = g.config.WorkDir
	}

	deployment := Deployment{
		APIVersion: "apps/v1",
		Kind:       "Deployment",
		Metadata: Metadata{
			Name:      g.options.AppName,
			Namespace: g.options.Namespace,
			Labels:    labels,
		},
		Spec: DeploymentSpec{
			Replicas: g.options.Replicas,
			Selector: LabelSelector{
				MatchLabels: labels,
			},
			Template: PodTemplateSpec{
				Metadata: PodMetadata{
					Labels: labels,
				},
				Spec: PodSpec{
					Containers: []Container{container},
				},
			},
		},
	}

	return yaml.Marshal(deployment)
}

// CreateService generates a Kubernetes Service manifest.
func (g *Generator) CreateService() ([]byte, error) {
	if len(g.config.Ports) == 0 {
		return nil, nil // No service needed if no ports exposed
	}

	labels := map[string]string{
		"app": g.options.AppName,
	}

	var servicePorts []ServicePort
	for i, port := range g.config.Ports {
		sp := ServicePort{
			Port:       port,
			TargetPort: port,
		}
		// Name ports if there are multiple
		if len(g.config.Ports) > 1 {
			sp.Name = fmt.Sprintf("port-%d", i+1)
		}
		servicePorts = append(servicePorts, sp)
	}

	service := Service{
		APIVersion: "v1",
		Kind:       "Service",
		Metadata: Metadata{
			Name:      g.options.AppName,
			Namespace: g.options.Namespace,
			Labels:    labels,
		},
		Spec: ServiceSpec{
			Selector: labels,
			Ports:    servicePorts,
			Type:     g.options.ServiceType,
		},
	}

	return yaml.Marshal(service)
}
