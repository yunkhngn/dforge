package docker

type MockClient struct {
	Statuses     []ServiceStatus
	AvailableErr error
	LogsService  string
	ShellService string
	CleanForce   bool
}

func (m *MockClient) Available() error { return m.AvailableErr }

func (m *MockClient) Status() ([]ServiceStatus, error) {
	return m.Statuses, m.AvailableErr
}

func (m *MockClient) Logs(service string, follow bool) error {
	m.LogsService = service
	return nil
}

func (m *MockClient) Shell(service string) error {
	m.ShellService = service
	return nil
}

func (m *MockClient) Clean(force bool) error {
	m.CleanForce = force
	return nil
}
