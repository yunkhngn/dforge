package doctor

type Severity string

const (
	Error Severity = "error"
	Warn  Severity = "warn"
	Info  Severity = "info"
)

type State struct {
	Root             string
	DockerInstalled  bool
	ComposeInstalled bool
	HasDockerfile    bool
	HasCompose       bool
	HasDockerignore  bool
	HasHealthcheck   bool
	HasRestart       bool
	UsesLatest       bool
	RunsAsRoot       bool
	PortConflicts    []string
	MissingVolumes   []string
	HasEnvFile       bool
}

type Result struct {
	ID       string
	Title    string
	Severity Severity
	Passed   bool
	Weight   int
}

type Check struct {
	ID       string
	Title    string
	Severity Severity
	Weight   int
	Eval     func(State) bool
}

type Report struct {
	Results []Result
	Score   int
}

func Run(state State, checks []Check) Report {
	rep := Report{Score: 100}
	for _, c := range checks {
		passed := false
		if c.Eval != nil {
			passed = c.Eval(state)
		}
		if !passed {
			rep.Score -= c.Weight
		}
		rep.Results = append(rep.Results, Result{
			ID:       c.ID,
			Title:    c.Title,
			Severity: c.Severity,
			Passed:   passed,
			Weight:   c.Weight,
		})
	}
	if rep.Score < 0 {
		rep.Score = 0
	}
	if rep.Score > 100 {
		rep.Score = 100
	}
	return rep
}
