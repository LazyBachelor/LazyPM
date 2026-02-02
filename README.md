# LazyPM

A lazy project management tool for developers.

## Overview

LazyPM is a lightweight project management system that provides three interfaces:
- **CLI** - Command-line interface for quick operations
- **TUI** - Terminal user interface for interactive use
- **Web** - Browser-based interface with HTML templates

## Features

- Issue tracking (bugs, features, tasks, epics, chores)
- Status management (open, in-progress, blocked, deferred, closed)
- Dependency tracking between issues
- Labels and comments
- Statistics and reporting
- SQLite storage via Beads library

## Quick Start

```bash
# Install dependencies
go mod tidy
```

### Install make

Linux:
```bash
sudo pacman -S make
```

Windows (Choco):
```bash
choco install make
```

### Available make commands

```bash
# Build all interfaces
make build

# Run CLI
make cli

# Run TUI
make tui

# Run Web interface
make web

# Development with hot reload
make dev
```

## Project Structure

```
├── cmd/           # Entry points (cli, tui, web)
├── internal/      # Core services and models
├── pkg/           # Public packages (CLI, TUI, Web)
└── .pm/           # Local data storage
```

## Requirements

- Go 1.25.6+
- SQLite

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [templ](https://github.com/a-h/templ) - HTML templating
- [beads](https://github.com/steveyegge/beads) - Issue tracking engine
