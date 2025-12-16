# AZD Doctor Extension

An `azd` extension that checks for necessary prerequisites based on the current project's configuration.

## Screenshot

![Example output for `azd doctor check`](docs/images/doctor-check.svg)

## Installation

### Install (Published Extension)

Once this extension is published to an `azd` extension source registry, you can install it with:

```bash
azd extension install spboyer.azd.doctor
```

### Install (Local Development)

Recommended workflow is to use the `azd` developer extension (`azd x`):

```bash
# If needed, add the dev extension source (one-time)
azd extension source add -n dev -t url -l "https://aka.ms/azd/extensions/registry/dev"

# Install the developer extension (one-time)
azd extension install microsoft.azd.extensions

# From this repo
cd spboyer.azd.doctor
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

### `context`

Displays the context of the current AZD project and environment.

```bash
azd doctor context
```

## Development

### Prerequisites

- Go 1.24+
- Azure Developer CLI (`azd`)

### Build

```bash
cd spboyer.azd.doctor
go build -o ../bin/azd-ext-doctor .
```

### Install (Local Development)

For local development, use:

```bash
cd spboyer.azd.doctor
azd x watch
```

