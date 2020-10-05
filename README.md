# multi-profile - `WIP`

Multi-profiling support package for Go.

This project was inspired by [pkg/profile](github.com/pkg/profile) but there is a fundamental difference: multi-profile offers the possibility to start multiple profiling at the same time.

## TODO list

- [x] align with profile original (last align 1.10.2020)
- [ ] testing
    - [ ] profile-test
    - [ ] multi-test
- [ ] github actions
- [ ] README
    - [ ] github actions badge in readme, e.g. [![Build Status](https://travis-ci.org/pkg/profile.svg?branch=master)](https://travis-ci.org/pkg/profile)
    - [ ] godoc badge in readme, e.g. [![GoDoc](http://godoc.org/github.com/pkg/profile?status.svg)](http://godoc.org/github.com/pkg/profile)
- [ ] external logger

## Installation

```shell script
go get github.com/bygui86/multi-profile
```

## Usage - `WIP`

Enabling profiling in your application is as simple as one line at the top of your main function.

For example:

```go
import "github.com/bygui86/multi-profile"

func main() {
    defer profile.Start().Stop()
    ...
}
```

## Options - `WIP`

What to profile is controlled by config value passed to profile.Start. 
By default CPU profiling is enabled.

```go
import "github.com/bygui86/multi-profile"

func main() {
    // p.Stop() must be called before the program exits to
    // ensure profiling information is written to disk.
    p := profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.NoShutdownHook)
    ...
    // You can enable different kinds of memory profiling, either Heap or Allocs where Heap
    // profiling is the default with profile.MemProfile.
    p := profile.Start(profile.MemProfileAllocs, profile.ProfilePath("."), profile.NoShutdownHook)
}
```

Several convenience package level values are provided for cpu, memory, and block (contention) profiling.

For more complex options, consult the [documentation](http://godoc.org/github.com/bygui86/multi-profile).

## Contributing

I welcome pull requests, bug fixes and issue reports.

Before proposing a change, please discuss it first by raising an issue.
