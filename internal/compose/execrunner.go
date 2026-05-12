package compose

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
)

// ExecRunner implements Runner by shelling out to `docker compose`.
type ExecRunner struct{}

// NewExecRunner returns a Runner that calls the real docker CLI.
func NewExecRunner() Runner {
	return &ExecRunner{}
}

func (r *ExecRunner) Up(composeFile string, services ...string) error {
	args := append([]string{"compose", "-f", composeFile, "up", "-d"}, services...)
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *ExecRunner) Stop(composeFile string, services ...string) error {
	args := append([]string{"compose", "-f", composeFile, "stop"}, services...)
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *ExecRunner) PS(composeFile string) ([]ServiceStatus, error) {
	out, err := exec.Command("docker", "compose", "-f", composeFile, "ps", "--format", "json").Output()
	if err != nil {
		return nil, err
	}
	return ParsePS(out)
}

func (r *ExecRunner) Exec(composeFile, service string, args []string) error {
	cmdArgs := append([]string{"compose", "-f", composeFile, "exec", service}, args...)
	cmd := exec.Command("docker", cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ParsePS handles both JSON array and JSONL output from `docker compose ps`.
// Exported for testing.
func ParsePS(data []byte) ([]ServiceStatus, error) {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return nil, nil
	}
	// Try JSON array first (Compose v2.20+).
	if bytes.HasPrefix(data, []byte("[")) {
		var arr []ServiceStatus
		if err := json.Unmarshal(data, &arr); err == nil {
			return arr, nil
		}
	}
	// Fall back to JSONL (one JSON object per line).
	var result []ServiceStatus
	for _, line := range bytes.Split(data, []byte("\n")) {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		var s ServiceStatus
		if err := json.Unmarshal(line, &s); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}
