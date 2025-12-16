# AZD Doctor Extension

An `azd` extension that checks for necessary prerequisites based on the current project's configuration.

## Features

- Checks for Container Runtime (Docker/Podman)
- Checks for Language Runtimes based on `azure.yaml` (Node.js, Python, .NET)

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

## Usage

Once installed, you can run the doctor check:

```bash
azd doctor check
```
