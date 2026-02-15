# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

IDL (Image Down Loader) is a Go CLI tool that downloads images from a CSV file in parallel and optionally preprocesses them (e.g., resizing) before saving to disk.

## Build & Development Commands

```bash
make build          # Build binary from ./cmd/idl
make lint           # Run golangci-lint
go test ./...       # Run all tests
go test ./internal/config/...  # Run tests for a single package
```

## Architecture

**Entry point:** `cmd/idl/main.go` — parses CLI config, sets up signal handling, runs the orchestrator.

**Core orchestrator:** `internal/idl/idl.go` — `IDL.Run()` implements a channel-based parallel pipeline:
1. CSV reading (goroutine) → 2. Multi-worker downloading with exponential backoff → 3. Multi-worker preprocessing → 4. Multi-worker file saving with progress bar

All stages communicate via Go channels and respect context cancellation for graceful shutdown.

**Key packages:**
- `internal/config` — CLI flag parsing via Kong (`--input`, `--dir`, `--workers`, `--timeout`, `--resize`, `--url-field`, `--file-field`, `--reload`, `--retries`)
- `internal/csv` — CSV reader with configurable field name mapping
- `internal/dl` — HTTP downloader with `cenkalti/backoff` exponential retry
- `internal/preprocess` — `Preprocessor` interface with implementations: `empty` (no-op), `resize` (via `disintegration/imaging`), `compose` (chains multiple preprocessors)
- `internal/utils` — `FileExists` helper

## Key Dependencies

- `github.com/alecthomas/kong` — CLI parsing
- `github.com/disintegration/imaging` — image resize/encode
- `github.com/cenkalti/backoff/v5` — retry with exponential backoff
- `github.com/schollz/progressbar/v3` — terminal progress bar

## Go Version

Go 1.24.0 (see `go.mod`)
