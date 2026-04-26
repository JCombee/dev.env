// fakedocker is a test double for the `docker` CLI.
// It handles the subset of `docker compose` commands used by devenv.
//
// State and call log are written to the directory in $FAKEDOCKER_DIR.
// Set that env var before running any `dev` command in e2e tests.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type state struct {
	Running map[string]bool `json:"running"`
}

type callEntry struct {
	Time string   `json:"time"`
	Args []string `json:"args"`
}

type servicePS struct {
	Name   string `json:"Name"`
	State  string `json:"State"`
	Status string `json:"Status"`
}

func stateDir() string {
	d := os.Getenv("FAKEDOCKER_DIR")
	if d == "" {
		fmt.Fprintln(os.Stderr, "fakedocker: FAKEDOCKER_DIR not set")
		os.Exit(1)
	}
	return d
}

func loadState(dir string) state {
	s := state{Running: map[string]bool{}}
	data, err := os.ReadFile(filepath.Join(dir, "state.json"))
	if err == nil {
		_ = json.Unmarshal(data, &s)
	}
	if s.Running == nil {
		s.Running = map[string]bool{}
	}
	return s
}

func saveState(dir string, s state) {
	data, _ := json.Marshal(s)
	_ = os.WriteFile(filepath.Join(dir, "state.json"), data, 0o644)
}

func logCall(dir string, args []string) {
	entry := callEntry{
		Time: time.Now().UTC().Format(time.RFC3339),
		Args: args,
	}
	data, _ := json.Marshal(entry)
	f, err := os.OpenFile(filepath.Join(dir, "calls.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintln(f, string(data))
}

func main() {
	args := os.Args[1:]
	dir := stateDir()

	logCall(dir, args)

	// Expect: docker compose [-f <file>] <command> [flags] [services...]
	if len(args) < 2 || args[0] != "compose" {
		fmt.Fprintln(os.Stderr, "fakedocker: only 'docker compose' is supported")
		os.Exit(1)
	}

	// Strip flags like -f <file> to find the subcommand.
	rest := args[1:]
	var composeFile string
	var positional []string
	for i := 0; i < len(rest); i++ {
		if rest[i] == "-f" && i+1 < len(rest) {
			composeFile = rest[i+1]
			i++
			continue
		}
		positional = append(positional, rest[i])
	}
	_ = composeFile

	if len(positional) == 0 {
		fmt.Fprintln(os.Stderr, "fakedocker: missing subcommand")
		os.Exit(1)
	}

	subcommand := positional[0]
	subargs := positional[1:]

	s := loadState(dir)

	switch subcommand {
	case "up":
		// Strip flags (-d, --detach, etc.)
		var services []string
		for _, a := range subargs {
			if !strings.HasPrefix(a, "-") {
				services = append(services, a)
			}
		}
		for _, svc := range services {
			s.Running[svc] = true
		}
		saveState(dir, s)

	case "stop":
		var services []string
		for _, a := range subargs {
			if !strings.HasPrefix(a, "-") {
				services = append(services, a)
			}
		}
		for _, svc := range services {
			s.Running[svc] = false
		}
		saveState(dir, s)

	case "ps":
		var rows []servicePS
		for name, running := range s.Running {
			st := "exited"
			if running {
				st = "running"
			}
			rows = append(rows, servicePS{Name: name, State: st, Status: st})
		}
		data, _ := json.Marshal(rows)
		fmt.Println(string(data))

	default:
		fmt.Fprintf(os.Stderr, "fakedocker: unknown subcommand %q\n", subcommand)
		os.Exit(1)
	}
}
