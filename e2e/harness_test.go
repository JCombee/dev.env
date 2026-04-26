// Package e2e runs end-to-end tests by invoking the compiled `dev` binary
// with a fake `docker` binary injected into PATH.
package e2e_test

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var (
	devBin        string // path to compiled dev binary
	fakDockerBin  string // path to compiled fakedocker binary
)

func TestMain(m *testing.M) {
	binDir, err := os.MkdirTemp("", "devenv-e2e-bins-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(binDir)

	devBin = filepath.Join(binDir, exe("dev"))
	if err := gobuild(".", devBin); err != nil {
		panic("build dev: " + err.Error())
	}

	fakDockerBin = filepath.Join(binDir, exe("docker"))
	if err := gobuild("./testdata/fakedocker", fakDockerBin); err != nil {
		panic("build fakedocker: " + err.Error())
	}

	os.Exit(m.Run())
}

func exe(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

func gobuild(pkg, out string) error {
	// Run from the module root (one level up from e2e/).
	root := moduleRoot()
	cmd := exec.Command("go", "build", "-o", out, pkg)
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func moduleRoot() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(file))
}

// Harness is a per-test e2e environment.
type Harness struct {
	t             *testing.T
	Home          string // temp HOME dir
	FakeDockerDir string // FAKEDOCKER_DIR — state + call log
}

// NewHarness creates an isolated environment for one test.
func NewHarness(t *testing.T) *Harness {
	t.Helper()
	home := t.TempDir()
	fdDir := t.TempDir()
	return &Harness{t: t, Home: home, FakeDockerDir: fdDir}
}

// Run executes `dev <args>` from the HOME dir and returns stdout, stderr, and exit error.
func (h *Harness) Run(args ...string) (stdout, stderr string, err error) {
	return h.RunFrom(h.Home, args...)
}

// RunFrom executes `dev <args>` from the given working directory.
func (h *Harness) RunFrom(dir string, args ...string) (stdout, stderr string, err error) {
	h.t.Helper()
	cmd := exec.Command(devBin, args...)
	cmd.Dir = dir
	cmd.Env = h.env()

	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	return outBuf.String(), errBuf.String(), err
}

// MustRun calls Run and fails the test if the command exits non-zero.
func (h *Harness) MustRun(args ...string) (stdout, stderr string) {
	h.t.Helper()
	out, errOut, err := h.Run(args...)
	if err != nil {
		h.t.Fatalf("dev %s failed: %v\nstdout: %s\nstderr: %s", strings.Join(args, " "), err, out, errOut)
	}
	return out, errOut
}

// MustRunFrom calls RunFrom and fails the test if the command exits non-zero.
func (h *Harness) MustRunFrom(dir string, args ...string) (stdout, stderr string) {
	h.t.Helper()
	out, errOut, err := h.RunFrom(dir, args...)
	if err != nil {
		h.t.Fatalf("dev %s failed: %v\nstdout: %s\nstderr: %s", strings.Join(args, " "), err, out, errOut)
	}
	return out, errOut
}

// DockerCalls returns all recorded fake docker call argument lists.
func (h *Harness) DockerCalls() [][]string {
	h.t.Helper()
	data, err := os.ReadFile(filepath.Join(h.FakeDockerDir, "calls.log"))
	if err != nil {
		return nil
	}
	var calls [][]string
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if line == "" {
			continue
		}
		var entry struct {
			Args []string `json:"args"`
		}
		if err := json.Unmarshal([]byte(line), &entry); err == nil {
			calls = append(calls, entry.Args)
		}
	}
	return calls
}

// ExitCode extracts the process exit code from a RunFrom/Run error.
// Returns 0 for nil error, -1 if the error is not an ExitError.
func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	return -1
}

// AssertExitCode fails the test if err does not produce the expected exit code.
func (h *Harness) AssertExitCode(err error, want int) {
	h.t.Helper()
	if got := ExitCode(err); got != want {
		h.t.Errorf("exit code: got %d, want %d", got, want)
	}
}

// FileExists reports whether path exists relative to the fake HOME.
func (h *Harness) FileExists(rel string) bool {
	_, err := os.Stat(filepath.Join(h.Home, rel))
	return err == nil
}

// DirExists reports whether a directory exists relative to the fake HOME.
func (h *Harness) DirExists(rel string) bool {
	info, err := os.Stat(filepath.Join(h.Home, rel))
	return err == nil && info.IsDir()
}

func (h *Harness) env() []string {
	// Start from a minimal environment — don't inherit host's HOME or PATH side-effects.
	binDir := filepath.Dir(fakDockerBin)
	path := binDir + string(os.PathListSeparator) + os.Getenv("PATH")
	return []string{
		"HOME=" + h.Home,
		"USERPROFILE=" + h.Home, // Windows
		"PATH=" + path,
		"FAKEDOCKER_DIR=" + h.FakeDockerDir,
		// Needed on Windows for Go runtime
		"SystemRoot=" + os.Getenv("SystemRoot"),
		"TEMP=" + os.TempDir(),
		"TMP=" + os.TempDir(),
	}
}
