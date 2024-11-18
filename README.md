# go-telemetry

This repository contains the solution to the Coding Challenge for Plenty One

## Environment variables

```YAML
<<<<<<< HEAD
GO_TELEMETRY_FILE_PATH=<DEFAULT:telemetry-config.yml> # The path to the go-telemetry configuration YAML file. Default location is project root.
=======
GO_TELEMETRY_FILE_PATH=<DEFAULT:telemetry-config.yml> # The file name of the telemetry configuration YAML file. This file should be placed in project root.
>>>>>>> 209bbb7e35dd2588c99fe44e37b6149cb8c881bd
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
