package checks

import (
	"context"
	"fmt"
	"testing"

	"github.com/azure/azure-dev/cli/azd/pkg/azdext"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

// MockRunner implements Runner for testing
type MockRunner struct {
	OutputFunc func(name string, args ...string) ([]byte, error)
	RunFunc    func(name string, args ...string) error
}

func (m *MockRunner) Output(name string, args ...string) ([]byte, error) {
	if m.OutputFunc != nil {
		return m.OutputFunc(name, args...)
	}
	return nil, nil
}

func (m *MockRunner) Run(name string, args ...string) error {
	if m.RunFunc != nil {
		return m.RunFunc(name, args...)
	}
	return nil
}

// MockAzdClient implements IAzdClient
type MockAzdClient struct {
	DeploymentClient azdext.DeploymentServiceClient
	ProjectClient    azdext.ProjectServiceClient
}

func (m *MockAzdClient) Deployment() azdext.DeploymentServiceClient {
	return m.DeploymentClient
}

func (m *MockAzdClient) Project() azdext.ProjectServiceClient {
	return m.ProjectClient
}

// MockDeploymentService implements azdext.DeploymentServiceClient
type MockDeploymentService struct {
	GetDeploymentFunc        func(ctx context.Context, in *azdext.EmptyRequest, opts ...grpc.CallOption) (*azdext.GetDeploymentResponse, error)
	GetDeploymentContextFunc func(ctx context.Context, in *azdext.EmptyRequest, opts ...grpc.CallOption) (*azdext.GetDeploymentContextResponse, error)
}

func (m *MockDeploymentService) GetDeployment(ctx context.Context, in *azdext.EmptyRequest, opts ...grpc.CallOption) (*azdext.GetDeploymentResponse, error) {
	if m.GetDeploymentFunc != nil {
		return m.GetDeploymentFunc(ctx, in, opts...)
	}
	return nil, nil
}

func (m *MockDeploymentService) GetDeploymentContext(ctx context.Context, in *azdext.EmptyRequest, opts ...grpc.CallOption) (*azdext.GetDeploymentContextResponse, error) {
	if m.GetDeploymentContextFunc != nil {
		return m.GetDeploymentContextFunc(ctx, in, opts...)
	}
	return nil, nil
}

// MockProjectService implements azdext.ProjectServiceClient
type MockProjectService struct {
	GetFunc                     func(ctx context.Context, in *azdext.EmptyRequest, opts ...grpc.CallOption) (*azdext.GetProjectResponse, error)
	AddServiceFunc              func(ctx context.Context, in *azdext.AddServiceRequest, opts ...grpc.CallOption) (*azdext.EmptyResponse, error)
	GetConfigSectionFunc        func(ctx context.Context, in *azdext.GetProjectConfigSectionRequest, opts ...grpc.CallOption) (*azdext.GetProjectConfigSectionResponse, error)
	GetConfigValueFunc          func(ctx context.Context, in *azdext.GetProjectConfigValueRequest, opts ...grpc.CallOption) (*azdext.GetProjectConfigValueResponse, error)
	GetResolvedServicesFunc     func(ctx context.Context, in *azdext.EmptyRequest, opts ...grpc.CallOption) (*azdext.GetResolvedServicesResponse, error)
	GetServiceConfigSectionFunc func(ctx context.Context, in *azdext.GetServiceConfigSectionRequest, opts ...grpc.CallOption) (*azdext.GetServiceConfigSectionResponse, error)
	GetServiceConfigValueFunc   func(ctx context.Context, in *azdext.GetServiceConfigValueRequest, opts ...grpc.CallOption) (*azdext.GetServiceConfigValueResponse, error)
	ParseGitHubUrlFunc          func(ctx context.Context, in *azdext.ParseGitHubUrlRequest, opts ...grpc.CallOption) (*azdext.ParseGitHubUrlResponse, error)
	SetConfigSectionFunc        func(ctx context.Context, in *azdext.SetProjectConfigSectionRequest, opts ...grpc.CallOption) (*azdext.EmptyResponse, error)
	SetConfigValueFunc          func(ctx context.Context, in *azdext.SetProjectConfigValueRequest, opts ...grpc.CallOption) (*azdext.EmptyResponse, error)
	SetServiceConfigSectionFunc func(ctx context.Context, in *azdext.SetServiceConfigSectionRequest, opts ...grpc.CallOption) (*azdext.EmptyResponse, error)
	SetServiceConfigValueFunc   func(ctx context.Context, in *azdext.SetServiceConfigValueRequest, opts ...grpc.CallOption) (*azdext.EmptyResponse, error)
	UnsetConfigFunc             func(ctx context.Context, in *azdext.UnsetProjectConfigRequest, opts ...grpc.CallOption) (*azdext.EmptyResponse, error)
	UnsetServiceConfigFunc      func(ctx context.Context, in *azdext.UnsetServiceConfigRequest, opts ...grpc.CallOption) (*azdext.EmptyResponse, error)
}

func (m *MockProjectService) Get(ctx context.Context, in *azdext.EmptyRequest, opts ...grpc.CallOption) (*azdext.GetProjectResponse, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, in, opts...)
	}
	return nil, nil
}

