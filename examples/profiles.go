package examples

import (
	"github.com/bygui86/multi-profile/v2"
)

func OneLine() {
	defer profile.CPUProfile(&profile.Config{}).Start().Stop()
}

func Multiple() {
	defer profile.CPUProfile(&profile.Config{}).Start().Stop()
	defer profile.MemProfile(&profile.Config{}).Start().Stop()
	defer profile.GoroutineProfile(&profile.Config{}).Start().Stop()
}

func CPUProfile() {
	cfg := &profile.Config{}
	prof := profile.CPUProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

// Default to type heap
func MemProfile() {
	cfg := &profile.Config{}
	prof := profile.MemProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

// Default to type heap
func MemProfileRate() {
	cfg := &profile.Config{
		MemProfileRate: 1024,
	}
	prof := profile.MemProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func MemProfileAllocs() {
	cfg := &profile.Config{
		MemProfileType: profile.MemProfileAllocs,
	}
	prof := profile.MemProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func MemProfileAllocsRate() {
	cfg := &profile.Config{
		MemProfileRate: 1024,
		MemProfileType: profile.MemProfileAllocs,
	}
	prof := profile.MemProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func MutexProfile() {
	cfg := &profile.Config{}
	prof := profile.MutexProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func BlockProfile() {
	cfg := &profile.Config{}
	prof := profile.BlockProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func TraceProfile() {
	cfg := &profile.Config{}
	prof := profile.TraceProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func ThreadCreationProfile() {
	cfg := &profile.Config{}
	prof := profile.ThreadCreationProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func GoroutineProfile() {
	cfg := &profile.Config{}
	prof := profile.GoroutineProfile(cfg)
	prof.Start()
	defer prof.Stop()
}
