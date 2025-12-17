package checks

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadProjectConfig(t *testing.T) {
	t.Run("Valid Config", func(t *testing.T) {
		content := `
name: test-project
services:
  api:
    language: js
    host: containerapp
    project: ./src/api
    docker:
      remote: true
hooks:
  postprovision:
    shell: sh
    run: ./scripts/postprovision.sh
  predeploy: echo "predeploy"
`
		tmpfile, err := os.CreateTemp("", "azure.yaml")
		require.NoError(t, err)
		defer os.Remove(tmpfile.Name())

		_, err = tmpfile.Write([]byte(content))
		require.NoError(t, err)
		tmpfile.Close()

		config, err := LoadProjectConfig(tmpfile.Name())
		require.NoError(t, err)

		assert.Equal(t, "test-project", config.Name)
		assert.Contains(t, config.Services, "api")
		assert.Equal(t, "js", config.Services["api"].Language)
		assert.Equal(t, "containerapp", config.Services["api"].Host)
		assert.True(t, config.Services["api"].Docker.Remote)

		// Check hooks
		assert.Contains(t, config.Hooks, "postprovision")
		assert.Equal(t, "sh", config.Hooks["postprovision"].Shell)
		assert.Equal(t, "./scripts/postprovision.sh", config.Hooks["postprovision"].Run)

		assert.Contains(t, config.Hooks, "predeploy")
		assert.Equal(t, "", config.Hooks["predeploy"].Shell)
		assert.Equal(t, "echo \"predeploy\"", config.Hooks["predeploy"].Run)
	})

	t.Run("Infra and RequiredVersions", func(t *testing.T) {
		content := `
name: test-infra
infra:
  provider: terraform
requiredVersions:
  azd: ">= 1.0.0"
  extensions:
    azure.ai.agents: ">= 0.1.0"
`
		tmpfile, err := os.CreateTemp("", "azure.yaml")
		require.NoError(t, err)
		defer os.Remove(tmpfile.Name())

		_, err = tmpfile.Write([]byte(content))
		require.NoError(t, err)
		tmpfile.Close()

		config, err := LoadProjectConfig(tmpfile.Name())
		require.NoError(t, err)

		assert.Equal(t, "test-infra", config.Name)
		assert.Equal(t, "terraform", config.Infra.Provider)
		assert.Equal(t, ">= 1.0.0", config.RequiredVersions.Azd)
		assert.Contains(t, config.RequiredVersions.Extensions, "azure.ai.agents")
		assert.Equal(t, ">= 0.1.0", config.RequiredVersions.Extensions["azure.ai.agents"])
	})

	t.Run("Invalid File", func(t *testing.T) {
		_, err := LoadProjectConfig("nonexistent.yaml")
		assert.Error(t, err)
	})

	t.Run("Invalid YAML", func(t *testing.T) {
		content := `
name: test-project
services:
  api:
    language: [invalid
`
		tmpfile, err := os.CreateTemp("", "azure.yaml")
		require.NoError(t, err)
		defer os.Remove(tmpfile.Name())

		_, err = tmpfile.Write([]byte(content))
		require.NoError(t, err)
		tmpfile.Close()

		_, err = LoadProjectConfig(tmpfile.Name())
		assert.Error(t, err)
	})
}
