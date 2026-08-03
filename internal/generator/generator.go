package generator

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/yunkhngn/dforge/internal/detect"
)

type spec struct {
	dockerfile   string // template path
	dockerignore string
}

var specs = map[detect.Framework]spec{
	detect.Go:         {"templates/dockerfiles/go.Dockerfile.tmpl", "templates/dockerignore/go.dockerignore.tmpl"},
	detect.Nextjs:     {"templates/dockerfiles/nextjs.Dockerfile.tmpl", "templates/dockerignore/node.dockerignore.tmpl"},
	detect.React:      {"templates/dockerfiles/react.Dockerfile.tmpl", "templates/dockerignore/node.dockerignore.tmpl"},
	detect.Nestjs:     {"templates/dockerfiles/nestjs.Dockerfile.tmpl", "templates/dockerignore/node.dockerignore.tmpl"},
	detect.Express:    {"templates/dockerfiles/express.Dockerfile.tmpl", "templates/dockerignore/node.dockerignore.tmpl"},
	detect.SpringBoot: {"templates/dockerfiles/springboot.Dockerfile.tmpl", "templates/dockerignore/java.dockerignore.tmpl"},
	detect.Rust:       {"templates/dockerfiles/rust.Dockerfile.tmpl", "templates/dockerignore/rust.dockerignore.tmpl"},
}

func render(path string, data any) (string, error) {
	raw, err := templatesFS.ReadFile(path)
	if err != nil {
		return "", err
	}
	tmpl, err := template.New(path).Parse(string(raw))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Render returns filename->content for the given framework.
func Render(fw detect.Framework) (map[string]string, error) {
	s, ok := specs[fw]
	if !ok {
		return nil, fmt.Errorf("unsupported framework: %s", fw)
	}
	data := map[string]string{"Framework": string(fw)}
	df, err := render(s.dockerfile, data)
	if err != nil {
		return nil, err
	}
	di, err := render(s.dockerignore, data)
	if err != nil {
		return nil, err
	}
	cf, err := render("templates/compose/base.compose.tmpl", data)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"Dockerfile":    df,
		".dockerignore": di,
		"compose.yaml":  cf,
	}, nil
}
