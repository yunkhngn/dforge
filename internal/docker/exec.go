package docker

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type execClient struct{}

func New() Client { return &execClient{} }

func (e *execClient) Available() error {
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker is not installed or not in PATH")
	}
	return nil
}

func (e *execClient) Status() ([]ServiceStatus, error) {
	if err := e.Available(); err != nil {
		return nil, err
	}
	// `docker compose ps --format json` emits one JSON object per line.
	out, err := exec.Command("docker", "compose", "ps", "--format", "json").Output()
	if err != nil {
		return nil, fmt.Errorf("docker compose not running or no project here: %w", err)
	}
	var statuses []ServiceStatus
	for _, line := range splitLines(out) {
		if len(line) == 0 {
			continue
		}
		var row struct {
			Service string `json:"Service"`
			State   string `json:"State"`
			Health  string `json:"Health"`
			Publishers []struct {
				PublishedPort int `json:"PublishedPort"`
			} `json:"Publishers"`
		}
		if err := json.Unmarshal(line, &row); err != nil {
			continue
		}
		s := ServiceStatus{Name: row.Service, State: row.State, Health: row.Health}
		for _, p := range row.Publishers {
			if p.PublishedPort != 0 {
				s.Ports = append(s.Ports, fmt.Sprintf("%d", p.PublishedPort))
			}
		}
		statuses = append(statuses, s)
	}
	return statuses, nil
}

func (e *execClient) Logs(service string, follow bool) error {
	args := []string{"compose", "logs"}
	if follow {
		args = append(args, "-f")
	}
	if service != "" {
		args = append(args, service)
	}
	return runInteractive("docker", args...)
}

func (e *execClient) Shell(service string) error {
	// try bash, fall back to sh
	if err := runInteractive("docker", "compose", "exec", service, "bash"); err != nil {
		return runInteractive("docker", "compose", "exec", service, "sh")
	}
	return nil
}

func (e *execClient) Clean(force bool) error {
	args := []string{"system", "prune", "--volumes"}
	if force {
		args = append(args, "-f")
	}
	return runInteractive("docker", args...)
}

func runInteractive(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

func splitLines(b []byte) [][]byte {
	var out [][]byte
	start := 0
	for i, c := range b {
		if c == '\n' {
			out = append(out, b[start:i])
			start = i + 1
		}
	}
	if start < len(b) {
		out = append(out, b[start:])
	}
	return out
}
