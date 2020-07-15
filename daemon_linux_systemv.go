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

func (d *daemonSystemV) pidFile() string {
	name := d.config.PIDName
	if d.config.PIDName == "" {
		name = d.config.Name
	}
	return d.config.PIDDir + "/" + name + ".pid"
}

func (d *daemonSystemV) startRunLevels() []string {
	return runLevels(d.config.StartRunLevels)
}

func (d *daemonSystemV) stopRunLevels() []string {
	return runLevels(d.config.StopRunLevels)
}

func (d *daemonSystemV) path() string {
	return "/etc/init.d/" + d.config.Name
}

func (d *daemonSystemV) installed() bool {
	return checkInstalled(d.path())
}

func (d *daemonSystemV) running() (string, bool) {
	return checkRunning(d.config.Name, "pid  ([0-9]+)", "service", d.config.Name, "status")
}

func (d *daemonSystemV) Install(args ...string) error {
	var (
		srvPath  = d.path()
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
	execPath, err = executablePath(d.config.Name)
	if err != nil {
		return err
	}
	templ, err = template.New("systemV").Parse(d.config.TemplateLinuxSystemV)
	if err != nil {
		return err
	}
	if err := templ.Execute(
		file,
		&struct {
			Name, Description, Path, Args, PIDFile, StartRunLevels, StopRunLevels string
		}{
			d.config.Name,
			d.config.Description,
			execPath,
			strings.Join(args, " "),
			d.pidFile(),
			strings.Join(d.startRunLevels(), " "),
			strings.Join(d.stopRunLevels(), " "),
		},
	); err != nil {
		return err
	}
	if err = os.Chmod(srvPath, 0755); err != nil {
		return err
	}
	for _, lvl = range d.startRunLevels() {
		if err = os.Symlink(srvPath, "/etc/rc"+lvl+".d/S87"+d.config.Name); err != nil {
			continue
		}
	}
	for _, lvl = range d.stopRunLevels() {
		if err = os.Symlink(srvPath, "/etc/rc"+lvl+".d/K17"+d.config.Name); err != nil {
			continue
		}
	}
	return nil
}

func (d *daemonSystemV) Uninstall() error {
	var (
		err error
		lvl string
	)
	if err = os.Remove(d.path()); err != nil {
		return err
	}
	for _, lvl = range d.startRunLevels() {
		if err = os.Remove("/etc/rc" + lvl + ".d/S87" + d.config.Name); err != nil {
			continue
		}
	}
	for _, lvl = range d.stopRunLevels() {
		if err = os.Remove("/etc/rc" + lvl + ".d/K17" + d.config.Name); err != nil {
			continue
		}
	}
	return nil
}

func (d *daemonSystemV) Restart() error {
	return exec.Command("service", d.config.Name, "restart").Run()
}

func (d *daemonSystemV) Start() error {
	return exec.Command("service", d.config.Name, "start").Run()
}

func (d *daemonSystemV) Stop() error {
	return exec.Command("service", d.config.Name, "stop").Run()
}

func (d *daemonSystemV) Reload() error {
	return exec.Command("service", d.config.Name, "reload").Run()
}
