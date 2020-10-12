package examples

import (
	"fmt"
	"os"

	"github.com/bygui86/multi-profile/v2"
)

// Example to write profile to default path (same as application)
func DefaultPath() {
	cfg := &profile.Config{}
	prof := profile.CPUProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

// Example to write profile to specific path
func CustomPath() {
	cfg := &profile.Config{Path: os.Getenv("HOME")}
	prof := profile.CPUProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

// Example to write profile to a temporary path
func TempDirPath() {
	cfg := &profile.Config{UseTempPath: true}
	prof := profile.CPUProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

// Example to enable the automatic shutdown hook based on interruption signals
func EnableInterruptHook() {
	cfg := &profile.Config{EnableInterruptHook: true}
	prof := profile.CPUProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

// Example to disable logs
func Quiet() {
	cfg := &profile.Config{Quiet: true}
	prof := profile.CPUProfile(cfg)
	prof.Start()
	defer prof.Stop()
}

// Example with a custom closer function
func CustomCloser() {
	cfg := &profile.Config{
		CloserHook: func() {
			fmt.Println("This is the custom closer executed after profile Stop")
		},
	}
	prof := profile.CPUProfile(cfg)
	prof.Start()
	defer prof.Stop()
}
