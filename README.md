# AZD Doctor Extension

An `azd` extension that checks for necessary prerequisites based on the current project's configuration.

## Screenshot

![Example output for `azd doctor check`](docs/images/doctor-check.png)

## Installation

### Enable Extensions & Install azd-doctor

```bash
# Add azd-doctor extension source
azd extension source add -n doc -t url -l https://raw.githubusercontent.com/spboyer/azd-ext-doctor/main/registry.json

# Install the extension
azd extension install spboyer.azd.doctor
```

### Upgrade

To upgrade the extension to the latest version:

```bash
azd extension upgrade spboyer.azd.doctor
```

> **Note:** If the upgrade command does not detect the latest version due to caching, use the install command with the `--force` flag:
> ```bash
> azd extension install spboyer.azd.doctor --force
> ```

To upgrade all extensions:

```bash
azd extension upgrade --all
```

### Install (Local Development)

Recommended workflow is to use the `azd` developer extension (`azd x`):

```bash
# If needed, add the dev extension source (one-time)
azd extension source add -n dev -t url -l "https://aka.ms/azd/extensions/registry/dev"

# Install the developer extension (one-time)
azd extension install microsoft.azd.extensions

# From this repo
azd x build
```

## Features

Checks for the presence and version of the following tools:

- **Container Runtimes**: Docker or Podman (checks if daemon is running)
- **Language Runtimes**:
  - Node.js
  - Python
  - .NET SDK
- **Shells**:
  - Bash
  - PowerShell (pwsh/powershell)
- **Azure Tools**:
  - Azure Functions Core Tools
  - Azure Static Web Apps CLI (swa)
- **Infrastructure as Code**:
  - Terraform
- **Extensions**:
  - Validates required extension versions specified in `azure.yaml`

## Commands

### `check`

Runs all the prerequisite checks.

```bash
azd doctor check
```

Flags:

- Skip auth check (avoids any auth-related delay):
  ```bash
  azd doctor check --skip-auth
  ```
- Limit auth check time (default is `5s`):
  ```bash
  azd doctor check --auth-timeout 2s
  ```

### `verify`

Verifies that the environment meets the requirements specified in `azure.yaml`. This command is designed to be used in CI/CD pipelines or as a pre-deployment check.

It checks:
- Required tools based on the project configuration (e.g., `swa` if using Static Web Apps, `terraform` if using Terraform).
- Required extension versions specified in `requiredVersions.extensions`.
- Suggests enabling `remoteBuild` if Docker is missing for container apps.

```bash
azd doctor verify
```

### `configure`

Helps configure project settings in `azure.yaml`.

#### `remote-build`

Enables remote build (`docker.remoteBuild: true`) for all services using `containerapp` or `aks` host. This is useful when local Docker is not available.

```bash
azd doctor configure remote-build
```

### `context`

Displays the context of the current AZD project and environment.

```bash
azd doctor context
```

## Lifecycle Hooks

The extension automatically registers a `predeploy` hook to run `azd doctor verify` before deployment. This ensures that the environment is correctly set up before attempting to deploy.

To enable this, ensure the extension is installed and enabled in your project.

### Bypassing Verification

You can bypass the verification check by setting the `AZD_DOCTOR_SKIP_VERIFY` environment variable.

- **Skip all checks**:
  ```bash
  export AZD_DOCTOR_SKIP_VERIFY=true
  # or
  export AZD_DOCTOR_SKIP_VERIFY=all
`  ```

- **Skip for specific commands**:
  ```bash
  # Skip only for deploy
  export AZD_DOCTOR_SKIP_VERIFY=deploy

  # Skip for provision and deploy
  export AZD_DOCTOR_SKIP_VERIFY=provision,deploy
  ```

Supported values for specific commands are `up`, `provision`, and `deploy`.


## Development

### Prerequisites

- Go 1.24+
- Azure Developer CLI (`azd`)

### Build

```bash
go build -o azd-ext-doctor .
```

### Install (Local Development)

For local development, use:

```bash
azd x watch
```

