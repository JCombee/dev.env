package compose

// Runner abstracts docker compose calls. The real implementation shells out
// to docker compose; tests inject a mock.
type Runner interface {
	Up(composeFile string, services ...string) error
	Stop(composeFile string, services ...string) error
	PS(composeFile string) ([]ServiceStatus, error)
	Exec(composeFile, service string, args []string) error
}

type ServiceStatus struct {
	Name   string `json:"Name"`
	State  string `json:"State"`
	Status string `json:"Status"`
}
