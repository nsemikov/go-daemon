package daemon

// Daemon interface.
type Daemon interface {
	Install(args ...string) (string, error)
	Pause() (string, error)
	Reload() (string, error)
	Restart() (string, error)
	Continue() (string, error)
	Status() (string, error)
	Start() (string, error)
	Stop() (string, error)
	Uninstall() (string, error)
	Run() error
}
