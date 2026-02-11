# LazyPM

A project management tool for developers. 

Made for bachelor project, used for gathering and comparing data between interfaces.

## Overview

LazyPM is a lightweight project management system that provides three interfaces:
- **CLI** - Command-line interface for quick operations
- **TUI** - Terminal user interface for interactive use
- **Web** - Browser-based interface with HTML templates

## Features

- Issue tracking (bugs, features, tasks, epics, chores)
- Status management (open, in-progress, closed)
- Dependency tracking between issues
- Labels and comments
- Statistics and reporting
- SQLite storage via Beads library

## Quick Start

### Prerequisites

- Go 1.25.6+
- Make

### Installation

```bash
# Install dependencies
go mod tidy
```

### Make Installation

Linux:
```bash
sudo pacman -S make
```

Windows (Choco):
```bash
# Install choco
winget install --id chocolatey.chocolatey --source winget

# Install make
choco install make
```

## Build Commands

```bash
# Build all binaries
make build              # Creates bin/pm, bin/tui, bin/web

# Run specific interfaces
make cli                # Run CLI interface
make tui                # Run TUI interface (interactive)
make web                # Run Web server (localhost:8080)

# Development with hot reload
make dev                # Watch templ files and auto-reload web
make tw                 # Watch Tailwind CSS changes
make watch              # Run both dev and tw in parallel

# Maintenance
make tidy               # Run go mod tidy
make clean              # Clean build artifacts

# Install
make install-cli        # Install CLI with shell completions
make completions        # Generate shell completion scripts
```

## Project Structure

```
├── bin/                # Compiled binaries
├── cmd/                # Entry points
│   ├── pm/            # CLI main.go
│   ├── tui/           # TUI main.go
│   └── web/           # Web server main.go
├── internal/           # Core implementation
│   ├── models/        # Data models (beads types)
│   ├── service/       # Business logic (beads, statistics)
│   └── storage/       # Data persistence
├── pkg/                # Public packages
│   ├── cli/           # CLI commands and REPL
│   ├── tui/           # TUI views and components
│   └── web/           # Web handlers, templates, assets
└── .pm/                # Local data storage (gitignored)
```

## Technology Stack

- [Go](https://golang.org/) 1.25.6 - Backend language
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Templ](https://github.com/a-h/templ) - HTML templating
- [Tailwind CSS v4](https://tailwindcss.com/) + DaisyUI - Web styling
- [Beads](https://github.com/steveyegge/beads) - Issue tracking engine

## Configuration

LazyPM stores data in a local `.pm` directory:
- `.pm/db.db` - SQLite database for issues
- `.pm/stats.json` - Statistics storage

The directory is automatically created on first run.