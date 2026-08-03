package doctor

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func DefaultChecks() []Check {
	return []Check{
		{ID: "docker", Title: "Docker installed", Severity: Error, Weight: 20, Eval: func(s State) bool { return s.DockerInstalled }},
		{ID: "compose", Title: "Docker Compose installed", Severity: Error, Weight: 15, Eval: func(s State) bool { return s.ComposeInstalled }},
		{ID: "dockerfile", Title: "Dockerfile exists", Severity: Error, Weight: 15, Eval: func(s State) bool { return s.HasDockerfile }},
		{ID: "composefile", Title: "compose.yaml exists", Severity: Error, Weight: 15, Eval: func(s State) bool { return s.HasCompose }},
		{ID: "dockerignore", Title: ".dockerignore exists", Severity: Warn, Weight: 5, Eval: func(s State) bool { return s.HasDockerignore }},
		{ID: "healthcheck", Title: "HEALTHCHECK present", Severity: Warn, Weight: 8, Eval: func(s State) bool { return s.HasHealthcheck }},
		{ID: "restart", Title: "restart policy set", Severity: Warn, Weight: 5, Eval: func(s State) bool { return s.HasRestart }},
		{ID: "latest", Title: "no :latest image tag", Severity: Warn, Weight: 8, Eval: func(s State) bool { return !s.UsesLatest }},
		{ID: "root", Title: "not running as root", Severity: Warn, Weight: 8, Eval: func(s State) bool { return !s.RunsAsRoot }},
		{ID: "ports", Title: "no port conflicts", Severity: Warn, Weight: 6, Eval: func(s State) bool { return len(s.PortConflicts) == 0 }},
		{ID: "volumes", Title: "no missing volumes", Severity: Warn, Weight: 5, Eval: func(s State) bool { return len(s.MissingVolumes) == 0 }},
		{ID: "envfile", Title: "environment file present", Severity: Info, Weight: 3, Eval: func(s State) bool { return s.HasEnvFile }},
	}
}

func CollectState(root string) State {
	exists := func(name string) bool {
		_, err := os.Stat(filepath.Join(root, name))
		return err == nil
	}
	dockerfile := readFile(filepath.Join(root, "Dockerfile"))
	compose := readFile(filepath.Join(root, "compose.yaml"))
	if compose == "" {
		compose = readFile(filepath.Join(root, "docker-compose.yml"))
	}
	_, dockerErr := exec.LookPath("docker")

	return State{
		Root:             root,
		DockerInstalled:  dockerErr == nil,
		ComposeInstalled: dockerErr == nil, // docker compose plugin ships with docker
		HasDockerfile:    dockerfile != "",
		HasCompose:       compose != "",
		HasDockerignore:  exists(".dockerignore"),
		HasHealthcheck:   strings.Contains(dockerfile, "HEALTHCHECK") || strings.Contains(compose, "healthcheck"),
		HasRestart:       strings.Contains(compose, "restart"),
		UsesLatest:       strings.Contains(dockerfile, ":latest") || strings.Contains(compose, ":latest"),
		RunsAsRoot:       dockerfile != "" && !strings.Contains(dockerfile, "USER "),
		HasEnvFile:       exists(".env"),
	}
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}
