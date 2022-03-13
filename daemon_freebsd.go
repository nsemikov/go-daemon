//go:build freebsd
// +build freebsd

package daemon

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type daemon struct {
	config *Config
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
	var err = c.check()
	if err != nil {
		return nil, err
	}
	return &daemon{c}, nil
}

func (d *daemon) pidFile() string {
	name := d.config.PIDName
	if d.config.PIDName == "" {
		name = d.config.Name
	}
	return d.config.PIDDir + "/" + name + ".pid"
}

func (d *daemon) path() string {
	return "/usr/local/etc/rc.d/" + d.config.Name
}

func (d *daemon) installed() bool {
	return checkInstalled(d.path())
}

func (d *daemon) enabled() (bool, error) {
	rcConf, err := os.Open("/etc/rc.conf")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return false, err
	}
	defer rcConf.Close()
	rcData, _ := ioutil.ReadAll(rcConf)
	r, _ := regexp.Compile(`.*` + d.config.Name + `_enable="YES".*`)
	v := string(r.Find(rcData))
	var chrFound, sharpFound bool
	for _, c := range v {
		if c == '#' && !chrFound {
			sharpFound = true // nolint:ineffassign
			break
		} else if !sharpFound && c != ' ' {
			chrFound = true
			break
		}
	}
	return chrFound, nil
}

func (d *daemon) cmd(cmd string) string {
	if ok, err := d.enabled(); !ok || err != nil {
		fmt.Println("Daemon is not enabled, using one" + cmd + " instead")
		cmd = "one" + cmd
	}
	return cmd
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
	templ, err := template.New("bsd").Parse(d.config.TemplateFreeBSDSystemV)
	if err != nil {
		return failed(action), err
	}
	if err := templ.Execute(
		file,
		&struct {
			Name, Description, Path, Args, PIDFile string
		}{d.config.Name, d.config.Description, execPath, strings.Join(args, " "), d.pidFile()},
	); err != nil {
		return failed(action), err
	}
	if err := os.Chmod(srvPath, 0755); err != nil {
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
	if err := exec.Command("service", d.config.Name, d.cmd("start")).Run(); err != nil {
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
	if err := exec.Command("service", d.config.Name, d.cmd("restart")).Run(); err != nil {
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
	if err := exec.Command("service", d.config.Name, d.cmd("stop")).Run(); err != nil {
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

// Run - Run service
func (d *daemon) Run() error {
	return d.config.RunHdlr()
}