func (m *MockProjectService) AddService(ctx context.Context, in *azdext.AddServiceRequest, opts ...grpc.CallOption) (*azdext.EmptyResponse, error) {
	if m.AddServiceFunc != nil {
		return m.AddServiceFunc(ctx, in, opts...)
	}
	return nil, nil
}

func (m *MockProjectService) GetConfigSection(ctx context.Context, in *azdext.GetProjectConfigSectionRequest, opts ...grpc.CallOption) (*azdext.GetProjectConfigSectionResponse, error) {
	if m.GetConfigSectionFunc != nil {
		return m.GetConfigSectionFunc(ctx, in, opts...)
	}
	return nil, nil
}

func (m *MockProjectService) GetConfigValue(ctx context.Context, in *azdext.GetProjectConfigValueRequest, opts ...grpc.CallOption) (*azdext.GetProjectConfigValueResponse, error) {
	if m.GetConfigValueFunc != nil {
		return m.GetConfigValueFunc(ctx, in, opts...)
	}
	return nil, nil
}

func (m *MockProjectService) GetResolvedServices(ctx context.Context, in *azdext.EmptyRequest, opts ...grpc.CallOption) (*azdext.GetResolvedServicesResponse, error) {
	if m.GetResolvedServicesFunc != nil {
		return m.GetResolvedServicesFunc(ctx, in, opts...)
	}
	return nil, nil
}

func (m *MockProjectService) GetServiceConfigSection(ctx context.Context, in *azdext.GetServiceConfigSectionRequest, opts ...grpc.CallOption) (*azdext.GetServiceConfigSectionResponse, error) {
	if m.GetServiceConfigSectionFunc != nil {
		return m.GetServiceConfigSectionFunc(ctx, in, opts...)
	}
	return nil, nil
}

func (m *MockProjectService) GetServiceConfigValue(ctx context.Context, in *azdext.GetServiceConfigValueRequest, opts ...grpc.CallOption) (*azdext.GetServiceConfigValueResponse, error) {
	if m.GetServiceConfigValueFunc != nil {
		return m.GetServiceConfigValueFunc(ctx, in, opts...)
	}
	return nil, nil
}

func (m *MockProjectService) ParseGitHubUrl(ctx context.Context, in *azdext.ParseGitHubUrlRequest, opts ...grpc.CallOption) (*azdext.ParseGitHubUrlResponse, error) {
	if m.ParseGitHubUrlFunc != nil {
		return m.ParseGitHubUrlFunc(ctx, in, opts...)
	}
	return nil, nil
}

func (m *MockProjectService) SetConfigSection(ctx context.Context, in *azdext.SetProjectConfigSectionRequest, opts ...grpc.CallOption) (*azdext.EmptyResponse, error) {
	if m.SetConfigSectionFunc != nil {
		return m.SetConfigSectionFunc(ctx, in, opts...)
	}
	return nil, nil
}

func (m *MockProjectService) SetConfigValue(ctx context.Context, in *azdext.SetProjectConfigValueRequest, opts ...grpc.CallOption) (*azdext.EmptyResponse, error) {
	if m.SetConfigValueFunc != nil {
		return m.SetConfigValueFunc(ctx, in, opts...)
	}
	return nil, nil
}

func (m *MockProjectService) SetServiceConfigSection(ctx context.Context, in *azdext.SetServiceConfigSectionRequest, opts ...grpc.CallOption) (*azdext.EmptyResponse, error) {
	if m.SetServiceConfigSectionFunc != nil {
		return m.SetServiceConfigSectionFunc(ctx, in, opts...)
	}
	return nil, nil
}

