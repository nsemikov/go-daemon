// +build windows

package daemon

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

// NewConfig is a Config constructor. Get the Config properly.
// Fills the generated Config with default templates.
// Sets WindowsStartMode to StartAutomatic.
func NewConfig() *Config {
	return &Config{
		TemplateFreeBSDSystemV:    defaultTemplateFreeBSDSystemV,
		TemplateLinuxSystemD:      defaultTemplateLinuxSystemD,
		TemplateLinuxSystemV:      defaultTemplateLinuxSystemV,
		TemplateLinuxUpstart:      defaultTemplateLinuxUpstart,
		TemplateMacOSPorpertyList: defaultTemplateMacOSPorpertyList,

		ErrorHdlr: func(format string, args ...interface{}) {
			fmt.Fprintf(os.Stderr, format+"\n", args...)
		},
		InfoHdlr: func(format string, args ...interface{}) {
			fmt.Printf(format+"\n", args...)
		},

		WindowsStartMode: mgr.StartAutomatic,
	}
}

func (c Config) acceptedCommands() svc.Accepted {
	if c.PauseHdlr == nil || c.ContinueHdlr == nil {
		return svc.AcceptStop | svc.AcceptShutdown
	}
	return svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
}
