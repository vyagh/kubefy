// Package cli implements the command-line interface for kubefy.
package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/vyagh/kubefy/internal/generator"
	"github.com/vyagh/kubefy/internal/parser"
)

var version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "kubefy [dockerfile]",
	Short: "Convert Dockerfiles to Kubernetes manifests",
	Long: `kubefy - Kubefy your Dockerfile

A CLI tool that parses Dockerfiles and generates Kubernetes
Deployment and Service YAML manifests.

Example:
  kubefy Dockerfile --name myapp
  kubefy Dockerfile --name myapp --namespace production --replicas 3`,
	Args: cobra.MaximumNArgs(1),
	RunE: runKubefy,
}

// CLI flags
var (
	appName     string
	namespace   string
	replicas    int
	outputDir   string
	serviceType string
	dryRun      bool
)

func init() {
	rootCmd.Flags().StringVarP(&appName, "name", "n", "", "Application name (required)")
	rootCmd.Flags().StringVar(&namespace, "namespace", "default", "Kubernetes namespace")
	rootCmd.Flags().IntVarP(&replicas, "replicas", "r", 1, "Number of replicas")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory")
	rootCmd.Flags().StringVar(&serviceType, "service-type", "ClusterIP", "Service type")
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview without writing files")
	rootCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("kubefy %s\n", version)
		},
	})
}

// Execute runs the CLI.
func Execute() error {
	return rootCmd.Execute()
}

func runKubefy(cmd *cobra.Command, args []string) error {
	dockerfilePath := "Dockerfile"
	if len(args) > 0 {
		dockerfilePath = args[0]
	}

	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		color.Red("✗ Dockerfile not found: %s", dockerfilePath)
		return fmt.Errorf("dockerfile not found: %s", dockerfilePath)
	}

	color.Cyan("⚡ Parsing %s...", dockerfilePath)

	config, err := parser.ParseDockerfile(dockerfilePath)
	if err != nil {
		color.Red("✗ Parse error: %v", err)
		return err
	}

	color.Green("✓ Parsed: %s:%s", config.Image, config.Tag)

	gen := generator.New(config, generator.Options{
		AppName:     appName,
		Namespace:   namespace,
		Replicas:    replicas,
		ServiceType: serviceType,
	})

	deploymentYAML, err := gen.CreateDeployment()
	if err != nil {
		return err
	}

	serviceYAML, _ := gen.CreateService()

	if dryRun {
		fmt.Println("---")
		fmt.Print(string(deploymentYAML))
		if serviceYAML != nil {
			fmt.Println("---")
			fmt.Print(string(serviceYAML))
		}
		color.Green("\n✓ Manifests generated (dry-run)")
		return nil
	}

	// Write files
	os.MkdirAll(outputDir, 0755)

	deployPath := filepath.Join(outputDir, "deployment.yaml")
	os.WriteFile(deployPath, deploymentYAML, 0644)
	color.Green("✓ Created %s", deployPath)

	if serviceYAML != nil {
		servicePath := filepath.Join(outputDir, "service.yaml")
		os.WriteFile(servicePath, serviceYAML, 0644)
		color.Green("✓ Created %s", servicePath)
	}

	color.Cyan("\n📦 Apply: kubectl apply -f %s/", outputDir)
	return nil
}
