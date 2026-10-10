# Payslip Automation System

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://go.dev/)
[![Wails](https://img.shields.io/badge/Wails-Desktop-3B82F6)](https://wails.io/)
[![SQLite](https://img.shields.io/badge/SQLite-Local%20DB-003B57?logo=sqlite)](https://sqlite.org/)
[![LibreOffice](https://img.shields.io/badge/LibreOffice-PDF%20Export-18A303?logo=libreoffice)](https://www.libreoffice.org/)

A payroll and payslip automation project for school staff data. It imports Excel-based payroll information, stores it in SQLite, and generates monthly and yearly payslips in both Excel and PDF form.

## Table of contents

- [Overview](#overview)
- [Features](#features)
- [Project structure](#project-structure)
- [Workflow](#workflow)
- [Getting started](#getting-started)
- [Run the importer](#run-the-importer)
- [Run the desktop app](#run-the-desktop-app)
- [Output layout](#output-layout)
- [Development notes](#development-notes)
- [Useful commands](#useful-commands)

## Overview

This repository contains two main execution paths:

- `cmd/importer` — batch payroll import and generation pipeline
- `cmd/desktop` — desktop app for browsing and generating payslips

Shared logic lives under `internal/` and handles parsing, database access, generation, and PDF conversion. The project is designed around school payroll processing and outputs structured payslip folders for each period.

## Features

- Import payroll data from Excel files
- Parse employee, school, and cluster information
- Store records in SQLite for easy querying and exports
- Generate monthly payslips from Excel templates
- Generate yearly payslips for financial-year reporting
- Convert generated XLSX files into PDFs using LibreOffice
- Support bulk generation by scope, cluster, or school
- Browse and generate payslips from a desktop interface

## Project structure

```text
.
├── cmd/
│   ├── desktop/             # Wails desktop application
│   └── importer/            # CLI payroll import/generation entrypoint
├── internal/
│   ├── concurrency/         # concurrent PDF conversion workers
│   ├── database/            # SQLite setup and queries
│   ├── excel/               # Excel reading and helpers
│   ├── generator/           # payslip generation orchestration
│   ├── genExcel/            # Excel template generation
│   ├── genPDF/              # PDF conversion logic
│   ├── importer/            # import workflow and reports
│   ├── models/              # domain data structures
│   ├── parser/              # parser logic for employee/school files
│   └── ...
├── data/                    # templates and static data files
├── source/                  # raw payroll source files
├── output/                  # generated payslips and PDFs
├── go.mod                   # root Go module
├── cmd/desktop/go.mod       # desktop app module
├── README.md                # project documentation
├── payslip.db              # DB created at runtime
└── ...
```

## Workflow

1. A payroll Excel file is loaded or placed in `source/`.
2. The importer detects the pay period and validates the file format.
3. Employee and salary records are stored in SQLite.
4. Monthly payslip templates are generated into `output/`.
5. The generated Excel files are converted to PDF.
6. Yearly payslips are generated using the financial year data.
7. The Wails app can then query and generate payslips for individuals or groups.

## Getting started

### Prerequisites

- Go 1.26+
- Wails CLI for the desktop app
- LibreOffice with `soffice` available in your `PATH`

### Install dependencies

```bash
go mod download
```

## Run the importer

The CLI importer is the main batch-processing entrypoint:

```bash
go run ./cmd/importer
```

This command reads the configured source file(s), stores imported data in SQLite, and generates monthly/yearly payslip output under `output/`.

## Run the desktop app

The desktop app is a separate Go module, so it should be run from its directory:

```bash
cd cmd/desktop
wails dev
```

To build a production desktop app:

```bash
cd cmd/desktop
wails build
```

## Output layout

Generated files are stored in a period-based folder structure such as:

```text
output/
├── August_2026/
│   └── ClusterName/
│       └── SchoolName/
│           ├── Employee.xlsx
│           └── PDF/
│               └── Employee.pdf
├── Financial_Year_2026_2027/
│   └── ClusterName/
│       └── SchoolName/
│           ├── Employee.xlsx
│           └── PDF/
│               └── Employee.pdf
```

This makes it easier to trace payouts by month, financial year, school, and cluster.

## Development notes

- The desktop app materializes bundled Excel templates into its runtime resources directory on startup.
- The project uses a multi-module Go layout on purpose:
  - root module for shared logic and CLI tooling
  - `cmd/desktop` for the Wails app
- Output directories and paths are configurable from the desktop app.
- This project is designed for a school payroll workflow and expects source templates and raw payroll files to be available locally.

## Useful commands

```bash
# root module checks
go test ./...

# desktop app build
cd cmd/desktop && go build ./...

# run importer
go run ./cmd/importer

# run desktop app
cd cmd/desktop && wails dev
```

## License

This project is intended for internal payroll processing and is not a general-purpose public package unless otherwise specified by the repository owner.
