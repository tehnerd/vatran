# Vatran

## Note from Author (aka tehnerd; the only non-vibecoded section of this project)
This project is fully vibecoded as i'm just fooling around with this brave new world of tooling. So use/play with
it at your own risk.

If you want to create a PR etc - dont. Since i could vibecode it as well. It is preferable just to add in Issue
w/ detailed feature you are missing/want to see implemented here.

TODOs (in no particular order):
1. Move work w/ UI (e.g. to define reals 
2. Add some fancy pictures of UI in Docs 
3. More examples on how it could be started
4. Proper auth workflow (https + login + oauth2 for cli clients etc)
5. Docker image for ease of use

## Abstract

Vatran is an HTTP server frontend for the Katran load balancer. It exposes a REST API and optional web UI for
managing Katran, which is an XDP-based L4 load balancer implemented as a C++ library plus BPF programs.

This repo provides:
- A Go HTTP/REST server that manages Katran via a CGO C API wrapper.
- A YAML-driven configuration flow for bringing up the load balancer and configuring VIPs/backends.
- An optional single-page UI served by the same HTTP server.

If you are looking for low-level Katran details (XDP/BPF behavior, environment requirements, etc.), start with
the upstream documentation in `external/katran/README.md` and `external/katran/USAGE.md`.
or visit https://github.com/facebookincubator/katran

## Architecture

Vatran stitches together three layers:
1. Katran (upstream) C++ library + BPF programs (in `external/katran/`)
2. A C API wrapper (`include/katran_capi.h`, `src/katran_capi.cpp`)
3. Go bindings and the REST server (`go/katran/`, `go/server/`)

The server exposes a REST API under `/api/v1` for lifecycle management, VIP/real management, and stats
queries (see `go/API.md`). It can also serve static UI assets if `static_dir` is configured.

## Features

- REST API for Katran load balancer lifecycle and configuration
- YAML config file support for server settings, load balancer config, VIPs, and backend pools
- TLS + optional mTLS support for the HTTP server
- Optional SPA UI served from a static directory

## Quick Start

### 1) Build the Katran BPF programs

Follow the upstream Katran requirements and build instructions:
- `external/katran/README.md`
- `external/katran/DEVELOPING.md`

From this repo, the helper script referenced in `BUILD.md` is:

```bash
./external/katran/build_bpf_modules_opensource.sh -s ./external/katran/
```

### 2) Build the Vatran server

```bash
go build ./go/cmd/katran-server
```

### 3) Configure and run

Copy and edit the example config:

```bash
cp go/config_example.yaml /tmp/vatran.yaml
$EDITOR /tmp/vatran.yaml
```

Start the server:

```bash
./katran-server -config /tmp/vatran.yaml
```

The REST API is served at `/api/v1`, and a health check is available at `/health`.

## Configuration

The canonical example configuration lives at `go/config_example.yaml`. It covers:
- HTTP server settings (bind address, TLS, timeouts, middleware)
- BPF program locations
- Katran load balancer settings
- Target groups (backend pools)
- VIPs and their target groups

## API

See `go/API.md` for endpoint details. The API includes:
- Load balancer lifecycle (create, load BPF, attach BPF, status, reload)
- VIP and backend management
- Stats, health checks, and feature flags

## UI

The UI lives in `ui/` and can be served by the Go server via the `static_dir` setting in the YAML config
or the `-static-dir` CLI flag. The design goals are outlined in `design_docs/ui.md`.

## License

Katran is licensed under GPLv2 (see `external/katran/LICENSE`). This repository includes Katran as a
submodule-style dependency in `external/katran/` and is governed by its licensing terms.
