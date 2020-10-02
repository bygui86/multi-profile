package profile

import (
	"os"
)

func ExampleOneLine() {
	defer CPUProfile(&ProfileConfig{}).Start().Stop()
}

func ExampleCPUProfile() {
	cfg := &ProfileConfig{}
	prof := CPUProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

// Default to type heap
func ExampleMemProfile() {
	cfg := &ProfileConfig{}
	prof := MemProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

// Default to type heap
func ExampleMemProfileRate() {
	cfg := &ProfileConfig{
		MemProfileRate: 1024,
	}
	prof := MemProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func ExampleMemProfileAllocs() {
	cfg := &ProfileConfig{
		MemProfileType: MemProfileAllocs,
	}
	prof := MemProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func ExampleMemProfileAllocsRate() {
	cfg := &ProfileConfig{
		MemProfileRate: 1024,
		MemProfileType: MemProfileAllocs,
	}
	prof := MemProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func ExampleMutexProfile() {
	cfg := &ProfileConfig{}
	prof := MutexProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func ExampleBlockProfile() {
	cfg := &ProfileConfig{}
	prof := BlockProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func ExampleTraceProfile() {
	cfg := &ProfileConfig{}
	prof := TraceProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func ExampleThreadCreationProfile() {
	cfg := &ProfileConfig{}
	prof := ThreadCreationProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func ExampleGoroutineProfile() {
	cfg := &ProfileConfig{}
	prof := GoroutineProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

// set the location that the profile will be written to
func ExampleProfilePath() {
	cfg := &ProfileConfig{Path: os.Getenv("HOME")}
	prof := CPUProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

// disable the automatic shutdown hook
func ExampleDisableShutdownHook() {
	cfg := &ProfileConfig{DisableShutdownHook: true}
	prof := CPUProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

// disable logs
func ExampleQuiet() {
	cfg := &ProfileConfig{Quiet: true}
	prof := CPUProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func ExampleMultipleProfile() {
	cpuCfg := &ProfileConfig{Path: "."}
	cpuProf := CPUProfile(cpuCfg)
	cpuProf.Start()
	defer cpuProf.Stop()

	memCfg := &ProfileConfig{Path: "."}
	memProf := MemProfile(memCfg)
	memProf.Start()
	defer memProf.Stop()

	cfg := &ProfileConfig{}
	prof := GoroutineProfile(cfg)
	prof.Start()
	defer prof.Stop()
}
