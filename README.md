
# multi-profile

[![PkgGoDev](https://pkg.go.dev/badge/github.com/bygui86/multi-profile)](https://pkg.go.dev/github.com/bygui86/multi-profile)

`/!\ WARNING` Please notice this package is still under development.

Multi-profiling support package for Go.

This project was inspired by [pkg/profile](github.com/pkg/profile) but there is a fundamental difference: multi-profile offers the possibility to start multiple profiling at the same time.

## TODO list

- [x] ~~align with profile original (last align 1.10.2020)~~
- [x] ~~testing~~
- [x] ~~improve os.exit codes~~
- [ ] more advanced logger (e.g. zap, logrus) through interface
- [ ] github actions - `WIP`
    - [ ] build stage for all all branches
    - [ ] test stage for all branches
    - [ ] `TBD` release stage for tags
- [ ] README - `WIP`
    - [x] ~~complete all sections~~
    - [ ] github actions badge in readme, e.g. from pkg/profile [![Build Status](https://travis-ci.org/pkg/profile.svg?branch=master)](https://travis-ci.org/pkg/profile)
    - [x] ~~godoc badge in readme, e.g.~~

## Installation

```shell script
go get github.com/bygui86/multi-profile
```

## Usage

Enabling profiling in your application is as simple as one line at the top of your main function.

For example:

```go
package main

import "github.com/bygui86/multi-profile"

func main() {
    defer profile.CPUProfile(&profile.ProfileConfig{}).Start().Stop()
    
    // ...
}
```

Using profile specific method, you can create the kind of profiling you want giving a ProfileConfig as input. 

```go
package main

import "github.com/bygui86/multi-profile"

func main() {
    defer profile.CPUProfile(&profile.ProfileConfig{}).Start().Stop()
    defer profile.MemProfile(&profile.ProfileConfig{}).Start().Stop()
    defer profile.GoroutineProfile(&profile.ProfileConfig{}).Start().Stop()

    // ...
}
```

`(i)️ INFO` see [examples](examples/) folder for all available profiles and samples.

## Options

`(i)️️ INFO` see [examples](examples/) for all usage samples.

### Path

You can customize the path in which a profile file is going to be written. Use field `Path` and `UseTempPath` in the ProfileConfig.

### Shutdown hook

You can disable the shutdown hook. The shutdown hook controls if profiling package should hook SIGINT to write profiles cleanly. Use `DisableShutdownHook` field in the ProfileConfig.

### Quiet mode

You can suppress all logs. Use field `Quiet` in the ProfileConfig.

### Closer function

You can call a function right after stopping the profiling. Use `CloserHook` field in the ProfileConfig.

## Error codes

| Code | Description |
| --- | --- |
| 11 | path preparation failed |
| 12 | file creation failed |
| 13 | cpu profile start failed |
| 14 | trace profile start failed |

## Contributing

I welcome pull requests, bug fixes and issue reports.

To propose an extensive change, please discuss it first by opening an issue.
