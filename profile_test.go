package profile_test

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfiles(t *testing.T) {
	for _, profTest := range profileTests {
		t.Logf("Run profile test '%s'", profTest.name)
		stdout, stderr, err := runTest(t, profTest.code)
		for _, check := range profTest.checks {
			check(t, stdout, stderr, err)
		}
	}

	checkPprofFiles(t, []string{
		"./cpu.pprof", "./mem.pprof", "./mutex.pprof", "./block.pprof",
		"./trace.pprof", "./thread.pprof", "./goroutine.pprof",
	})

	cleanupPprofFiles(t, []string{
		"./cpu.pprof", "./mem.pprof", "./mutex.pprof", "./block.pprof",
		"./trace.pprof", "./thread.pprof", "./goroutine.pprof",
	})
}

func TestOptions(t *testing.T) {
	for _, profTest := range optionsTests {
		t.Logf("Run option test '%s'", profTest.name)
		stdout, stderr, err := runTest(t, profTest.code)
		for _, check := range profTest.checks {
			check(t, stdout, stderr, err)
		}
	}

	checkPprofFiles(t, []string{
		"./cpu.pprof", os.Getenv("HOME") + "/cpu.pprof",
	})

	cleanupPprofFiles(t, []string{
		"./cpu.pprof", os.Getenv("HOME") + "/cpu.pprof",
	})
}

type profileTest struct {
	name   string
	code   string
	checks []checkFn
}

type checkFn func(t *testing.T, stdout, stderr []byte, err error)

// Stdout verifies that the given lines match the output from stdout
func Stdout(lines ...string) checkFn {
	return func(t *testing.T, stdout, stderr []byte, err error) {
		r := bytes.NewReader(stdout)
		if !validateOutput(r, lines) {
			t.Errorf("stdout: expected '%s', actual '%s'", strings.Join(lines, ", "), stdout)
		}
	}
}

// NotInStdout verifies that the given lines do not match the output from stdout
func NotInStdout(lines ...string) checkFn {
	return func(t *testing.T, stdout, stderr []byte, err error) {
		r := bytes.NewReader(stdout)
		if validateOutput(r, lines) {
			t.Errorf("stdout: '%s' was not expected, but found in stdout '%s'", strings.Join(lines, ", "), stdout)
		}
	}
}

// NoStdout checks that stdout was blank
func NoStdout(t *testing.T, stdout, stderr []byte, err error) {
	if len(stdout) > 0 {
		t.Errorf("stdout: expected 0 bytes, actual %d bytes - bytes to string: '%s'", len(stdout), string(stdout))
	}
}

// Stderr verifies that the given lines match the output from stderr
func Stderr(lines ...string) checkFn {
	return func(t *testing.T, stdout, stderr []byte, err error) {
		r := bytes.NewReader(stderr)
		if !validateOutput(r, lines) {
			t.Errorf("stderr: expected '%s', actual '%s'", strings.Join(lines, ", "), stderr)
		}
	}
}

// NoStderr checks that stderr was blank
func NoStderr(t *testing.T, stdout, stderr []byte, err error) {
	if len(stderr) > 0 {
		t.Errorf("stderr: expected 0 bytes, actual %d bytes - bytes to string: '%s'", len(stderr), string(stderr))
	}
}

// Err checks that there was an error returned
func Err(t *testing.T, stdout, stderr []byte, err error) {
	if err == nil {
		t.Errorf("expected error")
	}
}

// NoErr checks that err was nil
func NoErr(t *testing.T, stdout, stderr []byte, err error) {
	if err != nil {
		t.Errorf("error: expected nil, actual '%v'", err)
	}
}

// validateOutput checks if the given slice of lines are among data from the given reader
func validateOutput(reader io.Reader, expected []string) bool {
	scanner := bufio.NewScanner(reader)
	for _, line := range expected {
		flag := true
		for scanner.Scan() {
			if strings.Contains(strings.ToLower(scanner.Text()), strings.ToLower(line)) {
				return true
			} else {
				flag = false
			}
		}
		return flag
	}
	return true
}

/*
	runTest executes the go program supplied and returns the contents of stdout,
	stderr and an error which may contain status information about the result of the execution.
*/
func runTest(t *testing.T, codeToTest string) ([]byte, []byte, error) {
	tempGopathDir, goPathErr := ioutil.TempDir("", "profile_tests_")
	checkErr(t, goPathErr)
	defer os.RemoveAll(tempGopathDir)

	tempSrcDir := filepath.Join(tempGopathDir, "src")
	mkdirErr := os.Mkdir(tempSrcDir, 0755)
	checkErr(t, mkdirErr)

	tempMainPath := filepath.Join(tempSrcDir, "main.go")
	mainErr := ioutil.WriteFile(tempMainPath, []byte(codeToTest), 0644)
	checkErr(t, mainErr)

	var stdout, stderr bytes.Buffer
	cmd := exec.Command("go", "run", tempMainPath)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	runErr := cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), runErr
}

/*
	checkErr checks if the error provided as input is different than nil.
	In case the error is not nil, the test will fail.
*/
func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

// checkPprofFile checks if input pprof files exist
func checkPprofFiles(t *testing.T, pprofFilesPath []string) {
	for _, pprof := range pprofFilesPath {
		info, err := os.Stat(pprof)
		assert.Nil(t, err)
		assert.False(t, os.IsNotExist(err))
		assert.False(t, info.IsDir())
	}
}

// cleanupPprofFiles deletes all specified pprof files
func cleanupPprofFiles(t *testing.T, pprofFilesPath []string) {
	for _, pprof := range pprofFilesPath {
		err := os.Remove(pprof)
		if err != nil {
			t.Fatal(err)
		}
	}
}
