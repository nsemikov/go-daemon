//go:build linux
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

func (d *daemonSystemD) pidFile() string {
	name := d.config.PIDName
	if d.config.PIDName == "" {
		name = d.config.Name
	}
	return d.config.PIDDir + "/" + name + ".pid"
}

func (d *daemonSystemD) path() string {
	return "/etc/systemd/system/" + d.config.Name + ".service"
}

func (d *daemonSystemD) installed() bool {
	return checkInstalled(d.path())
}

func (d *daemonSystemD) running() (string, bool) {
	return checkRunning("Active: active", "Main PID: ([0-9]+)", "systemctl", "status", d.config.Name+".service")
}

func (d *daemonSystemD) Install(args ...string) error {
	var (
		srvPath  = d.path()
		file     *os.File
		execPath string
		templ    *template.Template
		err      error
	)
	if file, err = os.Create(srvPath); err != nil {
		return err
	}
	defer file.Close()
	if execPath, err = executablePath(d.config.Name); err != nil {
		return err
	}
	if templ, err = template.New("systemd").Parse(d.config.TemplateLinuxSystemD); err != nil {
		return err
	}
	if err = templ.Execute(
		file,
		&struct {
			Name, Description, Dependencies, Path, Args, PIDFile string
		}{
			d.config.Name,
			d.config.Description,
			strings.Join(d.config.Dependencies, " "),
			execPath,
			strings.Join(args, " "),
			d.pidFile(),
		},
	); err != nil {
		return err
	}
	if err = exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		return err
	}
	err = exec.Command("systemctl", "enable", d.config.Name+".service").Run()
	return err
}

func (d *daemonSystemD) Uninstall() error {
	var err = exec.Command("systemctl", "disable", d.config.Name+".service").Run()
	if err != nil {
		return err
	}
	err = os.Remove(d.path())
	return err
}

func (d *daemonSystemD) Restart() error {
	return exec.Command("systemctl", "restart", d.config.Name+".service").Run()
}

func (d *daemonSystemD) Start() error {
	return exec.Command("systemctl", "start", d.config.Name+".service").Run()
}

func (d *daemonSystemD) Stop() error {
	return exec.Command("systemctl", "stop", d.config.Name+".service").Run()
}

func (d *daemonSystemD) Reload() error {
	return exec.Command("systemctl", "reload", d.config.Name+".service").Run()
}
