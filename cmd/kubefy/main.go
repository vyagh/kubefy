// Package main is the entry point for the kubefy CLI.
package main

import (
	"os"

	"github.com/vyagh/kubefy/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
