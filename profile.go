// Package profile provides a simple way to manage multiple runtime/pprof profiling of your Go application
package profile

import (
	"fmt"
	"io/ioutil"
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
	// Supported profiles
	cpuMode       profileMode = "CPU"
	memMode       profileMode = "Memory"
	mutexMode     profileMode = "Mutex"
	blockMode     profileMode = "Block"
	traceMode     profileMode = "Trace"
	threadMode    profileMode = "Thread"
	goroutineMode profileMode = "Goroutine"

	// DefaultPath holds the default path where to create pprof file
	DefaultPath = "./"

	/*
		DefaultMemProfileRate holds the default memory profiling rate
		See also http://golang.org/pkg/runtime/#pkg-variables
	*/
	DefaultMemProfileRate = 4096

	// DefaultMemProfileRate holds the default memory profiling type
	DefaultMemProfileType = MemProfileHeap

	// Supported memory profiles
	MemProfileHeap   MemProfileType = "heap"
	MemProfileAllocs MemProfileType = "allocs"

	// Supported logging level
	debugLevel logLevel = "debug"
	infoLevel  logLevel = "info"
	warnLevel  logLevel = "warn"
	errorLevel logLevel = "error"
	fatalLevel logLevel = "fatal"
)

// Profile represents a profiling session
type Profile struct {
	// mode holds the type of profiling that will be made
	mode profileMode

	// lookupName holds the lookup name used to stop and flush the profile using pprof package
	lookupName string

	/*
		path holds the base path where various profiling files will be written.
		If blank, the base path will be the current directory "./"
	*/
	path string

	// useTempPath let the path be generated by "ioutil.TempDir"
	useTempPath bool

	// fileName holds the name of the file created by the profile
	fileName string

	// filePath holds the path to the file created by the profile
	filePath string

	// file holds the reference to the file created by the profile
	file *os.File

	// panicIfFail holds the flag to decide whether a profile failure causes a panic
	panicIfFail bool

	// enableInterruptHook controls whether to start a goroutine to wait for interruption signals to stop profiling
	enableInterruptHook bool

	// quiet suppresses informational messages during profiling
	quiet bool

	/*
		memProfileRate holds the rate for the memory profile
		See DefaultMemProfileRate for default value
	*/
	memProfileRate int

	/*
		memProfileType holds the type for the memory profile
		Available values:   heap | allocs
		See DefaultMemProfileType for default
	*/
	memProfileType MemProfileType

	/*
		internalCloser holds the internal cleanup function that run after profiling Stop
		This function is specific for each profile (CPU, MEM, GoRoutines, etc)
	*/
	internalCloser func()

	// closerHook holds a custom cleanup function that run after profiling Stop
	closerHook func()

	// Logger offers the possibility to inject a custom logger
	logger Logger

	// previousMemProfileRate keeps track of the previous runtime.MemProfileRate value
	previousMemProfileRate int

	// started records if a call to profile.Start has already been made
	started uint32
}

// Config holds configurations to create a new Profile
type Config struct {
	/*
		Path holds the base path where various profiling files are  written
		If blank, the base path will be the current directory "./"
	*/
	Path string

	// UseTempPath let the path be generated by "ioutil.TempDir"
	UseTempPath bool

	// PanicIfFail holds the flag to decide whether a profile failure causes a panic
	PanicIfFail bool

	// EnableInterruptHook controls whether to start a goroutine to wait for interruption signals to stop profiling
	EnableInterruptHook bool

	// Quiet suppresses informational messages during profiling
	Quiet bool

	/*
		MemProfileRate holds the rate for the memory profile
		See DefaultMemProfileRate for default value
	*/
	MemProfileRate int

	/*
		MemProfileType holds the type for the memory profile
		Available values:   heap | allocs
		See DefaultMemProfileType for default
	*/
	MemProfileType MemProfileType

	// CloserHook holds a custom cleanup function that run after profiling Stop
	CloserHook func()

	// Logger offers the possibility to inject a custom logger
	Logger Logger
}

