# AZD Doctor Extension

An `azd` extension that checks for necessary prerequisites based on the current project's configuration.

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

### `context`

Displays the context of the current AZD project and environment.

```bash
azd doctor context
```

## Development

### Prerequisites

- Go 1.21+
- Azure Developer CLI (`azd`)

### Build

```bash
go build -o azd-ext-doctor .
```

### Install (Local Development)

To develop and test this extension, it is recommended to use the `azd` developer extension tools (`azd x`).

1. Initialize the extension (if not already done):
   ```bash
   azd x init
   ```

2. Build and install locally:
   ```bash
   azd x build
   ```

3. Watch for changes:
   ```bash
   azd x watch
   ```

