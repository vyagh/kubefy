.PHONY: build test clean install lint

VERSION ?= 0.1.0
LDFLAGS := -ldflags "-X github.com/vyagh/kubefy/internal/cli.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o kubefy ./cmd/kubefy

test:
	go test -v ./...

lint:
	go vet ./...
	go fmt ./...

clean:
	rm -f kubefy

install: build
	cp kubefy $(GOPATH)/bin/kubefy
