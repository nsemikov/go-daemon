package daemon

// Config type
type Config struct {
	// Name of daemon
	Name string
	// Description of daemon
	Description string
	// Dependencies of daemon
	Dependencies []string

	// PIDDir is path to the directory containing the pid file.
	// Default value is "/run".
	// Suitable for for SystemV only
	PIDDir string
	// PIDName is empty by default.
	// In this case, Name will be used instead of PIDName
	// Suitable for for SystemV only
	PIDName string

	// StartRunLevels is start run levels, [2, 3, 4, 5] by default.
	// Suitable for for Linux SystemV and UpStart only
	StartRunLevels []int
	// StopRunLevels is stop run levels, [0, 1, 6] by default.
	// Suitable for for Linux SystemV and UpStart only
	StopRunLevels []int

	// TemplateLinuxUpstart contains template for Linux UpStart service file
	TemplateLinuxUpstart string
	// TemplateLinuxSystemV contains template for Linux SystemV service file
	TemplateLinuxSystemV string
	// TemplateLinuxSystemD contains template for Linux SystemD service unit
	TemplateLinuxSystemD string
	// TemplateMacOSPorpertyList contains template for MacOS property list
	TemplateMacOSPorpertyList string
	// TemplateFreeBSDSystemV contains template for FreeBSD SystemV system
	TemplateFreeBSDSystemV string

	// HideMethodsWarning will hide ErrSomeMethodsNotSpecified warning for Pause,
	// Continue and Reload
	HideMethodsWarning bool

	// WindowsStartMode is a mgr.Config.StartType:
	// https://godoc.org/golang.org/x/sys/windows/svc/mgr#Config
	WindowsStartMode uint32
	// WindowsStartAccountName is a mgr.Config.ServiceStartName:
	// https://godoc.org/golang.org/x/sys/windows/svc/mgr#Config
	WindowsStartAccountName string
	// WindowsStartAccountPassword is a mgr.Config.Password:
	// https://godoc.org/golang.org/x/sys/windows/svc/mgr#Config
	WindowsStartAccountPassword string

	// StartHdlr is non-blocking start service handler
	StartHdlr func() error
	// StopHdlr is non-blocking stop service handler
	StopHdlr func() error
	// Pause is non-blocking pause service handler
	PauseHdlr func() error
	// ContinueHdlr is non-blocking continue service handler
	ContinueHdlr func() error
	// ReloadHdlr is non-blocking reload service handler
	ReloadHdlr func() error
	// RunHdlr is blocking service entry-point.
	// Daemon will ignore other handlers if this setted.
	// Not suitable for Windows services.
	RunHdlr func() error

	// ErrorHdlr is non-blocking error message handler
	ErrorHdlr func(format string, args ...interface{})
	// InfoHdlr is non-blocking information message handler
	InfoHdlr func(format string, args ...interface{})
}

func (c *Config) check() error {
	if c == nil {
		return ErrConfigNotSpecified
	}
	if c.ErrorHdlr == nil {
		return ErrMethodErrorNotSpecified
	}
	if c.InfoHdlr == nil {
		return ErrMethodInfoNotSpecified
	}
	var info []string
	if c.PauseHdlr == nil {
		info = append(info, "Pause")
	}
	if c.ContinueHdlr == nil {
		info = append(info, "Continue")
	}
	if c.ReloadHdlr == nil {
		info = append(info, "Reload")
	}
	if len(info) > 0 && !c.HideMethodsWarning {
		c.InfoHdlr("%v: %s", ErrSomeMethodsNotSpecified, info)
	}
	var err []string
	if c.StartHdlr == nil {
		err = append(err, "Start")
	}
	if c.StopHdlr == nil {
		err = append(err, "Stop")
	}
	if len(err) > 0 && !c.HideMethodsWarning {
		c.ErrorHdlr("%v: %s", ErrSomeMethodsNotSpecified, err)
		return ErrSomeMethodsNotSpecified
	}
	return nil
}
