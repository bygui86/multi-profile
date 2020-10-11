
# multi-profile

[![PkgGoDev](https://pkg.go.dev/badge/github.com/bygui86/multi-profile)](https://pkg.go.dev/github.com/bygui86/multi-profile)
![GoVersion](https://img.shields.io/github/go-mod/go-version/bygui86/multi-profile)
![License](https://img.shields.io/github/license/bygui86/multi-profile)

[![BuildStatus](https://github.com/bygui86/multi-profile/workflows/build/badge.svg)](https://github.com/bygui86/multi-profile/actions)
![LatestRelease](https://img.shields.io/github/v/release/bygui86/multi-profile)

![LatestTag](https://img.shields.io/github/v/tag/bygui86/multi-profile)
![LastCommit](https://img.shields.io/github/last-commit/bygui86/multi-profile)
![Issues](https://img.shields.io/github/issues/bygui86/multi-profile)
![PullRequests](https://img.shields.io/github/issues-pr/bygui86/multi-profile)

Multi-profiling support package for Go.

This project was inspired by [pkg/profile](https://github.com/pkg/profile) but there is a fundamental difference: 
multi-profile offers the possibility to start multiple profiling at the same time.

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
    defer profile.CPUProfile(&profile.Config{}).Start().Stop()
    
    // ...
}
```

Using profile specific method, you can create the kind of profiling you want giving a Config as input. 

```go
package main

import "github.com/bygui86/multi-profile"

func main() {
    defer profile.CPUProfile(&profile.Config{}).Start().Stop()
    defer profile.MemProfile(&profile.Config{}).Start().Stop()
    defer profile.GoroutineProfile(&profile.Config{}).Start().Stop()

    // ...
}
```

`(i)️ INFO` see [examples](examples/) folder for all available profiles and samples.

`/!\ WARN` if not using `EnableInterruptHook` option (see below) ALWAYS remember to defer `Stop()` function, 
otherwise the profile won't stop and flush to file properly.

## Options

`(i)️️ INFO` see [examples](examples/) for all usage samples.

### Path

You can customize the path in which a profile file is going to be written.

Use field `Path` and `UseTempPath` in the Config.

### Interruption hook

You can enable an interruption hook that runs a new goroutine waiting for interruption signals (syscall.SIGTERM, 
syscall.SIGINT and os.Interrupt). If one of those signals arrives, the profiling packages stop the Profile and flushes 
results to file. Enabling this option, you can avoid deferring Stop() function in the main.

Use `EnableInterruptHook` field in the Config.

### Quiet mode

You can suppress all logs. 

Use field `Quiet` in the Config.

### Closer function

You can call a function right after stopping the profiling.

Use `CloserHook` field in the Config.

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

## Thanks

- [shields.io](https://shields.io) for providing badges
