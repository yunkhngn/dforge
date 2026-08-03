package cmd

import (
	"strings"
	"testing"

	"github.com/yunkhngn/dforge/internal/docker"
)

func TestFormatStatusListsServicesAndPorts(t *testing.T) {
	out := formatStatus([]docker.ServiceStatus{
		{Name: "api", State: "running", Ports: []string{"3000"}},
		{Name: "postgres", State: "running", Health: "healthy", Ports: []string{"5432"}},
	})
	for _, want := range []string{"api", "postgres", "3000", "5432"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q:\n%s", want, out)
		}
	}
}
