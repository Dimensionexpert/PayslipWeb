# Payslip

This repository contains a Go project with two separate executable entrypoints:

- a CLI tool under `cmd/importer`
- a Wails desktop app under `cmd/desktop`
- shared application logic under `internal/`

## Project layout

- `cmd/importer` - command-line payroll import/generation tool
- `cmd/desktop` - Wails desktop application
- `internal/` - shared Go packages for database, generation, parsing, and conversion logic
- `data/` - templates and input data
- `source/` - source Excel files
- `output/` - generated output files

## Run the CLI

```bash
go run ./cmd/importer
```

## Run the desktop app

```bash
cd cmd/desktop
wails dev
```

For production build:

```bash
cd cmd/desktop
wails build
```

## Important note about Go modules

This repo intentionally uses a multi-app Go layout:

- the root module is for shared packages and CLI code
- `cmd/desktop` is a separate Go module for the Wails app

Do not create a second Wails app at the repository root. Keep only one Wails app entrypoint in `cmd/desktop`.

## Common commands

```bash
# from repo root
 go test ./...

# desktop app
 cd cmd/desktop && go build ./...
```
