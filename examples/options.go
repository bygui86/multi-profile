package examples

import (
	"fmt"
	"os"

	"github.com/bygui86/multi-profile"
)

// Example to write profile to default path (same as application)
func DefaultPath() {
	cfg := &profile.ProfileConfig{}
	prof := profile.CPUProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

// Example to write profile to specific path
func CustomPath() {
	cfg := &profile.ProfileConfig{Path: os.Getenv("HOME")}
	prof := profile.CPUProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

// Example to write profile to a temporary path
func TempDirPath() {
	cfg := &profile.ProfileConfig{UseTempPath: true}
	prof := profile.CPUProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

// Example to disable the automatic shutdown hook
func DisableShutdownHook() {
	cfg := &profile.ProfileConfig{DisableShutdownHook: true}
	prof := profile.CPUProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

// Example to disable logs
func Quiet() {
	cfg := &profile.ProfileConfig{Quiet: true}
	prof := profile.CPUProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

// Example with a custom closer function
func CustomCloser() {
	cfg := &profile.ProfileConfig{
		CloserHook: func() {
			fmt.Println("This is the custom closer executed after profile Stop")
		},
	}
	prof := profile.CPUProfile(cfg)
	prof.Start()
	defer prof.Stop()
}
