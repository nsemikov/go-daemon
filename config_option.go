package daemon

type ConfigOption func(*Config)

func WithName(name string) ConfigOption {
	return func(c *Config) {
		c.Name = name
	}
}

func WithDescription(description string) ConfigOption {
	return func(c *Config) {
		c.Description = description
	}
}

func WithDependencies(dependencies ...string) ConfigOption {
	return func(c *Config) {
		c.Dependencies = dependencies
	}
}

func WithPIDDir(dir string) ConfigOption {
	return func(c *Config) {
		c.PIDDir = dir
	}
}

func WithPIDName(name string) ConfigOption {
	return func(c *Config) {
		c.PIDName = name
	}
}

func WithStartRunLevels(lvls ...int) ConfigOption {
	return func(c *Config) {
		c.StartRunLevels = lvls
	}
}

func WithStopRunLevels(lvls ...int) ConfigOption {
	return func(c *Config) {
		c.StopRunLevels = lvls
	}
}

func WithTemplateLinuxUpstart(template string) ConfigOption {
	return func(c *Config) {
		c.TemplateLinuxUpstart = template
	}
}

func WithTemplateLinuxSystemV(template string) ConfigOption {
	return func(c *Config) {
		c.TemplateLinuxSystemV = template
	}
}

func WithTemplateLinuxSystemD(template string) ConfigOption {
	return func(c *Config) {
		c.TemplateLinuxSystemD = template
	}
}

func WithTemplateMacOSPorpertyList(template string) ConfigOption {
	return func(c *Config) {
		c.TemplateMacOSPorpertyList = template
	}
}

func WithTemplateFreeBSDSystemV(template string) ConfigOption {
	return func(c *Config) {
		c.TemplateFreeBSDSystemV = template
	}
}

func WithHideMethodsWarning(hideMethodsWarning bool) ConfigOption {
	return func(c *Config) {
		c.HideMethodsWarning = hideMethodsWarning
	}
}

func WithWindowsStartMode(mode uint32) ConfigOption {
	return func(c *Config) {
		c.WindowsStartMode = mode
	}
}

func WithWindowsStartAccountName(name string) ConfigOption {
	return func(c *Config) {
		c.WindowsStartAccountName = name
	}
}

func WithWindowsStartAccountPassword(password string) ConfigOption {
	return func(c *Config) {
		c.WindowsStartAccountPassword = password
	}
}

func WithStartHdlr(hdlr func() error) ConfigOption {
	return func(c *Config) {
		c.StartHdlr = hdlr
	}
}

func WithStopHdlr(hdlr func() error) ConfigOption {
	return func(c *Config) {
		c.StopHdlr = hdlr
	}
}

func WithPauseHdlr(hdlr func() error) ConfigOption {
	return func(c *Config) {
		c.PauseHdlr = hdlr
	}
}

func WithContinueHdlr(hdlr func() error) ConfigOption {
	return func(c *Config) {
		c.ContinueHdlr = hdlr
	}
}

func WithReloadHdlr(hdlr func() error) ConfigOption {
	return func(c *Config) {
		c.ReloadHdlr = hdlr
	}
}

func WithRunHdlr(hdlr func() error) ConfigOption {
	return func(c *Config) {
		c.RunHdlr = hdlr
	}
}

func WithErrorHdlr(hdlr func(format string, args ...interface{})) ConfigOption {
	return func(c *Config) {
		c.ErrorHdlr = hdlr
	}
}

func WithInfoHdlr(hdlr func(format string, args ...interface{})) ConfigOption {
	return func(c *Config) {
		c.InfoHdlr = hdlr
	}
}
