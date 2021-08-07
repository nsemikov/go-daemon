// +build linux

package daemon

import (
	"os"
	"strconv"
)

type daemon struct {
	specific daemonSpecific
	config   *Config
}

// Must create Daemon or panic.
func Must(c *Config) Daemon {
	d, err := New(c)
	if err != nil {
		panic(err)
	}
	return d
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

func (d *daemon) Install(args ...string) (string, error) {
	var (
		action = "Install " + d.config.Description + ":"
		ok     bool
		err    error
	)
	if ok, err = checkPrivileges(); !ok {
		return failed(action), err
	}
	if d.specific.installed() {
		return failed(action), ErrAlreadyInstalled
	}
	if err = d.specific.Install(args...); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (d *daemon) Uninstall() (string, error) {
	var (
		action = "Uninstalling " + d.config.Description + ":"
		ok     bool
		err    error
	)
	if ok, err = checkPrivileges(); !ok {
		return failed(action), err
	}
	if !d.specific.installed() {
		return failed(action), ErrNotInstalled
	}
	if err = d.specific.Uninstall(); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (d *daemon) Restart() (string, error) {
	var (
		action = "Restarting " + d.config.Description + ":"
		ok     bool
		err    error
	)
	if ok, err = checkPrivileges(); !ok {
		return failed(action), err
	}
	if !d.specific.installed() {
		return failed(action), ErrNotInstalled
	}
	if err = d.specific.Restart(); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (d *daemon) Start() (string, error) {
	var (
		action = "Starting " + d.config.Description + ":"
		ok     bool
		err    error
	)
	if ok, err = checkPrivileges(); !ok {
		return failed(action), err
	}
	if !d.specific.installed() {
		return failed(action), ErrNotInstalled
	}
	if _, ok = d.specific.running(); ok {
		return failed(action), ErrAlreadyRunning
	}
	if err = d.specific.Start(); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (d *daemon) Stop() (string, error) {
	var (
		action = "Stopping " + d.config.Description + ":"
		ok     bool
		err    error
	)
	if ok, err = checkPrivileges(); !ok {
		return failed(action), err
	}
	if !d.specific.installed() {
		return failed(action), ErrNotInstalled
	}
	if _, ok = d.specific.running(); !ok {
		return failed(action), ErrAlreadyStopped
	}
	if err = d.specific.Stop(); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (d *daemon) Status() (string, error) {
	if ok, err := checkPrivileges(); !ok {
		return "", err
	}
	if !d.specific.installed() {
		return statusNotInstalled, ErrNotInstalled
	}
	status, _ := d.specific.running()
	return status, nil
}

func (d *daemon) Reload() (string, error) {
	var (
		action = "Stopping " + d.config.Description + ":"
		ok     bool
		err    error
	)
	if ok, err = checkPrivileges(); !ok {
		return failed(action), err
	}
	if !d.specific.installed() {
		return failed(action), ErrNotInstalled
	}
	if _, ok = d.specific.running(); !ok {
		return failed(action), ErrNotStarted
	}
	if err = d.specific.Reload(); err != nil {
		return failed(action), err
	}
	return success(action), nil
}

func (d *daemon) Pause() (string, error) {
	return "", ErrUnsupportedSystem
}
func (d *daemon) Continue() (string, error) {
	return "", ErrUnsupportedSystem
}

// Run - Run service
func (d *daemon) Run() error {
	return d.config.RunHdlr()
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

func runLevels(levels []int) []string {
	var result []string
	for _, lvl := range levels {
		result = append(result, strconv.FormatInt(int64(lvl), 10))
	}
	return result
}
