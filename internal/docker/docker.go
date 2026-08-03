package docker

type ServiceStatus struct {
	Name   string
	State  string
	Health string
	Ports  []string
}

type Client interface {
	Available() error
	Status() ([]ServiceStatus, error)
	Logs(service string, follow bool) error
	Shell(service string) error
	Clean(force bool) error
}