// MemProfileType defines which type of memory profiling you want to start
type MemProfileType string

// profileMode defined which profiling mode has to be run
type profileMode string

// logLevel defines the level at which a message has to be logged
type logLevel string

// Logger defines the interface an external logger have to implement, to be passed and used by multi-profile
type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})

	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}

// CPUProfile creates a CPU profiling object
func CPUProfile(cfg *Config) *Profile {
	// INFO: lookupName not required
	return buildProfile(cpuMode, "", "cpu.pprof", cfg)
}

// MemProfile creates a memory profiling object
func MemProfile(cfg *Config) *Profile {
	memRate := DefaultMemProfileRate
	memType := DefaultMemProfileType
	if cfg.MemProfileRate > 0 {
		memRate = cfg.MemProfileRate
	}
	if cfg.MemProfileType != "" {
		memType = cfg.MemProfileType
	}

	memPprof := buildProfile(memMode, string(memType), "mem.pprof", cfg)
	memPprof.memProfileRate = memRate
	memPprof.memProfileType = memType
	return memPprof
}

// MutexProfile creates a mutex profiling object
func MutexProfile(cfg *Config) *Profile {
	return buildProfile(mutexMode, "mutex", "mutex.pprof", cfg)
}

// BlockProfile creates a block (contention) profiling object
func BlockProfile(cfg *Config) *Profile {
	return buildProfile(blockMode, "block", "block.pprof", cfg)
}

// TraceProfile creates an execution tracing profiling object
func TraceProfile(cfg *Config) *Profile {
	// INFO: lookupName not required
	return buildProfile(traceMode, "", "trace.pprof", cfg)
}

// ThreadCreationProfile creates a thread creation profiling object
func ThreadCreationProfile(cfg *Config) *Profile {
	return buildProfile(threadMode, "thread", "thread.pprof", cfg)
}

// GoroutineProfile creates a goroutine profiling object
func GoroutineProfile(cfg *Config) *Profile {
	return buildProfile(goroutineMode, "goroutine", "goroutine.pprof", cfg)
}

// Start starts a new profiling session
func (p *Profile) Start() *Profile {
	if !atomic.CompareAndSwapUint32(&p.started, 0, 1) {
		// no-op, profiling already started
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

	case threadMode:
		p.startThreadCreationMode()

	case goroutineMode:
		p.startGoroutineMode()
	}

	p.startInterruptHook()

	return p
}

/*
	Stop stops the profiling and flushes any unwritten data.
	The caller should call the Stop method on the value returned to cleanly stop profiling.
*/
func (p *Profile) Stop() {
	if !atomic.CompareAndSwapUint32(&p.started, 1, 0) {
		// no-op, profiling already stopped
		return
	}

	if p.internalCloser != nil {
		p.internalCloser()
	}

	if p.closerHook != nil {
		p.closerHook()
	}
}

// startCpuMode starts cpu profiling
func (p *Profile) startCpuMode() {
	p.createFile()

	err := pprof.StartCPUProfile(p.file)
	if err != nil {
		p.logf(errorLevel, "CPU profiling start failed: %s", err.Error())
		if p.panicIfFail {
			panic(err)
		}
	}

	p.internalCloser = p.stopCpuMode

	p.logf(infoLevel, "CPU profiling enabled, file %s", p.filePath)
}

// startMemMode starts memory profiling
func (p *Profile) startMemMode() {
	p.createFile()

	p.previousMemProfileRate = runtime.MemProfileRate
	runtime.MemProfileRate = p.memProfileRate
	p.internalCloser = p.stopMemMode

	p.logf(infoLevel, "Memory profiling (%s) enabled at rate %d, file %s",
		p.memProfileType, runtime.MemProfileRate, p.filePath)
}

// startMutexMode starts mutes profiling
func (p *Profile) startMutexMode() {
	p.createFile()

	runtime.SetMutexProfileFraction(1)
	p.internalCloser = p.stopMutexMode

	p.logf(infoLevel, "Mutex profiling enabled, file %s", p.filePath)
}

