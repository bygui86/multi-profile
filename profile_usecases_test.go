package profile_test

import (
	"os"
)

var profileTests = []profileTest{
	{
		name: "cpu profile",
		code: `
			package main
			
			import "github.com/bygui86/multi-profile"
			
			func main() {
				defer profile.CPUProfile(&profile.Config{}).Start().Stop()
			}
			`,
		checks: []checkFn{
			Stdout("cpu profiling enabled", "cpu profiling disabled"),
			NotInStdout("panic situation recovered"),
			NoStderr,
			NoErr,
		},
	},
	{
		name: "heap memory profile",
		code: `
			package main
	
			import "github.com/bygui86/multi-profile"
	
			func main() {
				defer profile.MemProfile(&profile.Config{}).Start().Stop()
			}
			`,
		checks: []checkFn{
			Stdout("memory profiling (heap) enabled", "memory profiling (heap) disabled"),
			NotInStdout("panic situation recovered"),
			NoStderr,
			NoErr,
		},
	},
	{
		name: "allocs memory profile",
		code: `
			package main
	
			import "github.com/bygui86/multi-profile"
	
			func main() {
				defer profile.MemProfile(&profile.Config{MemProfileType: profile.MemProfileAllocs}).Start().Stop()
			}
			`,
		checks: []checkFn{
			Stdout("memory profiling (allocs) enabled", "memory profiling (allocs) disabled"),
			NotInStdout("panic situation recovered"),
			NoStderr,
			NoErr,
		},
	},
	{
		name: "rate memory profile",
		code: `
			package main
	
			import "github.com/bygui86/multi-profile"
	
			func main() {
				defer profile.MemProfile(&profile.Config{MemProfileRate: 1024}).Start().Stop()
			}
			`,
		checks: []checkFn{
			Stdout("memory profiling (heap) enabled at rate 1024", "memory profiling (heap) disabled"),
			NotInStdout("panic situation recovered"),
			NoStderr,
			NoErr,
		},
	},
	{
		name: "mutex profile",
		code: `
			package main
	
			import "github.com/bygui86/multi-profile"
	
			func main() {
				defer profile.MutexProfile(&profile.Config{}).Start().Stop()
			}
			`,
		checks: []checkFn{
			Stdout("mutex profiling enabled", "mutex profiling disabled"),
			NotInStdout("panic situation recovered"),
			NoStderr,
			NoErr,
		},
	},
	{
		name: "block profile",
		code: `
			package main
	
			import "github.com/bygui86/multi-profile"
	
			func main() {
				defer profile.BlockProfile(&profile.Config{}).Start().Stop()
			}
			`,
		checks: []checkFn{
			Stdout("block profiling enabled", "block profiling disabled"),
			NotInStdout("panic situation recovered"),
			NoStderr,
			NoErr,
		},
	},
	{
		name: "trace profile",
		code: `
			package main
	
			import "github.com/bygui86/multi-profile"
	
			func main() {
				defer profile.TraceProfile(&profile.Config{}).Start().Stop()
			}
			`,
		checks: []checkFn{
			Stdout("trace profiling enabled", "trace profiling disabled"),
			NotInStdout("panic situation recovered"),
			NoStderr,
			NoErr,
		},
	},
	{
		name: "thread creation profile",
		code: `
			package main
	
			import "github.com/bygui86/multi-profile"
	
			func main() {
				defer profile.ThreadCreationProfile(&profile.Config{}).Start().Stop()
			}
			`,
		checks: []checkFn{
			Stdout("thread creation profiling enabled", "thread creation profiling disabled"),
			NotInStdout("panic situation recovered"),
			NoStderr,
			NoErr,
		},
	},
	{
		name: "goroutine profile",
		code: `
			package main
	
			import "github.com/bygui86/multi-profile"
	
			func main() {
				defer profile.GoroutineProfile(&profile.Config{}).Start().Stop()
			}
			`,
		checks: []checkFn{
			Stdout("goroutine profiling enabled", "goroutine profiling disabled"),
			NotInStdout("panic situation recovered"),
			NoStderr,
			NoErr,
		},
	},
	{
		name: "multi profile",
		code: `
			package main
			
			import "github.com/bygui86/multi-profile"
			
			func main() {
				defer profile.CPUProfile(&profile.Config{}).Start().Stop()
				defer profile.MemProfile(&profile.Config{}).Start().Stop()
			}
			`,
		checks: []checkFn{
			Stdout("cpu profiling enabled", "cpu profiling disabled",
				"memory profiling (heap) enabled", "memory profiling (heap) disabled"),
			NotInStdout("panic situation recovered"),
			NoStderr,
			NoErr,
		},
	},
	{
		name: "profile panic",
		code: `
			package main
	
			import "github.com/bygui86/multi-profile"
	
			func main() {
				defer profile.CPUProfile(&profile.Config{Path: "/private"}).Start().Stop()
			}
			`,
		checks: []checkFn{
			Stdout("permission denied"),
			Err,
		},
	},
}

var optionsTests = []profileTest{
	{
		name: "quiet option",
		code: `
			package main
	
			import "github.com/bygui86/multi-profile"
	
			func main() {
				defer profile.CPUProfile(&profile.Config{Quiet: true}).Start().Stop()
			}
			`,
		checks: []checkFn{
			NoStdout,
			NoStderr,
			NoErr,
		},
	},
	{
		name: "custom path option",
		code: `
			package main
	
			import (
				"os"
				"github.com/bygui86/multi-profile"
			)
	
			func main() {
				defer profile.CPUProfile(&profile.Config{Path: os.Getenv("HOME")}).Start().Stop()
			}
			`,
		checks: []checkFn{
			Stdout("cpu profiling enabled", "cpu profiling disabled", os.Getenv("HOME")+"/cpu.pprof"),
			NotInStdout("panic situation recovered"),
			NoStderr,
			NoErr,
		},
	},
	{
		name: "temp path option",
		code: `
			package main
	
			import "github.com/bygui86/multi-profile"
	
			func main() {
				defer profile.CPUProfile(&profile.Config{UseTempPath: true}).Start().Stop()
			}
			`,
		checks: []checkFn{
			Stdout("cpu profiling enabled", "cpu profiling disabled", "profile_"),
			NotInStdout("panic situation recovered"),
			NoStderr,
			NoErr,
		},
	},
	{
		name: "custom closer option",
		code: `
			package main
	
			import (
				"fmt"
				"github.com/bygui86/multi-profile"
			)
	
			func main() {
				defer profile.CPUProfile(&profile.Config{CloserHook: closerFn}).Start().Stop()
			}
	
			func closerFn() {
				fmt.Println("custom closer")
			}
			`,
		checks: []checkFn{
			Stdout("cpu profiling enabled", "cpu profiling disabled", "custom closer"),
			NotInStdout("panic situation recovered"),
			NoErr,
		},
	},
	{
		name: "custom path error",
		code: `
			package main
	
			import "github.com/bygui86/multi-profile"
	
			func main() {
				defer profile.CPUProfile(&profile.Config{Path: "/private"}).Start().Stop()
			}
			`,
		checks: []checkFn{
			Stdout("permission denied"),
			NotInStdout("panic situation recovered"),
			Stderr("exit status"),
			Err,
		},
	},
	{
		name: "no exit",
		code: `
			package main
	
			import "github.com/bygui86/multi-profile"
	
			func main() {
				defer profile.CPUProfile(&profile.Config{Path: "/private", NoExit: true}).Start().Stop()
			}
			`,
		checks: []checkFn{
			Stdout("permission denied", "panic situation recovered"),
			Err,
		},
	},
}
