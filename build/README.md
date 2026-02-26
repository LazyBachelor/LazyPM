# LazyOS Build Directory

This directory contains files for building the LazyOS desktop environment - a containerized web-based desktop for running the LazyPM survey system.

## Contents

- **Dockerfile** - Container definition based on LinuxServer Webtop
- **setup-desktop.sh** - Initialization script that creates desktop shortcuts
- **survey** - Survey binary (copied from bin/ during build)

## LazyOS Overview

LazyOS provides a browser-based desktop environment where users can run the LazyPM survey tasks in an isolated container.

It's useful for:
- Running the survey system in a controlled environment
- Providing a desktop-like experience through a web browser
- Testing the survey interface without local setup

## Building LazyOS

From the project root:

```bash
make os-build
```

This will:
1. Build the survey binary
2. Copy it to this directory
3. Build the Docker image tagged as `lazyos:latest`

## Running LazyOS

```bash
make os-run
```

Access the desktop environment at http://localhost:3000

The container exposes :3000 for http and :3001 for https

## Stopping LazyOS

```bash
make os-stop
```

This stops and removes the lazyos container.

## Dockerfile Details

The Dockerfile:
- Uses `lscr.io/linuxserver/webtop:latest` as base
- Installs gcompat and libc6-compat for binary compatibility
- Copies the survey binary to `/usr/local/bin/survey`
- Includes setup-desktop.sh for desktop shortcut creation

## Desktop Setup

The setup-desktop.sh script runs on container startup and:
- Creates a Desktop folder for the 'abc' user
- Creates a desktop shortcut named "PM Survey"