// startBlockMode starts block profiling
func (p *Profile) startBlockMode() {
	p.createFile()

	runtime.SetBlockProfileRate(1)
	p.internalCloser = p.stopBlockMode

	p.logf(infoLevel, "Block profiling enabled, file %s", p.filePath)
}

// startTraceMode starts trace profiling
func (p *Profile) startTraceMode() {
	p.createFile()

	err := trace.Start(p.file)
	if err != nil {
		p.logf(errorLevel, "Trace profiling start failed: %s", err.Error())
		if p.panicIfFail {
			panic(err)
		}
	}

	p.internalCloser = p.stopTraceMode

	p.logf(infoLevel, "Trace profiling enabled, file %s", p.filePath)
}

// startThreadCreationMode starts thread creation profiling
func (p *Profile) startThreadCreationMode() {
	p.createFile()

	p.internalCloser = p.stopThreadCreationMode

	p.logf(infoLevel, "Thread profiling enabled, file %s", p.filePath)
}

// startGoroutineMode starts goroutine profiling
func (p *Profile) startGoroutineMode() {
	p.createFile()

	p.internalCloser = p.stopGoroutineMode

	p.logf(infoLevel, "Goroutine profiling enabled, file %s", p.filePath)
}

// stopCpuMode stops cpu profiling
func (p *Profile) stopCpuMode() {
	p.logf(infoLevel, "Stop and flush CPU profiling to file %s", p.filePath)

	pprof.StopCPUProfile()
	err := p.file.Close()
	if err != nil {
		p.logf(errorLevel, "CPU profiling flushing data to file %q failed: %s", p.filePath, err.Error())
	}

	p.log(infoLevel, "CPU profiling disabled")
}

// stopMemMode stops memory profiling
func (p *Profile) stopMemMode() {
	p.stopAndFlush()

	runtime.MemProfileRate = p.previousMemProfileRate
	p.previousMemProfileRate = -1
}

// stopMutexMode stops mutex profiling
func (p *Profile) stopMutexMode() {
	p.stopAndFlush()

	runtime.SetMutexProfileFraction(0)
}

// stopBlockMode stops block profiling
func (p *Profile) stopBlockMode() {
	p.stopAndFlush()

	runtime.SetBlockProfileRate(0)
}

// stopTraceMode stops trace profiling
func (p *Profile) stopTraceMode() {
	p.logf(infoLevel, "Stop and flush trace profiling to file %s", p.filePath)

	trace.Stop()

	p.log(infoLevel, "Trace profiling disabled")
}

// stopThreadCreationMode stops thread creation profiling
func (p *Profile) stopThreadCreationMode() {
	p.stopAndFlush()
}

// stopGoroutineMode stops goroutine profiling
func (p *Profile) stopGoroutineMode() {
	p.stopAndFlush()
}

// startInterruptHook starts the interruptHook function in a separate goroutine
func (p *Profile) startInterruptHook() {
	if p.enableInterruptHook {
		p.logf(infoLevel, "Start interrupt hook for %s profiling", string(p.mode))
		go p.interruptHook()
	}
}

