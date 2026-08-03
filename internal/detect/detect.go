package detect

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Framework string

const (
	Nextjs     Framework = "nextjs"
	React      Framework = "react"
	Nestjs     Framework = "nestjs"
	Express    Framework = "express"
	SpringBoot Framework = "springboot"
	Go         Framework = "go"
	Rust       Framework = "rust"
)

type Candidate struct {
	Framework Framework
	Dir       string
}

type pkgJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// Detect inspects marker files in root and returns detected frameworks.
func Detect(root string) ([]Candidate, error) {
	var out []Candidate

	if data, err := os.ReadFile(filepath.Join(root, "package.json")); err == nil {
		var p pkgJSON
		if err := json.Unmarshal(data, &p); err != nil {
			return nil, err
		}
		out = append(out, Candidate{Framework: nodeFramework(p), Dir: root})
	}
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
		out = append(out, Candidate{Framework: Go, Dir: root})
	}
	if _, err := os.Stat(filepath.Join(root, "Cargo.toml")); err == nil {
		out = append(out, Candidate{Framework: Rust, Dir: root})
	}
	if hasSpring(root) {
		out = append(out, Candidate{Framework: SpringBoot, Dir: root})
	}
	return out, nil
}

func nodeFramework(p pkgJSON) Framework {
	has := func(dep string) bool {
		_, a := p.Dependencies[dep]
		_, b := p.DevDependencies[dep]
		return a || b
	}
	switch {
	case has("next"):
		return Nextjs
	case has("@nestjs/core"):
		return Nestjs
	case has("express"):
		return Express
	case has("react"):
		return React
	default:
		return Express // sensible node default
	}
}

func hasSpring(root string) bool {
	for _, f := range []string{"pom.xml", "build.gradle", "build.gradle.kts"} {
		if _, err := os.Stat(filepath.Join(root, f)); err == nil {
			return true
		}
	}
	return false
}
