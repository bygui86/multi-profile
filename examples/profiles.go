package examples

import (
	"github.com/bygui86/multi-profile"
)

func OneLine() {
	defer profile.CPUProfile(&profile.ProfileConfig{}).Start().Stop()
}

func Multiple() {
	defer profile.CPUProfile(&profile.ProfileConfig{}).Start().Stop()
	defer profile.MemProfile(&profile.ProfileConfig{}).Start().Stop()
	defer profile.GoroutineProfile(&profile.ProfileConfig{}).Start().Stop()
}

func CPUProfile() {
	cfg := &profile.ProfileConfig{}
	prof := profile.CPUProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

// Default to type heap
func MemProfile() {
	cfg := &profile.ProfileConfig{}
	prof := profile.MemProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

// Default to type heap
func MemProfileRate() {
	cfg := &profile.ProfileConfig{
		MemProfileRate: 1024,
	}
	prof := profile.MemProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func MemProfileAllocs() {
	cfg := &profile.ProfileConfig{
		MemProfileType: profile.MemProfileAllocs,
	}
	prof := profile.MemProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func MemProfileAllocsRate() {
	cfg := &profile.ProfileConfig{
		MemProfileRate: 1024,
		MemProfileType: profile.MemProfileAllocs,
	}
	prof := profile.MemProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func MutexProfile() {
	cfg := &profile.ProfileConfig{}
	prof := profile.MutexProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func BlockProfile() {
	cfg := &profile.ProfileConfig{}
	prof := profile.BlockProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func TraceProfile() {
	cfg := &profile.ProfileConfig{}
	prof := profile.TraceProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func ThreadCreationProfile() {
	cfg := &profile.ProfileConfig{}
	prof := profile.ThreadCreationProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

func GoroutineProfile() {
	cfg := &profile.ProfileConfig{}
	prof := profile.GoroutineProfile(cfg)
	prof.Start()
	defer prof.Stop()
}
