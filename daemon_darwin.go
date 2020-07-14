// +build darwin

package daemon

import (
	"html/template"
	"os"
	"os/exec"
)

type daemon struct {
	config *Config
}

// New is Daemon constructor. Get the Daemon properly.
func New(c *Config) (Daemon, error) {
	var err = c.check()
	if err != nil {
		return nil, err
	}
	return &daemon{c}, nil
}

func (s *daemon) path() string {
	return "/etc/systemd/system/" + s.config.Name + ".service"
}

func (s *daemon) installed() bool {
	return checkInstalled(s.path())
}

func (s *daemon) running() (string, bool) {
	return checkRunning(s.config.Name, "PID\" = ([0-9]+);", "launchctl", "list", s.config.Name)
}

func (s *daemon) Install(args ...string) (string, error) {
	action := "Install " + s.config.Description + ":"
	if ok, err := checkPrivileges(); !ok {
		return failed(action), err
	}
	srvPath := s.path()
	if s.installed() {
		return failed(action), ErrAlreadyInstalled
	}
	file, err := os.Create(srvPath)
	if err != nil {
		return failed(action), err
	}
	defer file.Close()
	execPath, err := executablePath(s.config.Name)
	if err != nil {
		return failed(action), err
	}
	templ, err := template.New("propertyList").Parse(s.config.TemplateMacOSPorpertyList)
	if err != nil {
		return failed(action), err
	}
	if err := templ.Execute(
		file,
		&struct {
			Name, Path string
			Args       []string
		}{s.config.Name, execPath, args},
	); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (s *daemon) Uninstall() (string, error) {
	action := "Uninstalling " + s.config.Description + ":"
	if ok, err := checkPrivileges(); !ok {
		return failed(action), err
	}
	if !s.installed() {
		return failed(action), ErrNotInstalled
	}
	if err := os.Remove(s.path()); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (s *daemon) Restart() (string, error) {
	action := "Restarting " + s.config.Description + ":"
	if ok, err := checkPrivileges(); !ok {
		return failed(action), err
	}
	if !s.installed() {
		return failed(action), ErrNotInstalled
	}
	if err := exec.Command("launchctl", "reload", s.path()+".service").Run(); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (s *daemon) Start() (string, error) {
	action := "Starting " + s.config.Description + ":"
	if ok, err := checkPrivileges(); !ok {
		return failed(action), err
	}
	if !s.installed() {
		return failed(action), ErrNotInstalled
	}
	if _, ok := s.running(); ok {
		return failed(action), ErrAlreadyRunning
	}
	if err := exec.Command("launchctl", "load", s.path()+".service").Run(); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (s *daemon) Stop() (string, error) {
	action := "Stopping " + s.config.Description + ":"
	if ok, err := checkPrivileges(); !ok {
		return failed(action), err
	}
	if !s.installed() {
		return failed(action), ErrNotInstalled
	}
	if _, ok := s.running(); !ok {
		return failed(action), ErrAlreadyStopped
	}
	if err := exec.Command("launchctl", "unload", s.path()+".service").Run(); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (s *daemon) Status() (string, error) {
	if ok, err := checkPrivileges(); !ok {
		return "", err
	}
	if !s.installed() {
		return statusNotInstalled, ErrNotInstalled
	}
	status, _ := s.running()
	return status, nil
}

func (s *daemon) Reload() (string, error) {
	return "", ErrUnsupportedSystem
}

func (s *daemon) Pause() (string, error) {
	return "", ErrUnsupportedSystem
}
func (s *daemon) Continue() (string, error) {
	return "", ErrUnsupportedSystem
}

// Run - Run daemon
func (s *daemon) Run() error {
	return s.config.RunHdlr()
}
