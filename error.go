package daemon

import (
	"errors"
)

var (
	// ErrConfigNotSpecified appears if service configuration is not specified
	ErrConfigNotSpecified = errors.New("service configuration is not specified")

	// ErrMethodErrorNotSpecified appears if error method missed in service config
	ErrMethodErrorNotSpecified = errors.New("error method are not specified in service configuration")

	// ErrMethodInfoNotSpecified appears if info method missed in service config
	ErrMethodInfoNotSpecified = errors.New("info method are not specified in service configuration")

	// ErrSomeMethodsNotSpecified appears if any methods missed in service config
	ErrSomeMethodsNotSpecified = errors.New("some methods are not specified in service configuration")

	// ErrUnsupportedSystem appears if try to use service on system which is not supported by this release
	ErrUnsupportedSystem = errors.New("unsupported system")

	// ErrRootPrivileges appears if run installation or deleting the service without root privileges
	ErrRootPrivileges = errors.New("you must have root user privileges. Possibly using 'sudo' command should help")

	// ErrAlreadyInstalled appears if service already installed on the system
	ErrAlreadyInstalled = errors.New("service has already been installed")

	// ErrNotInstalled appears if try to delete service which was not been installed
	ErrNotInstalled = errors.New("service is not installed")

	// ErrAlreadyRunning appears if try to start already running service
	ErrAlreadyRunning = errors.New("service is already running")

	// ErrAlreadyStopped appears if try to stop already stopped service
	ErrAlreadyStopped = errors.New("service has already been stopped")

	// ErrNotStarted appears if try to reload stopped service
	ErrNotStarted = errors.New("service is not started")
)
