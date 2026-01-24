// Package generator creates Kubernetes manifests from parsed Dockerfile config.
package generator

// Deployment represents a Kubernetes Deployment resource.
type Deployment struct {
	APIVersion string         `yaml:"apiVersion"`
	Kind       string         `yaml:"kind"`
	Metadata   Metadata       `yaml:"metadata"`
	Spec       DeploymentSpec `yaml:"spec"`
}

// DeploymentSpec defines the desired state of a Deployment.
type DeploymentSpec struct {
	Replicas int             `yaml:"replicas"`
	Selector LabelSelector   `yaml:"selector"`
	Template PodTemplateSpec `yaml:"template"`
}

// LabelSelector matches pods by labels.
type LabelSelector struct {
	MatchLabels map[string]string `yaml:"matchLabels"`
}

// PodTemplateSpec describes the pods that will be created.
type PodTemplateSpec struct {
	Metadata PodMetadata `yaml:"metadata"`
	Spec     PodSpec     `yaml:"spec"`
}

// PodMetadata contains metadata for pod templates (no name field).
type PodMetadata struct {
	Labels map[string]string `yaml:"labels,omitempty"`
}

// PodSpec defines the containers and volumes in a pod.
type PodSpec struct {
	Containers []Container `yaml:"containers"`
}

// Container represents a single container within a pod.
type Container struct {
	Name       string        `yaml:"name"`
	Image      string        `yaml:"image"`
	Ports      []PortMapping `yaml:"ports,omitempty"`
	Env        []EnvVar      `yaml:"env,omitempty"`
	Command    []string      `yaml:"command,omitempty"`
	Args       []string      `yaml:"args,omitempty"`
	WorkingDir string        `yaml:"workingDir,omitempty"`
}

// PortMapping defines a container port.
type PortMapping struct {
	ContainerPort int `yaml:"containerPort"`
}

// EnvVar represents an environment variable.
type EnvVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// Service represents a Kubernetes Service resource.
type Service struct {
	APIVersion string      `yaml:"apiVersion"`
	Kind       string      `yaml:"kind"`
	Metadata   Metadata    `yaml:"metadata"`
	Spec       ServiceSpec `yaml:"spec"`
}

// ServiceSpec defines the desired state of a Service.
type ServiceSpec struct {
	Selector map[string]string `yaml:"selector"`
	Ports    []ServicePort     `yaml:"ports"`
	Type     string            `yaml:"type"`
}

// ServicePort defines a port exposed by a Service.
type ServicePort struct {
	Port       int    `yaml:"port"`
	TargetPort int    `yaml:"targetPort"`
	Name       string `yaml:"name,omitempty"`
}

// Metadata contains common Kubernetes object metadata.
type Metadata struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace,omitempty"`
	Labels    map[string]string `yaml:"labels,omitempty"`
}
