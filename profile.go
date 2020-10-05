// Package profile provides a simple way to manage multiple runtime/pprof profiling of your Go application.
package profile

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sync/atomic"
	"syscall"
)

const (
	cpuMode = iota
	memMode
	mutexMode
	blockMode
	traceMode
	threadCreationMode
	goroutineMode

	// DefaultMemProfileRate is the default memory profiling rate.
	// See also http://golang.org/pkg/runtime/#pkg-variables
	DefaultMemProfileRate = 4096

	// DefaultMemProfileRate is the default memory profiling type.
	DefaultMemProfileType = MemProfileHeap

	MemProfileHeap   MemProfileType = "heap"
	MemProfileAllocs MemProfileType = "allocs"

	cpuPprofDefaultFilename       = "cpu.pprof"
	memPprofDefaultFilename       = "mem.pprof"
	mutexPprofDefaultFilename     = "mutex.pprof"
	blockPprofDefaultFilename     = "block.pprof"
	tracePprofDefaultFilename     = "trace.pprof"
	threadPprofDefaultFilename    = "threadcreate.pprof"
	goroutingPprofDefaultFilename = "goroutine.pprof"

	mutexPprof     = "mutex"
	blockPprof     = "block"
	threadPprof    = "threadcreate"
	goroutinePprof = "goroutine"
)

type MemProfileType string

// Profile represents a profiling session.
type Profile struct {
	// mode holds the type of profiling that will be made.
	mode int

	// path holds the base path where various profiling files are  written.
	// If blank, the base path will be generated by "ioutil.TempDir"
	path string

	// filename holds the filename created by starting the profiling
	filename string

	// disableShutdownHook controls whether the profiling package should hook SIGINT to write profiles cleanly.
	disableShutdownHook bool

	// quiet suppresses informational messages during profiling.
	quiet bool

	// memProfileRate holds the rate for the memory profiler.
	memProfileRate int

	// memProfileType holds the type for the memory profiler.
	memProfileType MemProfileType

	// closer holds the internal cleanup function that run after profiling Stop.
	closer func()

	// closerHook holds a custom cleanup function that run after profiling Stop.
	closerHook func()

	// started records if a call to profile.Start has already been made.
	started uint32
}

// ProfileConfig holds configurations to create a new Profile
type ProfileConfig struct {
	// Path holds the base path where various profiling files are  written.
	// If blank, the base path will be generated by "ioutil.TempDir"
	Path string

	// DisableShutdownHook controls whether the profiling package should hook SIGINT to write profiles cleanly.
	DisableShutdownHook bool

	// Quiet suppresses informational messages during profiling.
	Quiet bool

	// MemProfileRate holds the rate for the memory profiler.
	// See DefaultMemProfileRate for default value.
	MemProfileRate int

	// MemProfileType holds the type for the memory profiler.
	// Available values:   heap | allocs
	// See DefaultMemProfileType for default.
	MemProfileType MemProfileType

	// CloserHook holds a custom cleanup function that run after profiling Stop.
	CloserHook func()
}

// TODO validate ProfileConfig

// CPUProfile creates a CPU profiling object
func CPUProfile(cfg *ProfileConfig) *Profile {
	return &Profile{
		mode:                cpuMode,
		path:                cfg.Path,
		disableShutdownHook: cfg.DisableShutdownHook,
		quiet:               cfg.Quiet,
		closerHook:          cfg.CloserHook,
		started:             0,
	}
}

// MemProfile creates a memory profiling object
func MemProfile(cfg *ProfileConfig) *Profile {
	memRate := DefaultMemProfileRate
	memType := DefaultMemProfileType
	if cfg.MemProfileRate > 0 {
		memRate = cfg.MemProfileRate
	}
	if cfg.MemProfileType != "" {
		memType = cfg.MemProfileType
	}

	return &Profile{
		mode:                memMode,
		path:                cfg.Path,
		disableShutdownHook: cfg.DisableShutdownHook,
		quiet:               cfg.Quiet,
		memProfileRate:      memRate,
		memProfileType:      memType,
		closerHook:          cfg.CloserHook,
		started:             0,
	}
}

