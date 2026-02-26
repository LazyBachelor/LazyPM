# LazyPM

A project management tool for developers. 

Made for bachelor project, used for gathering and comparing data between interfaces.

## Overview

LazyPM is a lightweight project management system that provides three interfaces:
- **CLI** - Command-line interface for quick operations
- **TUI** - Terminal user interface for interactive use
- **Web** - Browser-based interface with HTML templates
- **Survey** - Interactive task-based survey system for data gathering

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
- npm (for Tailwind CSS)

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
make build              # Creates bin/pm, bin/tui, bin/web, bin/survey

# Run specific interfaces
make cli                # Run CLI interface
make tui                # Run TUI interface (interactive)
make web                # Run Web server (localhost:8080)
make start              # Run survey interface

# Development with hot reload
make dev                # Watch templ files and auto-reload web
make tw                 # Watch Tailwind CSS changes

# Docker (lazyos desktop environment)
make os-build           # Build lazyos Docker image
make os-run             # Run lazyos container (localhost:3000-3001)
make os-stop            # Stop and remove lazyos container

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
├── build/              # Docker build files for lazyos desktop environment
│   ├── Dockerfile     # Webtop-based desktop container
│   ├── setup-desktop.sh # Desktop shortcut setup script
│   └── survey         # Survey binary (copied during build)
├── cmd/                # Entry points
│   ├── pm/            # CLI main
│   ├── tui/           # TUI main
│   ├── web/           # Web server main
│   └── survey/        # Survey system main
├── internal/           # Core implementation
│   ├── models/        # Data models (beads types)
│   ├── service/       # Business logic
│   ├── storage/       # Data persistence
│   └── commands/      # Cobra commands
│       ├── issues/    # Issue management commands
│       └── survey/    # Survey commands
├── pkg/                # Public packages
│   ├── cli/           # CLI framework
│   ├── tui/           # TUI views and components
│   ├── web/           # Web handlers, templates, assets
│   ├── repl/          # REPL implementation
│   └── task/          # Task registration and runner
└── .pm/                # Local data storage (gitignored)
```

## Technology Stack

- [Go](https://golang.org/) 1.25.6 - Backend language
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Fang](https://github.com/charmbracelet/fang) - Colored CLI output
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Huh](https://github.com/charmbracelet/huh) - Terminal forms
- [Templ](https://github.com/a-h/templ) - HTML templating
- [HTMX](https://htmx.org/) - Dynamic web UI
- [Tailwind CSS v4](https://tailwindcss.com/) + DaisyUI - Web styling
- [Beads](https://github.com/steveyegge/beads) - Issue tracking engine
- [go-git](https://github.com/go-git/go-git) v6 - Git operations

## Configuration

LazyPM stores data in a local `.pm` directory:
- `.pm/db.db` - SQLite database for issues
- `.pm/stats.json` - Statistics storage

The directory is automatically created on first run.

## Development

### Running the Survey System

The survey system provides interactive task-based data gathering:

```bash
make start              # Start the survey interface
```

### LazyOS Desktop Environment

LazyOS is a containerized desktop environment for running the survey in a browser-based desktop:

```bash
make os-build           # Build lazyos Docker image
make os-run             # Run lazyos container (access at localhost:3000)
make os-stop            # Stop and remove lazyos container
```

The lazyos container provides:
- Web-based desktop environment (LinuxServer Webtop)
- Pre-configured desktop shortcut for the survey
- Runs at http://localhost:3000

### Web Development

```bash
# Terminal 1 - Auto-reload Go server on templ changes
make dev

# Terminal 2 - Watch Tailwind CSS
make tw
```

### Available Tasks

- `create_issue` - Create a new issue
- `coding_task` - Coding-related task
- `git_task` - Git operations task
- `sprint_planning` - Sprint planning task
- `issue_triage` - Issue triage task
- `milestone_tracking` - Milestone tracking task
- `dependency_management` - Dependency management task
- `team_capacity` - Team capacity task
- `report_generation` - Report generation task
- `stakeholder_update` - Stakeholder update task
- `priority_management` - Priority management task
- `backlog_refinement` - Backlog refinement task

## CI/CD

This project uses GitHub Actions for:
- Building Go binaries on push/PR
- Docker image publishing
- Release automation with GoReleaser

See `.github/workflows/` for details.

## License

See [LICENSE](LICENSE) file for details.