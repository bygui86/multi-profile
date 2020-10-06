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
)

// TODO memMode, mutexMode, blockMode, traceMode, threadCreationMode, goroutineMode

type profileTest struct {
	name   string
	code   string
	checks []checkFn
}

type checkFn func(t *testing.T, stdout, stderr []byte, err error)

var profileTests = []profileTest{
	{
		name: "CPU profile",
		code: `
			package main
			
			import "github.com/bygui86/multi-profile"
			
			func main() {
				defer profile.CPUProfile(&profile.ProfileConfig{}).Start().Stop()
			}
			`,
		checks: []checkFn{
			NoStdout,
			Stderr("cpu profiling enabled"),
			NoErr,
		},
	},
}

func TestProfiles(t *testing.T) {
	for _, profTest := range profileTests {
		t.Log(profTest.name)
		stdout, stderr, err := runTest(t, profTest.code)
		for _, check := range profTest.checks {
			check(t, stdout, stderr, err)
		}
	}
}

// NoStdout checks that stdout was blank.
func NoStdout(t *testing.T, stdout, _ []byte, _ error) {
	if length := len(stdout); length > 0 {
		t.Errorf("stdout: wanted 0 bytes, got %d", length)
	}
}

// Stderr verifies that the given lines match the output from stderr
func Stderr(lines ...string) checkFn {
	return func(t *testing.T, _, stderr []byte, _ error) {
		r := bytes.NewReader(stderr)
		if !validateOutput(r, lines) {
			t.Errorf("stderr: expected '%s', actual '%s'", lines, stderr)
		}
	}
}

// NoErr checks that err was nil
func NoErr(t *testing.T, stdout, stderr []byte, err error) {
	if err != nil {
		// t.Errorf("error: expected nil, got '%v'", err)
		t.Errorf("error: expected nil, got '%v'", err)
	}
}

// validatedOutput validates the given slice of lines against data from the given reader.
func validateOutput(reader io.Reader, expected []string) bool {
	s := bufio.NewScanner(reader)
	for _, line := range expected {
		if !s.Scan() || !strings.Contains(s.Text(), line) {
			return false
		}
	}
	return true
}

// runTest executes the go program supplied and returns the contents of stdout,
// stderr and an error which may contain status information about the result
// of the execution.
func runTest(t *testing.T, codeToTest string) ([]byte, []byte, error) {
	// TODO try to replace with an external function
	checkErr := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	tempGopathDir, err := ioutil.TempDir("", "profile-tests-gopath")
	checkErr(err)
	// defer os.RemoveAll(tempGopathDir)

	srcDir := filepath.Join(tempGopathDir, "src")
	err = os.Mkdir(srcDir, 0755)
	checkErr(err)

	srcMainPath := filepath.Join(srcDir, "main.go")
	err = ioutil.WriteFile(srcMainPath, []byte(codeToTest), 0644)
	checkErr(err)

	var stdout, stderr bytes.Buffer
	cmd := exec.Command("go", "run", srcMainPath)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}
