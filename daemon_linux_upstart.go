// +build linux

package daemon

import (
	"html/template"
	"os"
	"os/exec"
	"strings"
)

type daemonUpstart struct {
	config *Config
}

func (s *daemonUpstart) startRunLevels() []string {
	return runLevels(s.config.StartRunLevels)
}

func (s *daemonUpstart) stopRunLevels() []string {
	return runLevels(s.config.StopRunLevels)
}

func (s *daemonUpstart) path() string {
	return "/etc/init/" + s.config.Name + ".conf"
}

func (s *daemonUpstart) installed() bool {
	return checkInstalled(s.path())
}

func (s *daemonUpstart) running() (string, bool) {
	return checkRunning(s.config.Name+" start/running", "process ([0-9]+)", "service", s.config.Name, "status")
}

func (s *daemonUpstart) Install(args ...string) error {
	var (
		srvPath  = s.path()
		file     *os.File
		execPath string
		templ    *template.Template
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
	templ, err = template.New("upstart").Parse(s.config.TemplateLinuxUpstart)
	if err != nil {
		return err
	}
	if err := templ.Execute(
		file,
		&struct {
			Name, Description, Path, Args, StartRunLevels, StopRunLevels string
		}{
			s.config.Name,
			s.config.Description,
			execPath,
			strings.Join(args, " "),
			strings.Join(s.startRunLevels(), ""),
			strings.Join(s.stopRunLevels(), ""),
		},
	); err != nil {
		return err
	}
	if err = os.Chmod(srvPath, 0755); err != nil {
		return err
	}
	return nil
}

func (s *daemonUpstart) Uninstall() error {
	return os.Remove(s.path())
}

func (s *daemonUpstart) Restart() error {
	return exec.Command("restart", s.config.Name).Run()
}

func (s *daemonUpstart) Start() error {
	return exec.Command("start", s.config.Name).Run()
}

func (s *daemonUpstart) Stop() error {
	return exec.Command("stop", s.config.Name).Run()
}

func (s *daemonUpstart) Reload() error {
	return exec.Command("reload", s.config.Name, "reload").Run()
}
