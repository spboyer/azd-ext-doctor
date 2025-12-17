package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewConfigureCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Configure the project",
		Long:  `Configure the project settings in azure.yaml.`,
	}

	cmd.AddCommand(newConfigureRemoteBuildCommand())

	return cmd
}

func newConfigureRemoteBuildCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "remote-build",
		Short: "Enable remote build for container services",
		Long:  `Updates azure.yaml to enable remote build (docker.remoteBuild: true) for all services using containerapp or aks host.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigureRemoteBuild()
		},
	}
}

func runConfigureRemoteBuild() error {
	projectFile := "azure.yaml"
	if _, err := os.Stat(projectFile); os.IsNotExist(err) {
		projectFile = "azure.yml"
		if _, err := os.Stat(projectFile); os.IsNotExist(err) {
			return fmt.Errorf("project file (azure.yaml/yml) not found")
		}
	}

	// Read file
	data, err := os.ReadFile(projectFile)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", projectFile, err)
	}

	// Parse into Node to preserve comments
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return fmt.Errorf("failed to parse %s: %w", projectFile, err)
	}

	// Traverse and update
	updated, err := enableRemoteBuildInNodes(&root)
	if err != nil {
		return err
	}

	if !updated {
		fmt.Println("No changes needed. Remote build is already enabled or no applicable services found.")
		return nil
	}

	// Write back
	// Use 2 spaces indentation
	encoder := yaml.NewEncoder(os.Stdout) // Temporary to check output? No, write to file.
	
	f, err := os.Create(projectFile)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %w", projectFile, err)
	}
	defer f.Close()

	encoder = yaml.NewEncoder(f)
	encoder.SetIndent(2)
	if err := encoder.Encode(&root); err != nil {
		return fmt.Errorf("failed to write %s: %w", projectFile, err)
	}

	fmt.Printf("Successfully enabled remote build in %s\n", projectFile)
	return nil
}

func enableRemoteBuildInNodes(root *yaml.Node) (bool, error) {
	if root.Kind != yaml.DocumentNode {
		return false, fmt.Errorf("expected document node")
	}
	if len(root.Content) == 0 {
		return false, nil
	}

	// Root content[0] is the map
	topMap := root.Content[0]
	if topMap.Kind != yaml.MappingNode {
		return false, fmt.Errorf("expected mapping node at root")
	}

	// Find "services" key
	var servicesNode *yaml.Node
	for i := 0; i < len(topMap.Content); i += 2 {
		key := topMap.Content[i]
		if key.Value == "services" {
			servicesNode = topMap.Content[i+1]
			break
		}
	}

	if servicesNode == nil {
		return false, nil
	}

	updated := false

	// Iterate over services
	for i := 0; i < len(servicesNode.Content); i += 2 {
		// serviceName := servicesNode.Content[i].Value
		serviceProps := servicesNode.Content[i+1]

		// Check host
		isContainer := false
		for j := 0; j < len(serviceProps.Content); j += 2 {
			k := serviceProps.Content[j]
			if k.Value == "host" {
				v := serviceProps.Content[j+1]
				if v.Value == "containerapp" || v.Value == "aks" {
					isContainer = true
					break
				}
			}
		}

		if isContainer {
			// Check/Add docker.remoteBuild
			if ensureRemoteBuild(serviceProps) {
				updated = true
			}
		}
	}

	return updated, nil
}

func ensureRemoteBuild(serviceProps *yaml.Node) bool {
	// Look for "docker" key
	var dockerNode *yaml.Node

	for i := 0; i < len(serviceProps.Content); i += 2 {
		if serviceProps.Content[i].Value == "docker" {
			dockerNode = serviceProps.Content[i+1]
			break
		}
	}

	if dockerNode == nil {
		// Create docker node
		dockerKey := &yaml.Node{Kind: yaml.ScalarNode, Value: "docker"}
		dockerVal := &yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "remoteBuild"},
			{Kind: yaml.ScalarNode, Value: "true", Tag: "!!bool"},
		}}
		serviceProps.Content = append(serviceProps.Content, dockerKey, dockerVal)
		return true
	}

	// Check remoteBuild in docker node
	for i := 0; i < len(dockerNode.Content); i += 2 {
		if dockerNode.Content[i].Value == "remoteBuild" {
			if dockerNode.Content[i+1].Value == "true" {
				return false // Already true
			}
			// Update to true
			dockerNode.Content[i+1].Value = "true"
			return true
		}
	}

	// Add remoteBuild to existing docker node
	dockerNode.Content = append(dockerNode.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "remoteBuild"},
		&yaml.Node{Kind: yaml.ScalarNode, Value: "true", Tag: "!!bool"},
	)
	return true
}
