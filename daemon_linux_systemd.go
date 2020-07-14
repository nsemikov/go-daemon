// +build linux

package daemon

import (
	"html/template"
	"os"
	"os/exec"
	"strings"
)

type daemonSystemD struct {
	config *Config
}

func (s *daemonSystemD) path() string {
	return "/etc/systemd/system/" + s.config.Name + ".service"
}

func (s *daemonSystemD) installed() bool {
	return checkInstalled(s.path())
}

func (s *daemonSystemD) running() (string, bool) {
	return checkRunning("Active: active", "Main PID: ([0-9]+)", "systemctl", "status", s.config.Name+".service")
}

func (s *daemonSystemD) Install(args ...string) error {
	var (
		srvPath  = s.path()
		file     *os.File
		execPath string
		templ    *template.Template
		err      error
	)
	if file, err = os.Create(srvPath); err != nil {
		return err
	}
	defer file.Close()
	if execPath, err = executablePath(s.config.Name); err != nil {
		return err
	}
	if templ, err = template.New("systemd").Parse(s.config.TemplateLinuxSystemD); err != nil {
		return err
	}
	if err = templ.Execute(
		file,
		&struct {
			Name, Description, Dependencies, Path, Args string
		}{
			s.config.Name,
			s.config.Description,
			strings.Join(s.config.Dependencies, " "),
			execPath,
			strings.Join(args, " "),
		},
	); err != nil {
		return err
	}
	if err = exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		return err
	}
	err = exec.Command("systemctl", "enable", s.config.Name+".service").Run()
	return err
}

func (s *daemonSystemD) Uninstall() error {
	var err = exec.Command("systemctl", "disable", s.config.Name+".service").Run()
	if err != nil {
		return err
	}
	err = os.Remove(s.path())
	return err
}

func (s *daemonSystemD) Restart() error {
	return exec.Command("systemctl", "restart", s.config.Name+".service").Run()
}

func (s *daemonSystemD) Start() error {
	return exec.Command("systemctl", "start", s.config.Name+".service").Run()
}

func (s *daemonSystemD) Stop() error {
	return exec.Command("systemctl", "stop", s.config.Name+".service").Run()
}

func (s *daemonSystemD) Reload() error {
	return exec.Command("systemctl", "reload", s.config.Name+".service").Run()
}
