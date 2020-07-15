// +build linux

package daemon

import (
	"html/template"
	"os"
	"os/exec"
	"strings"
)

type daemonSystemV struct {
	config *Config
}

func (s *daemonSystemV) pidFile() string {
	name := s.config.PIDName
	if s.config.PIDName == "" {
		name = s.config.Name
	}
	return s.config.PIDDir + "/" + name + ".pid"
}

func (s *daemonSystemV) startRunLevels() []string {
	return runLevels(s.config.StartRunLevels)
}

func (s *daemonSystemV) stopRunLevels() []string {
	return runLevels(s.config.StopRunLevels)
}

func (s *daemonSystemV) path() string {
	return "/etc/init.d/" + s.config.Name
}

func (s *daemonSystemV) installed() bool {
	return checkInstalled(s.path())
}

func (s *daemonSystemV) running() (string, bool) {
	return checkRunning(s.config.Name, "pid  ([0-9]+)", "service", s.config.Name, "status")
}

func (s *daemonSystemV) Install(args ...string) error {
	var (
		srvPath  = s.path()
		file     *os.File
		execPath string
		templ    *template.Template
		lvl      string
		err      error
	)
	file, err = os.Create(srvPath)
	if err != nil {
		return err
	}
	defer file.Close()
	execPath, err = executablePath(s.config.Name)
	if err != nil {
		return err
	}
	templ, err = template.New("systemV").Parse(s.config.TemplateLinuxSystemV)
	if err != nil {
		return err
	}
	if err := templ.Execute(
		file,
		&struct {
			Name, Description, Path, Args, PIDFile, StartRunLevels, StopRunLevels string
		}{
			s.config.Name,
			s.config.Description,
			execPath,
			strings.Join(args, " "),
			s.pidFile(),
			strings.Join(s.startRunLevels(), " "),
			strings.Join(s.stopRunLevels(), " "),
		},
	); err != nil {
		return err
	}
	if err = os.Chmod(srvPath, 0755); err != nil {
		return err
	}
	for _, lvl = range s.startRunLevels() {
		if err = os.Symlink(srvPath, "/etc/rc"+lvl+".d/S87"+s.config.Name); err != nil {
			continue
		}
	}
	for _, lvl = range s.stopRunLevels() {
		if err = os.Symlink(srvPath, "/etc/rc"+lvl+".d/K17"+s.config.Name); err != nil {
			continue
		}
	}
	return nil
}

func (s *daemonSystemV) Uninstall() error {
	var (
		err error
		lvl string
	)
	if err = os.Remove(s.path()); err != nil {
		return err
	}
	for _, lvl = range s.startRunLevels() {
		if err = os.Remove("/etc/rc" + lvl + ".d/S87" + s.config.Name); err != nil {
			continue
		}
	}
	for _, lvl = range s.stopRunLevels() {
		if err = os.Remove("/etc/rc" + lvl + ".d/K17" + s.config.Name); err != nil {
			continue
		}
	}
	return nil
}

func (s *daemonSystemV) Restart() error {
	return exec.Command("service", s.config.Name, "restart").Run()
}

func (s *daemonSystemV) Start() error {
	return exec.Command("service", s.config.Name, "start").Run()
}

func (s *daemonSystemV) Stop() error {
	return exec.Command("service", s.config.Name, "stop").Run()
}

func (s *daemonSystemV) Reload() error {
	return exec.Command("service", s.config.Name, "reload").Run()
}