func (m *MockProjectService) SetServiceConfigValue(ctx context.Context, in *azdext.SetServiceConfigValueRequest, opts ...grpc.CallOption) (*azdext.EmptyResponse, error) {
	if m.SetServiceConfigValueFunc != nil {
		return m.SetServiceConfigValueFunc(ctx, in, opts...)
	}
	return nil, nil
}

func (m *MockProjectService) UnsetConfig(ctx context.Context, in *azdext.UnsetProjectConfigRequest, opts ...grpc.CallOption) (*azdext.EmptyResponse, error) {
	if m.UnsetConfigFunc != nil {
		return m.UnsetConfigFunc(ctx, in, opts...)
	}
	return nil, nil
}

func (m *MockProjectService) UnsetServiceConfig(ctx context.Context, in *azdext.UnsetServiceConfigRequest, opts ...grpc.CallOption) (*azdext.EmptyResponse, error) {
	if m.UnsetServiceConfigFunc != nil {
		return m.UnsetServiceConfigFunc(ctx, in, opts...)
	}
	return nil, nil
}

func TestCheckTool(t *testing.T) {
	// Save original runner and restore after test
	origRunner := CommandRunner
	defer func() { CommandRunner = origRunner }()

	t.Run("Tool Installed", func(t *testing.T) {
		CommandRunner = &MockRunner{
			OutputFunc: func(name string, args ...string) ([]byte, error) {
				return []byte("1.0.0\n"), nil
			},
		}

		res := CheckTool("test-tool", "--version")
		assert.True(t, res.Installed)
		assert.Equal(t, "1.0.0", res.Version)
		assert.True(t, res.Running)
		assert.NoError(t, res.Error)
	})

	t.Run("Tool Not Installed", func(t *testing.T) {
		CommandRunner = &MockRunner{
			OutputFunc: func(name string, args ...string) ([]byte, error) {
				return nil, fmt.Errorf("executable file not found in $PATH")
			},
		}

		res := CheckTool("test-tool", "--version")
		assert.False(t, res.Installed)
		assert.Error(t, res.Error)
	})
}

func TestCheckAzdLogin(t *testing.T) {
	ctx := context.Background()
	// Save original runner and restore after test
	origRunner := CommandRunner
	defer func() { CommandRunner = origRunner }()

	t.Run("Logged In", func(t *testing.T) {
		CommandRunner = &MockRunner{
			OutputFunc: func(name string, args ...string) ([]byte, error) {
				if name == "azd" && args[0] == "auth" && args[1] == "login" && args[2] == "--check-status" {
					return []byte("Logged in to Azure as user@example.com\n"), nil
				}
				return nil, fmt.Errorf("unexpected command: %s %v", name, args)
			},
		}

		// We don't need the client for this check anymore, so passing nil or a mock is fine
		// The function signature still requires it, so we pass a mock
		mockClient := &MockAzdClient{}

		res := CheckAzdLogin(ctx, mockClient)
		assert.True(t, res.Installed)
		assert.True(t, res.Running)
		assert.Contains(t, res.Version, "Logged in to Azure as user@example.com")
	})

	t.Run("Not Logged In (Error)", func(t *testing.T) {
		CommandRunner = &MockRunner{
			OutputFunc: func(name string, args ...string) ([]byte, error) {
				if name == "azd" && args[0] == "auth" && args[1] == "login" && args[2] == "--check-status" {
					return []byte("Not logged in, run `azd auth login` to login to Azure\n"), fmt.Errorf("exit status 1")
				}
				return nil, fmt.Errorf("unexpected command: %s %v", name, args)
			},
		}

		mockClient := &MockAzdClient{}

		res := CheckAzdLogin(ctx, mockClient)
		assert.True(t, res.Installed)
		assert.False(t, res.Running)
		assert.Contains(t, res.Version, "Not logged in")
		assert.Error(t, res.Error)
	})
}

