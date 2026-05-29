# DNSwitch

DNSwitch is a small Go-based utility for managing and switching system DNS settings across multiple operating systems.

## Overview

DNSwitch simplifies switching DNS servers temporarily or permanently. The project is designed to be cross-platform (Windows, Linux, macOS) and separates DNS logic from platform-specific network operations.

## Features

- Manage system DNS settings
- Cross-platform support (Windows, Linux, macOS)
- Load and apply custom DNS configurations
- Simple, modular structure for easy extension

## Project Structure

- `main.go`: program entry point.
- `go.mod`: Go module and dependency management.
- `dns/`: DNS-related logic
  - `dns.go`: public interfaces and general DNS functions.
  - `system.go`: system-level DNS application logic that calls platform-specific implementations.
  - `custom_dns.go`: read/write custom DNS configuration files.
  - `utils.go`: helper utilities related to DNS.
  - `vars.go`: global and default variables.
  - `custom_dns.json`: example or default custom DNS configuration ([dns/custom_dns.json](dns/custom_dns.json#L1)).
- `network/`: platform-specific network implementations
  - `windows.go`, `linux.go`, `darwin.go`: platform-specific routines to apply network/DNS changes.
  - `network.go`: public interface for network operations.
- `pkg/`: utility packages
  - `input.go`: input/argument handling (e.g., CLI args or interactive input).

## How It Works (high level)

1. The program reads configuration or user input to determine which DNS settings to apply.
2. The `dns` layer manages DNS details and delegates platform-specific work to the `network` layer.
3. Custom configurations can be provided via `dns/custom_dns.json` or interactive input.

## Configuration

- Example configuration file: [dns/custom_dns.json](dns/custom_dns.json#L1). It can contain a list of DNS servers, priorities, and optional labels.

## Build & Run

To build the binary:

```bash
go build -o DNSwitch main.go
```

To run directly during development:

```bash
go run main.go
```

Command-line arguments and runtime behavior are defined in the code; consult `pkg/input.go` for how inputs are parsed.

## Contributing

- To add support for a new OS or improve DNS behavior, extend the appropriate files in `network/` and `dns/`.
- Implementations should follow the interfaces in `dns.go` and `network.go`.

## Issues

- Please open issues on the repository for bugs or feature requests.

## License

No license file is included in this repository. Add a suitable license if you plan to publish or distribute the project.

---

If you want, I can add: CLI usage examples, a concrete `dns/custom_dns.json` sample, or a short usage guide. Tell me which you'd like next.
