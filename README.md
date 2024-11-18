# go-telemetry

This repository contains the solution to the Coding Challenge for Plenty One

## Environment variables

```YAML
GO_TELEMETRY_FILE_PATH=<DEFAULT:telemetry-config.yml> # The file name of the telemetry configuration YAML file. This file should be placed in project root.
```

## YAML Configuration file

```YAML
logger:
  level: <off|info|warning|debug|error> # Default: info
  outputWriter: <cli|jsonFile|textFile> # Default: cli
  outputDir: .                          # Default: root dir
```

## Test

`go test -v ./... [--cover]`
