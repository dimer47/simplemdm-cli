# simplemdm-cli

CLI for the SimpleMDM API — manage your Apple devices from the terminal.

[Documentation en francais](README.fr.md)

## Features

- **116 MCP tools** and **130+ CLI commands** covering the full SimpleMDM API
- **Secure API key storage**: stored in your system's credential manager (macOS Keychain, Windows Credential Manager, Linux Secret Service)
- **Multi-context**: manage multiple SimpleMDM accounts
- **Flexible output**: table, JSON, YAML, CSV
- **Cross-platform**: macOS, Linux, Windows (amd64 and arm64)
- **MCP integration**: 116 tools for Claude Code, VS Code, JetBrains
- **Self-update**: automatic update checks and one-command updates

## Prerequisites

- A [SimpleMDM](https://simplemdm.com) account
- A **SimpleMDM API key** (created from [Settings > API](https://a.simplemdm.com/admin/settings/api))

## Installation

### Method 1: Download the binary (recommended)

Go to the [Releases](https://github.com/dimer47/simplemdm-cli/releases/latest) page and download the archive for your platform.

Or with a single command:

**macOS (Apple Silicon -- M1/M2/M3/M4):**

```bash
curl -sL https://github.com/dimer47/simplemdm-cli/releases/latest/download/simplemdm-cli_darwin_arm64.tar.gz | tar xz
sudo mv simplemdm-cli /usr/local/bin/
```

**macOS (Intel):**

```bash
curl -sL https://github.com/dimer47/simplemdm-cli/releases/latest/download/simplemdm-cli_darwin_amd64.tar.gz | tar xz
sudo mv simplemdm-cli /usr/local/bin/
```

**Linux (amd64):**

```bash
curl -sL https://github.com/dimer47/simplemdm-cli/releases/latest/download/simplemdm-cli_linux_amd64.tar.gz | tar xz
sudo mv simplemdm-cli /usr/local/bin/
```

**Linux (arm64 -- Raspberry Pi, etc.):**

```bash
curl -sL https://github.com/dimer47/simplemdm-cli/releases/latest/download/simplemdm-cli_linux_arm64.tar.gz | tar xz
sudo mv simplemdm-cli /usr/local/bin/
```

**Windows:**

Download `simplemdm-cli_windows_amd64.zip` from the [Releases](https://github.com/dimer47/simplemdm-cli/releases/latest) page, extract and add the folder to your `PATH`.

### Method 2: From source (requires Go 1.25+)

```bash
go install github.com/dimer47/simplemdm-cli@latest
```

The binary will be installed in `$GOPATH/bin/` (usually `~/go/bin/`). Make sure this directory is in your `PATH`.

### Method 3: Build locally

```bash
git clone https://github.com/dimer47/simplemdm-cli.git
cd simplemdm-cli
go build -o simplemdm-cli .
./simplemdm-cli version
```

### Verify installation

```bash
simplemdm-cli version
# simplemdm-cli 1.0.0 (abc1234) built 2026-05-11T10:00:00Z
```

## Updating

The CLI automatically checks for new versions at startup and notifies you when an update is available.

```bash
# Update to the latest version
simplemdm-cli self-update
```

The update is downloaded from GitHub Releases and replaces the current binary in place. If the binary is in a protected directory (e.g. `/usr/local/bin/`), try running with `sudo`.

## Quick Start

### 1. Get a SimpleMDM API key

1. Log in to the [SimpleMDM admin console](https://a.simplemdm.com)
2. Go to **Settings > API**
3. Copy your API key

### 2. Configure the CLI

```bash
simplemdm-cli auth login
```

Answer the 2 prompts:
```
Context name (default): Enter          # Press Enter for "default"
SimpleMDM API key: ••••••••••••        # Paste your API key (hidden)
```

The API key is stored in your **system credential manager** (macOS Keychain, Windows Credential Manager, or Linux Secret Service) -- encrypted, never written in plain text on disk.

### 3. Test it

```bash
# Get your account details
simplemdm-cli account get

# List your devices
simplemdm-cli device list

# JSON output
simplemdm-cli device list --json
```

## Configuration

### API key resolution priority

| Priority | Source | Use case |
|----------|--------|----------|
| 1 | Flag `--api-key` / `-k` | One-off tests |
| 2 | Env var `SMDM_API_KEY` | CI/CD, scripts |
| 3 | System credential manager | Daily use (via `auth login`) |

### Multi-context (multiple SimpleMDM accounts)

```bash
# Configure a "production" context
simplemdm-cli auth login
# -> Enter "production" as context name

# Configure a "staging" context
simplemdm-cli auth login
# -> Enter "staging" as context name

# List all contexts (* = active)
simplemdm-cli auth list
# * production (****abcd1234)
#   staging    (****efgh5678)

# Switch context
simplemdm-cli auth switch staging

# Use a context for a single command
simplemdm-cli device list --context production

# Check status (is the key still valid?)
simplemdm-cli auth status

# Remove a context
simplemdm-cli auth remove staging
```

## Usage

### Account

```bash
simplemdm-cli account get                                    # Account details
simplemdm-cli account update --name "My Company"             # Update account
```

### Devices

```bash
simplemdm-cli device list                                    # List all devices
simplemdm-cli device list --json                             # JSON output
simplemdm-cli device list --search "MacBook"                 # Search
simplemdm-cli device get <id>                                # Device details
simplemdm-cli device create --name "New Device"              # Create
simplemdm-cli device update <id> --name "Renamed"            # Update name
simplemdm-cli device delete <id>                             # Delete
simplemdm-cli device refresh <id>                            # Refresh info
simplemdm-cli device push-apps <id>                          # Push apps
simplemdm-cli device lock <id> --message "Lost" --pin 1234   # Lock
simplemdm-cli device wipe <id>                               # Wipe
simplemdm-cli device restart <id>                            # Restart
simplemdm-cli device shutdown <id>                           # Shutdown
simplemdm-cli device clear-passcode <id>                     # Clear passcode
simplemdm-cli device update-os <id>                          # Update OS
simplemdm-cli device unenroll <id>                           # Unenroll
```

### Device -- Lost Mode (iOS supervised)

```bash
simplemdm-cli device lost-mode-enable <id> --message "Call IT" --phone-number "+1234567890"
simplemdm-cli device lost-mode-disable <id>
simplemdm-cli device lost-mode-play-sound <id>
simplemdm-cli device lost-mode-update-location <id>
```

### Device -- Advanced commands

```bash
simplemdm-cli device bluetooth-enable <id>                   # Enable Bluetooth
simplemdm-cli device bluetooth-disable <id>                  # Disable Bluetooth
simplemdm-cli device remote-desktop-enable <id>              # Enable Remote Desktop
simplemdm-cli device remote-desktop-disable <id>             # Disable Remote Desktop
simplemdm-cli device rotate-firmware-password <id>           # Rotate firmware password
simplemdm-cli device rotate-recovery-lock <id>               # Rotate recovery lock
simplemdm-cli device rotate-filevault-key <id>               # Rotate FileVault key
simplemdm-cli device set-admin-password <id> --new-password "..."
simplemdm-cli device rotate-admin-password <id>
simplemdm-cli device clear-firmware-password <id>
simplemdm-cli device clear-recovery-lock <id>
simplemdm-cli device clear-restrictions-password <id>
simplemdm-cli device set-timezone <id> --timezone "America/New_York"
```

### Device -- Custom Attributes & Users

```bash
simplemdm-cli device custom-attributes <id>                  # List attribute values
simplemdm-cli device set-custom-attribute <id> department --value "Engineering"
simplemdm-cli device profiles <id>                           # List device profiles
simplemdm-cli device users <id>                              # List device users
simplemdm-cli device delete-user <id> <user-id>              # Delete user
```

### Apps

```bash
simplemdm-cli app list                                       # List all apps
simplemdm-cli app get <id>                                   # App details
simplemdm-cli app create --name "My App" --app-store-id 123  # Create from App Store
simplemdm-cli app create --binary ./app.ipa                  # Upload enterprise app
simplemdm-cli app update <id> --binary ./app.ipa             # Update binary
simplemdm-cli app delete <id>                                # Delete
simplemdm-cli app installs <id>                              # List installs
simplemdm-cli app managed-configs <id>                       # List managed configs
simplemdm-cli app managed-config-create <id> --key "url" --value "https://..."
simplemdm-cli app managed-configs-push <id>                  # Push configs to devices
simplemdm-cli app managed-config-delete <id> <config-id>     # Delete config
```

### Assignment Groups

```bash
simplemdm-cli assignment-group list                          # List groups
simplemdm-cli assignment-group get <id>                      # Group details
simplemdm-cli assignment-group create --name "Team"          # Create
simplemdm-cli assignment-group update <id> --name "New Name" # Update
simplemdm-cli assignment-group delete <id>                   # Delete
simplemdm-cli assignment-group assign-app <id> --app-id 10   # Assign app
simplemdm-cli assignment-group assign-device <id> --device-id 121
simplemdm-cli assignment-group push-apps <id>                # Push apps
simplemdm-cli assignment-group clone <id>                    # Clone group
```

### Profiles

```bash
simplemdm-cli profile list                                   # List profiles
simplemdm-cli profile get <id>                               # Profile details
simplemdm-cli profile assign-device <id> --device-id 121     # Assign to device
simplemdm-cli profile unassign-device <id> --device-id 121   # Unassign
```

### Custom Configuration Profiles

```bash
simplemdm-cli custom-configuration-profile list              # List profiles
simplemdm-cli custom-configuration-profile delete <id>       # Delete
simplemdm-cli custom-configuration-profile push-device <id> --device-id 121
simplemdm-cli custom-configuration-profile remove-device <id> --device-id 121
```

### Custom Attributes

```bash
simplemdm-cli custom-attribute list                          # List attributes
simplemdm-cli custom-attribute get <id>                      # Attribute details
simplemdm-cli custom-attribute create --name "department"    # Create
simplemdm-cli custom-attribute delete <id>                   # Delete
```

### Custom Declarations

```bash
simplemdm-cli custom-declaration list                        # List declarations
simplemdm-cli custom-declaration delete <id>                 # Delete
```

### Device Groups

```bash
simplemdm-cli device-group list                              # List groups
simplemdm-cli device-group get <id>                          # Group details
simplemdm-cli device-group assign-device <id> --device-id 121
```

### DEP Servers

```bash
simplemdm-cli dep-server list                                # List DEP servers
simplemdm-cli dep-server get <id>                            # Server details
simplemdm-cli dep-server devices <id>                        # List DEP devices
simplemdm-cli dep-server sync <id>                           # Sync server
```

### Enrollments

```bash
simplemdm-cli enrollment list                                # List enrollments
simplemdm-cli enrollment get <id>                            # Enrollment details
simplemdm-cli enrollment delete <id>                         # Delete
simplemdm-cli enrollment send-invitation <id> --contact "user@company.com"
```

### Installed Apps

```bash
simplemdm-cli installed-app get <id>                         # App details
simplemdm-cli installed-app delete <id>                      # Delete (uninstall)
simplemdm-cli installed-app update <id>                      # Update app
```

### Scripts & Script Jobs

```bash
simplemdm-cli script list                                    # List scripts
simplemdm-cli script get <id>                                # Script details
simplemdm-cli script delete <id>                             # Delete
simplemdm-cli script-job list                                # List jobs
simplemdm-cli script-job get <id>                            # Job details
simplemdm-cli script-job create --script-id 100 --device-ids "121,122"
simplemdm-cli script-job cancel <id>                         # Cancel job
```

### Logs

```bash
simplemdm-cli log list                                       # List logs
simplemdm-cli log get <id>                                   # Log details
```

### Push Certificate

```bash
simplemdm-cli push-certificate get                           # Certificate details
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SMDM_API_KEY` | SimpleMDM API key | -- |
| `SMDM_OUTPUT` | Output format: `table`, `json`, `yaml`, `csv` | `table` |
| `SMDM_DEBUG` | Debug mode (`true`/`false`) | `false` |
| `SMDM_CONTEXT` | Active context | `default` |
| `NO_COLOR` | Disable colors | -- |

## Shell Completion

```bash
# Bash
simplemdm-cli completion bash > /etc/bash_completion.d/simplemdm-cli

# Zsh (add to your .zshrc)
simplemdm-cli completion zsh > "${fpath[1]}/_simplemdm-cli"

# Fish
simplemdm-cli completion fish > ~/.config/fish/completions/simplemdm-cli.fish

# PowerShell
simplemdm-cli completion powershell > simplemdm-cli.ps1
```

## MCP Integration (Claude Code, VS Code, JetBrains)

The CLI includes a built-in [MCP (Model Context Protocol)](https://modelcontextprotocol.io/) server exposing **116 tools** for AI assistants.

### Setup

Add to your Claude Code settings (or VS Code / JetBrains with the Claude extension):

```json
{
  "mcpServers": {
    "simplemdm": {
      "command": "simplemdm-cli",
      "args": ["mcp-serve"],
      "env": {
        "SMDM_API_KEY": "your-api-key-here"
      }
    }
  }
}
```

> If you already configured the API key via `simplemdm-cli auth login`, the MCP server will automatically use your system credential store -- no need for the `SMDM_API_KEY` env var.

### Available MCP Tools

| Category | Tools | Description |
|----------|-------|-------------|
| **Account** | `account-get`, `account-update` | Manage account details |
| **Apps** | `app-list`, `app-get`, `app-create`, `app-update`, `app-delete`, `app-installs`, `app-managed-configs`, `app-managed-config-create`, `app-managed-configs-push`, `app-managed-config-delete` | Full app lifecycle |
| **Assignment Groups** | `assignment-group-list`, `assignment-group-get`, `assignment-group-create`, `assignment-group-update`, `assignment-group-delete`, `assignment-group-assign-app`, `assignment-group-assign-device`, `assignment-group-push-apps`, `assignment-group-clone`, ... | Group management |
| **Custom Attributes** | `custom-attribute-list`, `custom-attribute-get`, `custom-attribute-create`, `custom-attribute-delete` | Custom attributes |
| **Custom Config Profiles** | `custom-configuration-profile-list`, `custom-configuration-profile-delete`, `custom-configuration-profile-push-device`, `custom-configuration-profile-remove-device` | Profile management |
| **Custom Declarations** | `custom-declaration-list`, `custom-declaration-delete` | Declaration management |
| **DEP Servers** | `dep-server-list`, `dep-server-get`, `dep-server-devices`, `dep-server-sync` | DEP management |
| **Device Groups** | `device-group-list`, `device-group-get`, `device-group-assign-device` | Group management |
| **Devices** | `device-list`, `device-get`, `device-create`, `device-update`, `device-delete`, `device-lock`, `device-wipe`, `device-restart`, `device-shutdown`, `device-refresh`, `device-push-apps`, `device-update-os`, `device-lost-mode-*`, ... | Full device management |
| **Enrollments** | `enrollment-list`, `enrollment-get`, `enrollment-delete`, `enrollment-send-invitation` | Enrollment management |
| **Installed Apps** | `installed-app-get`, `installed-app-delete`, `installed-app-update` | App management |
| **Logs** | `log-list`, `log-get` | Log access |
| **Profiles** | `profile-list`, `profile-get`, `profile-assign-device`, `profile-unassign-device` | Profile management |
| **Push Certificate** | `push-certificate-get` | Certificate info |
| **Scripts** | `script-list`, `script-get`, `script-delete` | Script management |
| **Script Jobs** | `script-job-list`, `script-job-get`, `script-job-create`, `script-job-cancel` | Job management |

### Usage in Claude Code

Once configured, you can simply say:

- *"List my SimpleMDM devices"*
- *"Lock device 121 with the message 'Contact IT'"*
- *"What apps are installed on device 122?"*
- *"Push all apps to the Engineering assignment group"*
- *"Create a script job for script 100 on devices 121 and 122"*

Claude will automatically call the right MCP tools.

## Development

```bash
# Clone
git clone https://github.com/dimer47/simplemdm-cli.git
cd simplemdm-cli

# Build
go build -o simplemdm-cli .

# Run tests
go test ./...

# Lint
go vet ./...
```

### Creating a new release

```bash
git tag v1.0.0
git push origin v1.0.0
# GitHub Actions builds and publishes automatically
```

## License

MIT
