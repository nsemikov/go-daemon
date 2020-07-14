// +build linux

package daemon

import (
	"os"
)

type daemon struct {
	specific daemonSpecific
	config   *Config
}

// New is Daemon constructor. Get the Daemon properly.
func New(c *Config) (Daemon, error) {
	var (
		specific daemonSpecific
		err      = c.check()
	)
	if err != nil {
		return nil, err
	}
	// newer subsystem must be checked first
	if _, err = os.Stat("/run/systemd/system"); err == nil {
		specific = &daemonSystemD{c}
	} else if _, err = os.Stat("/sbin/initctl"); err == nil {
		specific = &daemonUpstart{c}
	} else {
		specific = &daemonSystemV{c}
	}
	return &daemon{specific, c}, nil
}

func (s *daemon) Install(args ...string) (string, error) {
	var (
		action = "Install " + s.config.Description + ":"
		ok     bool
		err    error
	)
	if ok, err = checkPrivileges(); !ok {
		return failed(action), err
	}
	if s.specific.installed() {
		return failed(action), ErrAlreadyInstalled
	}
	if err = s.specific.Install(); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (s *daemon) Uninstall() (string, error) {
	var (
		action = "Uninstalling " + s.config.Description + ":"
		ok     bool
		err    error
	)
	if ok, err = checkPrivileges(); !ok {
		return failed(action), err
	}
	if !s.specific.installed() {
		return failed(action), ErrNotInstalled
	}
	if err = s.specific.Uninstall(); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (s *daemon) Restart() (string, error) {
	var (
		action = "Restarting " + s.config.Description + ":"
		ok     bool
		err    error
	)
	if ok, err = checkPrivileges(); !ok {
		return failed(action), err
	}
	if !s.specific.installed() {
		return failed(action), ErrNotInstalled
	}
	if err = s.specific.Restart(); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (s *daemon) Start() (string, error) {
	var (
		action = "Starting " + s.config.Description + ":"
		ok     bool
		err    error
	)
	if ok, err = checkPrivileges(); !ok {
		return failed(action), err
	}
	if !s.specific.installed() {
		return failed(action), ErrNotInstalled
	}
	if _, ok = s.specific.running(); ok {
		return failed(action), ErrAlreadyRunning
	}
	if err = s.specific.Start(); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (s *daemon) Stop() (string, error) {
	var (
		action = "Stopping " + s.config.Description + ":"
		ok     bool
		err    error
	)
	if ok, err = checkPrivileges(); !ok {
		return failed(action), err
	}
	if !s.specific.installed() {
		return failed(action), ErrNotInstalled
	}
	if _, ok = s.specific.running(); !ok {
		return failed(action), ErrAlreadyStopped
	}
	if err = s.specific.Stop(); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (s *daemon) Status() (string, error) {
	if ok, err := checkPrivileges(); !ok {
		return "", err
	}
	if !s.specific.installed() {
		return statusNotInstalled, ErrNotInstalled
	}
	status, _ := s.specific.running()
	return status, nil
}

func (s *daemon) Reload() (string, error) {
	var (
		action = "Stopping " + s.config.Description + ":"
		ok     bool
		err    error
	)
	if ok, err = checkPrivileges(); !ok {
		return failed(action), err
	}
	if !s.specific.installed() {
		return failed(action), ErrNotInstalled
	}
	if _, ok = s.specific.running(); !ok {
		return failed(action), ErrNotStarted
	}
	if err = s.specific.Reload(); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (s *daemon) Pause() (string, error) {
	return "", ErrUnsupportedSystem
}
func (s *daemon) Continue() (string, error) {
	return "", ErrUnsupportedSystem
}

// Run - Run service
func (s *daemon) Run() error {
	return s.config.RunHdlr()
}

type daemonSpecific interface {
	Install(args ...string) error
	Reload() error
	Restart() error
	Start() error
	Stop() error
	Uninstall() error

	path() string
	installed() bool
	running() (string, bool)
}