// MutexProfile creates a mutex profiling object
func MutexProfile(cfg *ProfileConfig) *Profile {
	return &Profile{
		mode:                mutexMode,
		path:                cfg.Path,
		disableShutdownHook: cfg.DisableShutdownHook,
		quiet:               cfg.Quiet,
		closerHook:          cfg.CloserHook,
		started:             0,
	}
}

// BlockProfile creates a block (contention) profiling object
func BlockProfile(cfg *ProfileConfig) *Profile {
	return &Profile{
		mode:                blockMode,
		path:                cfg.Path,
		disableShutdownHook: cfg.DisableShutdownHook,
		quiet:               cfg.Quiet,
		closerHook:          cfg.CloserHook,
		started:             0,
	}
}

// TraceProfile creates an execution tracing profiling object
func TraceProfile(cfg *ProfileConfig) *Profile {
	return &Profile{
		mode:                traceMode,
		path:                cfg.Path,
		disableShutdownHook: cfg.DisableShutdownHook,
		quiet:               cfg.Quiet,
		closerHook:          cfg.CloserHook,
		started:             0,
	}
}

// ThreadCreationProfile creates a thread creation profiling object
func ThreadCreationProfile(cfg *ProfileConfig) *Profile {
	return &Profile{
		mode:                threadCreationMode,
		path:                cfg.Path,
		disableShutdownHook: cfg.DisableShutdownHook,
		quiet:               cfg.Quiet,
		closerHook:          cfg.CloserHook,
		started:             0,
	}
}

// GoroutineProfile creates a goroutine profiling object
func GoroutineProfile(cfg *ProfileConfig) *Profile {
	return &Profile{
		mode:                goroutineMode,
		path:                cfg.Path,
		disableShutdownHook: cfg.DisableShutdownHook,
		quiet:               cfg.Quiet,
		closerHook:          cfg.CloserHook,
		started:             0,
	}
}

// Start starts a new profiling session.
func (p *Profile) Start() *Profile {
	if !atomic.CompareAndSwapUint32(&p.started, 0, 1) {
		// profile already started
		return p
	}

	p.preparePath()

	switch p.mode {
	case cpuMode:
		p.startCpuMode()

	case memMode:
		p.startMemMode()

	case mutexMode:
		p.startMutexMode()

	case blockMode:
		p.startBlockMode()

	case traceMode:
		p.startTraceMode()

	case threadCreationMode:
		p.startThreadCreationMode()

	case goroutineMode:
		p.startGoroutineMode()
	}

	p.startShutdownHook()
	return p
}

// Stop stops the profile and flushes any unwritten data.
// The caller should call the Stop method on the value returned to cleanly stop profiling.
func (p *Profile) Stop() {
	if !atomic.CompareAndSwapUint32(&p.started, 1, 0) {
		// profiling already stopped
		return
	}
	p.closer()
	p.closerHook()
}

// preparePath prepares the file path to flush data into when profiling will be stopped.
func (p *Profile) preparePath() {
	var pathErr error
	if p.path != "" {
		pathErr = os.MkdirAll(p.path, 0777)
	}
	p.path, pathErr = ioutil.TempDir("", "profile")
	if pathErr != nil {
		log.Fatalf("profiling start aborted, could not create initial output directory: %s", pathErr.Error())
	}
}

// startCpuMode start cpu profiling
func (p *Profile) startCpuMode() {
	p.filename = filepath.Join(p.path, cpuPprofDefaultFilename)
	file, fileErr := os.Create(p.filename)
	if fileErr != nil {
		log.Fatalf("could not create cpu profile %q: %s", p.filename, fileErr.Error())
	}
	p.logger("cpu profiling enabled (%q)", p.filename)
	startErr := pprof.StartCPUProfile(file)
	if startErr != nil {
		p.logger("cpu profiling start failed: %s", startErr.Error())
	}
	p.closer = p.stopAndFlush(file, -1)
}

