package doctor

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
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

	portConflicts, missingVolumes := analyzeCompose(compose)

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
		PortConflicts:    portConflicts,
		MissingVolumes:   missingVolumes,
		HasEnvFile:       exists(".env"),
	}
}

// composeDoc is a minimal view of compose.yaml for static analysis.
type composeDoc struct {
	Services map[string]struct {
		Ports   []string `yaml:"ports"`
		Volumes []string `yaml:"volumes"`
	} `yaml:"services"`
	Volumes map[string]any `yaml:"volumes"`
}

// analyzeCompose returns host ports bound by more than one mapping and named
// volumes referenced by a service but never declared under top-level `volumes:`.
func analyzeCompose(compose string) (portConflicts, missingVolumes []string) {
	if compose == "" {
		return nil, nil
	}
	var doc composeDoc
	if err := yaml.Unmarshal([]byte(compose), &doc); err != nil {
		return nil, nil
	}

	portCount := map[string]int{}
	missingSet := map[string]bool{}
	for _, svc := range doc.Services {
		for _, p := range svc.Ports {
			if host := hostPort(p); host != "" {
				portCount[host]++
			}
		}
		for _, v := range svc.Volumes {
			if name, ok := namedVolume(v); ok {
				if _, declared := doc.Volumes[name]; !declared {
					missingSet[name] = true
				}
			}
		}
	}

	for host, n := range portCount {
		if n > 1 {
			portConflicts = append(portConflicts, host)
		}
	}
	for name := range missingSet {
		missingVolumes = append(missingVolumes, name)
	}
	sort.Strings(portConflicts)
	sort.Strings(missingVolumes)
	return portConflicts, missingVolumes
}

// hostPort extracts the host-side port from a compose short-syntax mapping,
// e.g. "3000:3000" -> "3000", "127.0.0.1:8080:80" -> "8080", "80/tcp" -> "80".
func hostPort(mapping string) string {
	m := strings.SplitN(mapping, "/", 2)[0]
	parts := strings.Split(m, ":")
	switch len(parts) {
	case 1:
		return parts[0]
	case 2:
		return parts[0]
	case 3:
		return parts[1]
	}
	return ""
}

// namedVolume reports whether a volume mapping references a named volume
// (as opposed to a bind mount starting with "." or "/") and returns its name.
func namedVolume(mapping string) (string, bool) {
	parts := strings.SplitN(mapping, ":", 2)
	if len(parts) == 2 && !strings.HasPrefix(parts[0], ".") && !strings.HasPrefix(parts[0], "/") {
		return parts[0], true
	}
	return "", false
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}