func TestCheckAzdInit(t *testing.T) {
	ctx := context.Background()

	t.Run("Initialized", func(t *testing.T) {
		mockClient := &MockAzdClient{
			ProjectClient: &MockProjectService{
				GetFunc: func(ctx context.Context, in *azdext.EmptyRequest, opts ...grpc.CallOption) (*azdext.GetProjectResponse, error) {
					return &azdext.GetProjectResponse{
						Project: &azdext.ProjectConfig{
							Name: "test-project",
						},
					}, nil
				},
			},
		}

		res := CheckAzdInit(ctx, mockClient)
		assert.True(t, res.Installed)
		assert.True(t, res.Running)
		assert.Contains(t, res.Version, "Initialized")
		assert.Contains(t, res.Version, "test-project")
	})

	t.Run("Not Initialized (Error)", func(t *testing.T) {
		mockClient := &MockAzdClient{
			ProjectClient: &MockProjectService{
				GetFunc: func(ctx context.Context, in *azdext.EmptyRequest, opts ...grpc.CallOption) (*azdext.GetProjectResponse, error) {
					return nil, fmt.Errorf("no project found")
				},
			},
		}

		res := CheckAzdInit(ctx, mockClient)
		assert.False(t, res.Installed)
		assert.Error(t, res.Error)
	})
}

func TestCheckSwaCli(t *testing.T) {
	origRunner := CommandRunner
	defer func() { CommandRunner = origRunner }()

	t.Run("Installed", func(t *testing.T) {
		CommandRunner = &MockRunner{
			OutputFunc: func(name string, args ...string) ([]byte, error) {
				if name == "swa" && args[0] == "--version" {
					return []byte("1.1.0\n"), nil
				}
				return nil, fmt.Errorf("unexpected command")
			},
		}
		res := CheckSwaCli()
		assert.True(t, res.Installed)
		assert.Equal(t, "1.1.0", res.Version)
	})

	t.Run("Not Installed", func(t *testing.T) {
		CommandRunner = &MockRunner{
			OutputFunc: func(name string, args ...string) ([]byte, error) {
				return nil, fmt.Errorf("not found")
			},
		}
		res := CheckSwaCli()
		assert.False(t, res.Installed)
	})
}

func TestCheckTerraform(t *testing.T) {
	origRunner := CommandRunner
	defer func() { CommandRunner = origRunner }()

	t.Run("Installed", func(t *testing.T) {
		CommandRunner = &MockRunner{
			OutputFunc: func(name string, args ...string) ([]byte, error) {
				if name == "terraform" && args[0] == "--version" {
					return []byte("Terraform v1.5.0\n"), nil
				}
				return nil, fmt.Errorf("unexpected command")
			},
		}
		res := CheckTerraform()
		assert.True(t, res.Installed)
		assert.Contains(t, res.Version, "v1.5.0")
	})
}

func TestCheckExtension(t *testing.T) {
	origRunner := CommandRunner
	defer func() { CommandRunner = origRunner }()

	t.Run("Extension Found and Valid", func(t *testing.T) {
		CommandRunner = &MockRunner{
			OutputFunc: func(name string, args ...string) ([]byte, error) {
				if name == "azd" && args[0] == "extension" && args[1] == "list" {
					return []byte(`[{"name": "test-ext", "version": "1.0.0"}]`), nil
				}
				return nil, fmt.Errorf("unexpected command")
			},
		}
		res := CheckExtension("test-ext", ">= 1.0.0")
		assert.True(t, res.Installed)
		assert.True(t, res.Running)
		assert.Equal(t, "1.0.0", res.Version)
	})

	t.Run("Extension Found but Invalid Version", func(t *testing.T) {
		CommandRunner = &MockRunner{
			OutputFunc: func(name string, args ...string) ([]byte, error) {
				if name == "azd" && args[0] == "extension" && args[1] == "list" {
					return []byte(`[{"name": "test-ext", "version": "0.5.0"}]`), nil
				}
				return nil, fmt.Errorf("unexpected command")
			},
		}
		res := CheckExtension("test-ext", ">= 1.0.0")
		assert.True(t, res.Installed)
		assert.False(t, res.Running)
		assert.Equal(t, "0.5.0", res.Version)
		assert.Error(t, res.Error)
	})

	t.Run("Extension Not Found", func(t *testing.T) {
		CommandRunner = &MockRunner{
			OutputFunc: func(name string, args ...string) ([]byte, error) {
				if name == "azd" && args[0] == "extension" && args[1] == "list" {
					return []byte(`[]`), nil
				}
				return nil, fmt.Errorf("unexpected command")
			},
		}
		res := CheckExtension("test-ext", ">= 1.0.0")
		assert.False(t, res.Installed)
		assert.Error(t, res.Error)
	})
}