// startMemMode starts memory profiling
func (p *Profile) startMemMode() {
	p.filename = filepath.Join(p.path, memPprofDefaultFilename)
	file, err := os.Create(p.filename)
	if err != nil {
		log.Fatalf("could not create memory profile %q: %s", p.filename, err.Error())
	}
	previous := runtime.MemProfileRate
	runtime.MemProfileRate = p.memProfileRate
	p.logger("memory profiling enabled at rate %d (%q)", runtime.MemProfileRate, p.filename)
	p.closer = p.stopAndFlush(file, previous)
}

// startMutexMode starts mutes profiling
func (p *Profile) startMutexMode() {
	p.filename = filepath.Join(p.path, mutexPprofDefaultFilename)
	file, err := os.Create(p.filename)
	if err != nil {
		log.Fatalf("could not create mutex profile %q: %s", p.filename, err.Error())
	}
	runtime.SetMutexProfileFraction(1)
	p.logger("mutex profiling enabled (%q)", p.filename)
	p.closer = p.stopAndFlush(file, -1)
}

// startBlockMode starts block profiling
func (p *Profile) startBlockMode() {
	p.filename = filepath.Join(p.path, blockPprofDefaultFilename)
	file, err := os.Create(p.filename)
	if err != nil {
		log.Fatalf("could not create block profile %q: %s", p.filename, err.Error())
	}
	runtime.SetBlockProfileRate(1)
	p.logger("block profiling enabled (%q)", p.filename)
	p.closer = p.stopAndFlush(file, -1)
}

// startTraceMode starts trace profiling
func (p *Profile) startTraceMode() {
	p.filename = filepath.Join(p.path, tracePprofDefaultFilename)
	file, fileErr := os.Create(p.filename)
	if fileErr != nil {
		log.Fatalf("could not create trace output file %q: %s", p.filename, fileErr.Error())
	}
	startErr := trace.Start(file)
	if startErr != nil {
		log.Fatalf("could not start trace: %s", startErr.Error())
	}
	p.logger("trace enabled (%q)", p.filename)
	p.closer = p.stopAndFlush(file, -1)
}

// startThreadCreationMode starts thread creation profiling
func (p *Profile) startThreadCreationMode() {
	p.filename = filepath.Join(p.path, threadPprofDefaultFilename)
	file, err := os.Create(p.filename)
	if err != nil {
		log.Fatalf("could not create thread creation profile %q: %s", p.filename, err.Error())
	}
	p.logger("thread creation profiling enabled (%q)", p.filename)
	p.closer = p.stopAndFlush(file, -1)
}

// startGoroutineMode starts goroutine profiling
func (p *Profile) startGoroutineMode() {
	p.filename = filepath.Join(p.path, goroutingPprofDefaultFilename)
	file, err := os.Create(p.filename)
	if err != nil {
		log.Fatalf("could not create goroutine profile %q: %s", p.filename, err.Error())
	}
	p.logger("goroutine profiling enabled (%q)", p.filename)
	p.closer = p.stopAndFlush(file, -1)
}

