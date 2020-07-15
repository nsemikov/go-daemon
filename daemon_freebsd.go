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

// New is Daemon constructor. Get the Daemon properly.
func New(c *Config) (Daemon, error) {
	var err = c.check()
	if err != nil {
		return nil, err
	}
	return &daemon{c}, nil
}

func (s *daemon) path() string {
	return "/usr/local/etc/rc.d/" + s.config.Name
}

func (s *daemon) installed() bool {
	return checkInstalled(s.path())
}

func (s *daemon) enabled() (bool, error) {
	rcConf, err := os.Open("/etc/rc.conf")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return false, err
	}
	defer rcConf.Close()
	rcData, _ := ioutil.ReadAll(rcConf)
	r, _ := regexp.Compile(`.*` + s.config.Name + `_enable="YES".*`)
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

func (s *daemon) cmd(cmd string) string {
	if ok, err := s.enabled(); !ok || err != nil {
		fmt.Println("Daemon is not enabled, using one" + cmd + " instead")
		cmd = "one" + cmd
	}
	return cmd
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
	templ, err := template.New("bsd").Parse(s.config.TemplateFreeBSDSystemV)
	if err != nil {
		return failed(action), err
	}
	if err := templ.Execute(
		file,
		&struct {
			Name, Description, Path, Args, PIDFile string
		}{s.config.Name, s.config.Description, execPath, strings.Join(args, " "), s.config.pidPath()},
	); err != nil {
		return failed(action), err
	}
	if err := os.Chmod(srvPath, 0755); err != nil {
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
	if err := exec.Command("service", s.config.Name, s.cmd("start")).Run(); err != nil {
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
	if err := exec.Command("service", s.config.Name, s.cmd("restart")).Run(); err != nil {
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
	if err := exec.Command("service", s.config.Name, s.cmd("stop")).Run(); err != nil {
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

// Run - Run service
func (s *daemon) Run() error {
	return s.config.RunHdlr()
}
