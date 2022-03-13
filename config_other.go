//go:build linux || darwin || freebsd
// +build linux darwin freebsd

package daemon

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

// NewConfig is a Config constructor. Get the Config properly.
// Fills the generated Config with default templates.
func NewConfig(opts ...ConfigOption) *Config {
	cfg := &Config{
		TemplateFreeBSDSystemV:    defaultTemplateFreeBSDSystemV,
		TemplateLinuxSystemD:      defaultTemplateLinuxSystemD,
		TemplateLinuxSystemV:      defaultTemplateLinuxSystemV,
		TemplateLinuxUpstart:      defaultTemplateLinuxUpstart,
		TemplateMacOSPorpertyList: defaultTemplateMacOSPorpertyList,

		PIDDir: "/run",

		StartRunLevels: []int{2, 3, 4, 5},
		StopRunLevels:  []int{0, 1, 6},

		ErrorHdlr: func(format string, args ...interface{}) {
			fmt.Fprintf(os.Stderr, format+"\n", args...)
		},
		InfoHdlr: func(format string, args ...interface{}) {
			fmt.Printf(format+"\n", args...)
		},
	}
	cfg.RunHdlr = defaultRunHdlr(cfg)

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

// DefaultRunHdlr will using instead of nil RunHdlr.
// It is using StartHdlr, StopHdlr and ReloadHdlr.
func defaultRunHdlr(c *Config) func() error {
	return func() (err error) {
		// signal.Ignore(syscall.SIGTSTP)
		var (
			sigChan  = make(chan os.Signal, 1)
			exitChan = make(chan int)
			signals  = []os.Signal{
				syscall.SIGINT,
				syscall.SIGTERM,
				syscall.SIGQUIT,
				syscall.SIGTSTP,
			}
		)
		if c.ReloadHdlr != nil {
			signals = append(signals, syscall.SIGHUP)
		}
		signal.Notify(sigChan, signals...)
		go func(ch <-chan os.Signal) {
			defer func() {
				str := recover()
				if str != nil {
					str = fmt.Sprintf("%s; %s", str, string(debug.Stack()))
					err = fmt.Errorf("run() panic: %s", str)
					exitChan <- 100
				}
			}()
			if err = c.StartHdlr(); err != nil {
				exitChan <- 1
			}
			for {
				sig := <-ch
				sigName, sigDesc := sigNameDesc(sig)
				c.InfoHdlr("SIGNAL %s catched (%s)", sigName, sigDesc)
				if sigName == "reload" {
					if err = c.ReloadHdlr(); err != nil {
						exitChan <- 1
					}
				} else {
					if err = c.StopHdlr(); err != nil {
						exitChan <- 1
					}
					exitChan <- 1
				}
			}
		}(sigChan)
		<-exitChan
		return err
	}
}

func sigNameDesc(sig os.Signal) (sigName, sigDesc string) {
	switch sig {
	case syscall.SIGHUP:
		return "reload", "kill -SIGHUP XXXX"
	case syscall.SIGINT:
		return "terminate", "kill -SIGINT XXXX or `Ctrl+C`"
	case syscall.SIGTERM:
		return "force terminate (kill)", "kill -SIGTERM XXXX or kill -9 XXXX"
	case syscall.SIGQUIT:
		return "terminate and core dump", "kill -SIGQUIT XXXX or `Ctrl+\\`"
	case syscall.SIGTSTP:
		return "console terminate", "kill -SIGTSTP XXXX or `Ctrl+Z`"
	}
	return "Unknown signal", sig.String()
}