// interruptHook waits for interruption signals and stop the profiling
func (p *Profile) interruptHook() {
	syscallCh := make(chan os.Signal)
	signal.Notify(syscallCh, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-syscallCh

	p.logf(warnLevel, "Caught interrupt signal, stop and flush %s profiling to file", string(p.mode))
	p.Stop()
}

// buildProfile builds a Profile using input parameters
func buildProfile(mode profileMode, lookupName, fileName string, cfg *Config) *Profile {
	return &Profile{
		mode:                mode,
		lookupName:          lookupName,
		path:                cfg.Path,
		useTempPath:         cfg.UseTempPath,
		fileName:            fileName,
		panicIfFail:         cfg.PanicIfFail,
		enableInterruptHook: cfg.EnableInterruptHook,
		quiet:               cfg.Quiet,
		logger:              cfg.Logger,
		closerHook:          cfg.CloserHook,
		started:             0,
	}
}

// createFile creates the file that the profile will use to flush results into
func (p *Profile) createFile() {
	p.filePath = filepath.Join(p.path, p.fileName)
	var err error
	p.file, err = os.Create(p.filePath)
	if err != nil {
		p.logf(errorLevel, "%s profiling file %s creation failed: %s",
			string(p.mode), p.filePath, err.Error())
		if p.panicIfFail {
			panic(err)
		}
	}
}

// stopAndFlush stops profiling and flushes results to file (valid for all modes except CPU and Trace)
func (p *Profile) stopAndFlush() {
	p.logf(infoLevel, "Stop and flush %s lookup for %s profiling to file %s", p.lookupName, string(p.mode), p.filePath)
	pprofile := pprof.Lookup(p.lookupName)
	if pprofile != nil {
		err := pprofile.WriteTo(p.file, 0)
		if err != nil {
			p.logf(errorLevel, "%s profiling flushing data to file %s failed: %s",
				string(p.mode), p.filePath, err.Error())
		}
	} else {
		p.logf(errorLevel, "%s profiling flushing data to file %s failed: pprof lookup returned nil profile",
			string(p.mode), p.filePath)
	}

	err := p.file.Close()
	if err != nil {
		p.logf(errorLevel, "%s profiling flushing data to file %s failed: %s",
			string(p.mode), p.filePath, err.Error())
	}

	p.logf(infoLevel, "%s profiling disabled", string(p.mode))
}

// preparePath prepares the file path to flush data into when profiling will be stopped
func (p *Profile) preparePath() {
	var err error
	if p.useTempPath {
		err = p.prepareTempPath()
	} else {
		err = p.prepareCustomPath()
	}
	if err != nil {
		p.logf(errorLevel, "%s profiling start aborted, could not create output directory: %s",
			string(p.mode), err.Error())
		if p.panicIfFail {
			panic(err)
		} else {
			if p.panicIfFail {
				panic(err)
			}
		}
	}
}

// prepareCustomPath prepares the file path using default or custom parameters to flush data into when profiling will be stopped
func (p *Profile) prepareCustomPath() error {
	if p.path == "" {
		p.path = DefaultPath
	}

	if p.path != DefaultPath {
		mkdirErr := os.MkdirAll(p.path, 0777)
		if mkdirErr != nil {
			return mkdirErr
		}
	}

	return nil
}

// prepareTempPath prepares the file path in 'tmp' folder to flush data into when profiling will be stopped
func (p *Profile) prepareTempPath() error {
	var err error
	p.path, err = ioutil.TempDir("", "profile_")
	if err != nil {
		return err
	}
	return nil
}

// log abstracts the complexity of using an external specific logger
func (p *Profile) log(level logLevel, args ...interface{}) {
	if !p.quiet {
		if p.logger != nil {
			switch level {
			case debugLevel:
				p.logger.Debug(args...)
			case infoLevel:
				p.logger.Info(args...)
			case warnLevel:
				p.logger.Warn(args...)
			case errorLevel:
				p.logger.Error(args...)
			case fatalLevel:
				p.logger.Fatal(args...)
			default:
				p.logger.Info(args...)
			}
		} else {
			fmt.Print(fmt.Sprintf("[%s]", level), args, "\n")
		}
	}
}

// logf abstracts the complexity of using an external specific logger
func (p *Profile) logf(level logLevel, template string, args ...interface{}) {
	if !p.quiet {
		if p.logger != nil {
			switch level {
			case debugLevel:
				p.logger.Debugf(template, args...)
			case infoLevel:
				p.logger.Infof(template, args...)
			case warnLevel:
				p.logger.Warnf(template, args...)
			case errorLevel:
				p.logger.Errorf(template, args...)
			case fatalLevel:
				p.logger.Fatalf(template, args...)
			default:
				p.logger.Infof(template, args...)
			}
		} else {
			fmt.Printf("[%s] %s\n", level, fmt.Sprintf(template, args...))
		}
	}
}