// stopAndFlush stops profiling and flush data to file.
func (p *Profile) stopAndFlush(file *os.File, previousMemRate int) func() {
	switch p.mode {

	case cpuMode:
		return func() {
			pprof.StopCPUProfile()
			err := file.Close()
			if err != nil {
				p.logger("cpu profiling error flushing data to file %q: %s", p.filename, err.Error())
			}
			p.logger("cpu profiling disabled (%q)", p.filename)
		}

	case memMode:
		return func() {
			pprofile := pprof.Lookup(string(p.memProfileType))
			if pprofile != nil {
				err := pprofile.WriteTo(file, 0)
				if err != nil {
					p.logger("memory profiling error flushing data to file %q: %s", p.filename, err.Error())
				}
			} else {
				p.logger("memory profiling error flushing data to file %q: pprof lookup returned null", p.filename)
			}
			err := file.Close()
			if err != nil {
				p.logger("memory profiling error flushing data to file %q: %s", p.filename, err.Error())
			}
			runtime.MemProfileRate = previousMemRate
			p.logger("memory profiling disabled (%q)", p.filename)
		}

	case mutexMode:
		return func() {
			pprofile := pprof.Lookup(mutexPprof)
			if pprofile != nil {
				err := pprofile.WriteTo(file, 0)
				if err != nil {
					p.logger("mutex profiling error flushing data to file %q: %s", p.filename, err.Error())
				}
			} else {
				p.logger("mutex profiling error flushing data to file %q: pprof lookup returned null", p.filename)
			}
			err := file.Close()
			if err != nil {
				p.logger("mutex profiling error flushing data to file %q: %s", p.filename, err.Error())
			}
			runtime.SetMutexProfileFraction(0)
			p.logger("mutex profiling disabled (%q)", p.filename)
		}

	case blockMode:
		return func() {
			pprofile := pprof.Lookup(blockPprof)
			if pprofile != nil {
				err := pprofile.WriteTo(file, 0)
				if err != nil {
					p.logger("block profiling error flushing data to file %q: %s", p.filename, err.Error())
				}
			} else {
				p.logger("block profiling error flushing data to file %q: pprof lookup returned null", p.filename)
			}
			err := file.Close()
			if err != nil {
				p.logger("block profiling error flushing data to file %q: %s", p.filename, err.Error())
			}
			runtime.SetBlockProfileRate(0)
			p.logger("block profiling disabled (%q)", p.filename)
		}

	case traceMode:
		return func() {
			trace.Stop()
			p.logger("trace disabled (%q)", p.filename)
		}

	case threadCreationMode:
		return func() {
			pprofile := pprof.Lookup(threadPprof)
			if pprofile != nil {
				err := pprofile.WriteTo(file, 0)
				if err != nil {
					p.logger("thread profiling error flushing data to file %q: %s", p.filename, err.Error())
				}
			} else {
				p.logger("thread profiling error flushing data to file %q: pprof lookup returned null", p.filename)
			}
			err := file.Close()
			if err != nil {
				p.logger("thread profiling error flushing data to file %q: %s", p.filename, err.Error())
			}
			p.logger("thread creation profiling disabled (%q)", p.filename)
		}

	case goroutineMode:
		return func() {
			pprofile := pprof.Lookup(goroutinePprof)
			if pprofile != nil {
				err := pprofile.WriteTo(file, 0)
				if err != nil {
					p.logger("goroutine profiling error flushing data to file %q: %s", p.filename, err.Error())
				}
			} else {
				p.logger("goroutine profiling error flushing data to file %q: pprof lookup returned null", p.filename)
			}
			err := file.Close()
			if err != nil {
				p.logger("goroutine profiling error flushing data to file %q: %s", p.filename, err.Error())
			}
			p.logger("goroutine profiling disabled (%q)", p.filename)
		}

	// WARN: we should never reach default!
	default:
		return func() {
			p.logger("unknown profiling disabled (%q)", p.filename)
		}
	}
}

// startShutdownHook starts a goroutine to wait for interruption signals and stop cleanly the profiling.
func (p *Profile) startShutdownHook() {
	if !p.disableShutdownHook {
		go func() {
			syscallCh := make(chan os.Signal)
			signal.Notify(syscallCh, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
			<-syscallCh

			log.Println("caught interrupt, stop profiling")
			p.Stop()

			os.Exit(0)
		}()
	}
}

func (p *Profile) logger(format string, args ...interface{}) {
	if !p.quiet {
		log.Printf(format, args...)
	}
}
