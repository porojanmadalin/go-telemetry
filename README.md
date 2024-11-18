# go-telemetry

This repository contains the solution to the Coding Challenge for Plenty One.

The library implements a minimal logging tool and a transaction logging tool.

The Transaction logging tool groups individual logs to form a transaction. (e.g. Wishlist Feature: Add items to cart -> log -> Browse Similar Items -> log -> Remove an item from cart -> log...)

## Environment Variables

Environment variables are used to set up go-telemetry in a custom way, independent of the YAML file configuration.

```YAML
GO_TELEMETRY_FILE_PATH=<DEFAULT:telemetry-config.yml> # The path to the go-telemetry configuration YAML file. Default location is project root.
```

## YAML Configuration File

The YAML configuration file can be adjusted in order to remove the need to modify source code when for e.g. log level is adjusted.

```YAML
logger:
  level: <off|info|warning|debug|error> # Default: info
  outputWriter: <cli|jsonFile|textFile> # Default: cli
  outputDir: .                          # Default: root dir
```

## Test

Unit test coverage of 86.3%.

`go test -v ./... [--cover]`
