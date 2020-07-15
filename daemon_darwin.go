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

func (d *daemon) path() string {
	return "/etc/systemd/system/" + d.config.Name + ".service"
}

func (d *daemon) installed() bool {
	return checkInstalled(d.path())
}

func (d *daemon) running() (string, bool) {
	return checkRunning(d.config.Name, "PID\" = ([0-9]+);", "launchctl", "list", d.config.Name)
}

func (d *daemon) Install(args ...string) (string, error) {
	action := "Install " + d.config.Description + ":"
	if ok, err := checkPrivileges(); !ok {
		return failed(action), err
	}
	srvPath := d.path()
	if d.installed() {
		return failed(action), ErrAlreadyInstalled
	}
	file, err := os.Create(srvPath)
	if err != nil {
		return failed(action), err
	}
	defer file.Close()
	execPath, err := executablePath(d.config.Name)
	if err != nil {
		return failed(action), err
	}
	templ, err := template.New("propertyList").Parse(d.config.TemplateMacOSPorpertyList)
	if err != nil {
		return failed(action), err
	}
	if err := templ.Execute(
		file,
		&struct {
			Name, Path string
			Args       []string
		}{d.config.Name, execPath, args},
	); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (d *daemon) Uninstall() (string, error) {
	action := "Uninstalling " + d.config.Description + ":"
	if ok, err := checkPrivileges(); !ok {
		return failed(action), err
	}
	if !d.installed() {
		return failed(action), ErrNotInstalled
	}
	if err := os.Remove(d.path()); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (d *daemon) Restart() (string, error) {
	action := "Restarting " + d.config.Description + ":"
	if ok, err := checkPrivileges(); !ok {
		return failed(action), err
	}
	if !d.installed() {
		return failed(action), ErrNotInstalled
	}
	if err := exec.Command("launchctl", "reload", d.path()+".service").Run(); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (d *daemon) Start() (string, error) {
	action := "Starting " + d.config.Description + ":"
	if ok, err := checkPrivileges(); !ok {
		return failed(action), err
	}
	if !d.installed() {
		return failed(action), ErrNotInstalled
	}
	if _, ok := d.running(); ok {
		return failed(action), ErrAlreadyRunning
	}
	if err := exec.Command("launchctl", "load", d.path()+".service").Run(); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (d *daemon) Stop() (string, error) {
	action := "Stopping " + d.config.Description + ":"
	if ok, err := checkPrivileges(); !ok {
		return failed(action), err
	}
	if !d.installed() {
		return failed(action), ErrNotInstalled
	}
	if _, ok := d.running(); !ok {
		return failed(action), ErrAlreadyStopped
	}
	if err := exec.Command("launchctl", "unload", d.path()+".service").Run(); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (d *daemon) Status() (string, error) {
	if ok, err := checkPrivileges(); !ok {
		return "", err
	}
	if !d.installed() {
		return statusNotInstalled, ErrNotInstalled
	}
	status, _ := d.running()
	return status, nil
}

func (d *daemon) Reload() (string, error) {
	return "", ErrUnsupportedSystem
}

func (d *daemon) Pause() (string, error) {
	return "", ErrUnsupportedSystem
}
func (d *daemon) Continue() (string, error) {
	return "", ErrUnsupportedSystem
}

// Run - Run daemon
func (d *daemon) Run() error {
	return d.config.RunHdlr()
}
