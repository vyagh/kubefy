# kubefy

> Convert Dockerfiles to Kubernetes manifests

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Installation

```bash
git clone https://github.com/vyagh/kubefy
cd kubefy
make build
```

## Usage

```bash
# Basic usage
kubefy Dockerfile --name myapp

# With options
kubefy Dockerfile --name myapp --namespace production --replicas 3

# Preview without writing files
kubefy Dockerfile --name myapp --dry-run

# Different service type
kubefy Dockerfile --name myapp --service-type LoadBalancer
```

## Example

Given this Dockerfile:

```dockerfile
FROM node:18-alpine
WORKDIR /app
ENV NODE_ENV=production
EXPOSE 3000
CMD ["npm", "start"]
```

Running `kubefy Dockerfile --name myapp` generates:

**deployment.yaml**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
  namespace: default
  labels:
    app: myapp
spec:
  replicas: 1
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
        - name: myapp
          image: node:18-alpine
          ports:
            - containerPort: 3000
          env:
            - name: NODE_ENV
              value: production
          args:
            - npm
            - start
          workingDir: /app
```

**service.yaml**

```yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp
  namespace: default
spec:
  selector:
    app: myapp
  ports:
    - port: 3000
      targetPort: 3000
  type: ClusterIP
```

## Supported Directives

| Directive    | Support                   |
| ------------ | ------------------------- |
| `FROM`       | Image and tag extraction  |
| `EXPOSE`     | Single and multiple ports |
| `ENV`        | KEY=value format          |
| `CMD`        | JSON and shell formats    |
| `ENTRYPOINT` | JSON and shell formats    |
| `WORKDIR`    | Working directory         |

## Flags

| Flag             | Short | Default     | Description             |
| ---------------- | ----- | ----------- | ----------------------- |
| `--name`         | `-n`  | (required)  | Application name        |
| `--namespace`    |       | `default`   | Kubernetes namespace    |
| `--replicas`     | `-r`  | `1`         | Number of replicas      |
| `--output`       | `-o`  | `.`         | Output directory        |
| `--service-type` |       | `ClusterIP` | Service type            |
| `--dry-run`      |       | `false`     | Preview without writing |

## Limitations

- Multi-stage builds not supported
- ARG substitution not supported
- No Ingress generation

## License

MIT
