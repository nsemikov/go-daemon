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

func (d *daemonUpstart) startRunLevels() []string {
	return runLevels(d.config.StartRunLevels)
}

func (d *daemonUpstart) stopRunLevels() []string {
	return runLevels(d.config.StopRunLevels)
}

func (d *daemonUpstart) path() string {
	return "/etc/init/" + d.config.Name + ".conf"
}

func (d *daemonUpstart) installed() bool {
	return checkInstalled(d.path())
}

func (d *daemonUpstart) running() (string, bool) {
	return checkRunning(d.config.Name+" start/running", "process ([0-9]+)", "service", d.config.Name, "status")
}

func (d *daemonUpstart) Install(args ...string) error {
	var (
		srvPath  = d.path()
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
	execPath, err = executablePath(d.config.Name)
	if err != nil {
		return err
	}
	templ, err = template.New("upstart").Parse(d.config.TemplateLinuxUpstart)
	if err != nil {
		return err
	}
	if err := templ.Execute(
		file,
		&struct {
			Name, Description, Path, Args, StartRunLevels, StopRunLevels string
		}{
			d.config.Name,
			d.config.Description,
			execPath,
			strings.Join(args, " "),
			strings.Join(d.startRunLevels(), ""),
			strings.Join(d.stopRunLevels(), ""),
		},
	); err != nil {
		return err
	}
	if err = os.Chmod(srvPath, 0755); err != nil {
		return err
	}
	return nil
}

func (d *daemonUpstart) Uninstall() error {
	return os.Remove(d.path())
}

func (d *daemonUpstart) Restart() error {
	return exec.Command("restart", d.config.Name).Run()
}

func (d *daemonUpstart) Start() error {
	return exec.Command("start", d.config.Name).Run()
}

func (d *daemonUpstart) Stop() error {
	return exec.Command("stop", d.config.Name).Run()
}

func (d *daemonUpstart) Reload() error {
	return exec.Command("reload", d.config.Name, "reload").Run()
}
