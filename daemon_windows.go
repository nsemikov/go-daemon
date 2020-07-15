// +build windows

package daemon

import (
	"fmt"
	runtimeDebug "runtime/debug"
	"strconv"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/mgr"
)

type daemon struct {
	config *Config
}

// New is Daemon constructor. Get the Daemon properly.
func New(c *Config) (Daemon, error) {
	var err = c.check()
	if err != nil {
		return nil, err
	}
	return &daemon{c}, nil
}

func (d *daemon) Install(args ...string) (string, error) {
	var (
		action = "Install " + d.config.Description + ":"
		execp  string
		m      *mgr.Mgr
		w      *mgr.Service
		err    error
	)
	if execp, err = execPath(); err != nil {
		return failed(action), err // or getWindowsError(err) ?
	}
	if m, err = mgr.Connect(); err != nil {
		return failed(action), getWindowsError(err)
	}
	defer func() {
		_ = m.Disconnect()
	}()
	if w, err = m.OpenService(d.config.Name); err == nil {
		w.Close()
		return failed(action), getWindowsError(err)
	}
	if w, err = m.CreateService(d.config.Name, execp, mgr.Config{
		DisplayName:      d.config.Name,
		Description:      d.config.Description,
		StartType:        d.config.WindowsStartMode,
		Dependencies:     d.config.Dependencies,
		ServiceStartName: d.config.WindowsStartAccountName,
		Password:         d.config.WindowsStartAccountPassword,
	}, args...); err != nil {
		return failed(action), getWindowsError(err)
	}
	defer w.Close()
	// set recovery action for service
	// restart after 5 seconds for the first 3 times
	// restart after 1 minute, otherwise
	r := []mgr.RecoveryAction{{
		Type:  mgr.ServiceRestart,
		Delay: 5000 * time.Millisecond,
	}, {
		Type:  mgr.ServiceRestart,
		Delay: 5000 * time.Millisecond,
	}, {
		Type:  mgr.ServiceRestart,
		Delay: 5000 * time.Millisecond,
	}, {
		Type:  mgr.ServiceRestart,
		Delay: 60000 * time.Millisecond,
	},
	}
	// set reset period as a day
	if err = w.SetRecoveryActions(r, uint32(86400)); err != nil {
		return failed(action), getWindowsError(err)
	}
	return success(action), nil
}

func (d *daemon) Uninstall() (string, error) {
	var (
		action = "Uninstalling " + d.config.Description + ":"
		m      *mgr.Mgr
		w      *mgr.Service
		err    error
	)
	if m, err = mgr.Connect(); err != nil {
		return failed(action), getWindowsError(err)
	}
	defer func() {
		_ = m.Disconnect()
	}()
	if w, err = m.OpenService(d.config.Name); err != nil {
		return failed(action), getWindowsError(err)
	}
	defer w.Close()
	if err = w.Delete(); err != nil {
		return failed(action), getWindowsError(err)
	}
	return success(action), nil
}

func (d *daemon) Restart() (string, error) {
	var (
		action = "Restarting " + d.config.Description + ":"
		m      *mgr.Mgr
		w      *mgr.Service
		err    error
	)
	if m, err = mgr.Connect(); err != nil {
		return failed(action), getWindowsError(err)
	}
	defer func() {
		_ = m.Disconnect()
	}()
	if w, err = m.OpenService(d.config.Name); err != nil {
		return failed(action), getWindowsError(err)
	}
	defer w.Close()
	if err = controlService(w, svc.Stop, svc.Stopped, getStopTimeout()); err != nil {
		return failed(action), getWindowsError(err)
	}
	if err = w.Start(); err != nil {
		return failed(action), getWindowsError(err)
	}
	return success(action), nil
}

func (d *daemon) Start() (string, error) {
	var (
		action = "Starting " + d.config.Description + ":"
		m      *mgr.Mgr
		w      *mgr.Service
		err    error
	)
	if m, err = mgr.Connect(); err != nil {
		return failed(action), getWindowsError(err)
	}
	defer func() {
		_ = m.Disconnect()
	}()
	if w, err = m.OpenService(d.config.Name); err != nil {
		return failed(action), getWindowsError(err)
	}
	defer w.Close()
	if err = w.Start(); err != nil {
		return failed(action), getWindowsError(err)
	}
	return success(action), nil
}

func (d *daemon) Stop() (string, error) {
	var (
		action = "Stopping " + d.config.Description + ":"
		m      *mgr.Mgr
		w      *mgr.Service
		err    error
	)
	if m, err = mgr.Connect(); err != nil {
		return failed(action), getWindowsError(err)
	}
	defer func() {
		_ = m.Disconnect()
	}()
	if w, err = m.OpenService(d.config.Name); err != nil {
		return failed(action), getWindowsError(err)
	}
	defer w.Close()
	if err = controlService(w, svc.Stop, svc.Stopped, getStopTimeout()); err != nil {
		return failed(action), getWindowsError(err)
	}
	return success(action), nil
}

func (d *daemon) Status() (string, error) {
	var (
		action = "Getting status " + d.config.Description + ":"
		m      *mgr.Mgr
		w      *mgr.Service
		status svc.Status
		err    error
	)
	if m, err = mgr.Connect(); err != nil {
		return failed(action), getWindowsError(err)
	}
	defer func() {
		_ = m.Disconnect()
	}()
	if w, err = m.OpenService(d.config.Name); err != nil {
		return failed(action), getWindowsError(err)
	}
	defer w.Close()
	if status, err = w.Query(); err != nil {
		return failed(action), getWindowsError(err)
	}
	return "Status " + d.config.Description + ":" + getWindowsServiceStateFromUint32(status.State), nil
}

func (d *daemon) Reload() (string, error) {
	return "", ErrUnsupportedSystem
}

func (d *daemon) Pause() (string, error) {
	var (
		action = "Pausing " + d.config.Description + ":"
		m      *mgr.Mgr
		w      *mgr.Service
		err    error
	)
	if m, err = mgr.Connect(); err != nil {
		return failed(action), getWindowsError(err)
	}
	defer func() {
		_ = m.Disconnect()
	}()
	if w, err = m.OpenService(d.config.Name); err != nil {
		return failed(action), getWindowsError(err)
	}
	defer w.Close()
	if err = controlService(w, svc.Pause, svc.Paused, 10*time.Second); err != nil {
		return failed(action), getWindowsError(err)
	}
	return success(action), nil
}
func (d *daemon) Continue() (string, error) {
	var (
		action = "Resuming " + d.config.Description + ":"
		m      *mgr.Mgr
		w      *mgr.Service
		err    error
	)
	if m, err = mgr.Connect(); err != nil {
		return failed(action), getWindowsError(err)
	}
	defer func() {
		_ = m.Disconnect()
	}()
	if w, err = m.OpenService(d.config.Name); err != nil {
		return failed(action), getWindowsError(err)
	}
	defer w.Close()
	if err = controlService(w, svc.Continue, svc.Running, 10*time.Second); err != nil {
		return failed(action), getWindowsError(err)
	}
	return success(action), nil
}

// Run - Run service
func (d *daemon) Run() error {
	var (
		interactive bool
		runit       func(string, svc.Handler) error
		err         error
	)
	interactive, err = svc.IsAnInteractiveSession()
	if err != nil {
		return getWindowsError(err)
	}
	runit = svc.Run
	if interactive {
		runit = debug.Run
	}
	if err = runit(d.config.Name, &serviceHandler{d.config}); err != nil {
		return getWindowsError(err)
	}
	return nil
}

func execPath() (string, error) {
	var n uint32
	b := make([]uint16, syscall.MAX_PATH)
	size := uint32(len(b))

	r0, _, e1 := syscall.MustLoadDLL(
		"kernel32.dll",
	).MustFindProc(
		"GetModuleFileNameW",
	).Call(0, uintptr(unsafe.Pointer(&b[0])), uintptr(size))
	n = uint32(r0)
	if n == 0 {
		return "", e1
	}
	return string(utf16.Decode(b[0:n])), nil
}

func controlService(s *mgr.Service, c svc.Cmd, to svc.State, d time.Duration) error {
	status, err := s.Control(c)
	if err != nil {
		return err
	}

	timeDuration := time.Millisecond * 50

	timeout := time.After(d + timeDuration*2)
	tick := time.NewTicker(timeDuration)
	defer tick.Stop()

	for status.State != to {
		select {
		case <-tick.C:
			status, err = s.Query()
			if err != nil {
				return err
			}
		case <-timeout:
			break
		}
	}
	return nil
}

func getStopTimeout() time.Duration {
	// For default and paths see https://support.microsoft.com/en-us/kb/146092
	defaultTimeout := time.Millisecond * 20000
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control`, registry.READ)
	if err != nil {
		return defaultTimeout
	}
	sv, _, err := key.GetStringValue("WaitToKillServiceTimeout")
	if err != nil {
		return defaultTimeout
	}
	v, err := strconv.Atoi(sv)
	if err != nil {
		return defaultTimeout
	}
	return time.Millisecond * time.Duration(v)
}

func getWindowsServiceStateFromUint32(state svc.State) string {
	switch state {
	case svc.Stopped:
		return "SERVICE_STOPPED"
	case svc.StartPending:
		return "SERVICE_START_PENDING"
	case svc.StopPending:
		return "SERVICE_STOP_PENDING"
	case svc.Running:
		return "SERVICE_RUNNING"
	case svc.ContinuePending:
		return "SERVICE_CONTINUE_PENDING"
	case svc.PausePending:
		return "SERVICE_PAUSE_PENDING"
	case svc.Paused:
		return "SERVICE_PAUSED"
	}
	return "SERVICE_UNKNOWN"
}

type serviceHandler struct {
	config *Config
}

func (hdlr *serviceHandler) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	var cmdsAccepted = hdlr.config.acceptedCommands()
	defer hdlr.executeRecover()
	changes <- svc.Status{State: svc.StartPending}

	fasttick := time.NewTicker(500 * time.Millisecond)
	slowtick := time.NewTicker(2 * time.Second)
	tick := fasttick

	var err error
	if err = hdlr.config.StartHdlr(); err != nil {
		return
	}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

loop:
	for {
		select {
		case <-tick.C:
			break
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// Testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				_ = hdlr.config.StopHdlr()
				changes <- svc.Status{State: svc.StopPending}
				break loop
			case svc.Pause:
				err = hdlr.config.PauseHdlr()
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
				tick = slowtick
			case svc.Continue:
				err = hdlr.config.ContinueHdlr()
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				tick = fasttick
			default:
				hdlr.config.ErrorHdlr("unexpected control request #%d", c)
				continue loop
			}
			if err != nil {
				return
			}
		}
	}
	return
}

func (hdlr *serviceHandler) executeRecover() {
	str := recover()
	if str != nil {
		str = fmt.Sprintf("%s; %s", str, string(runtimeDebug.Stack()))
		hdlr.config.ErrorHdlr("execute panic: %v", str)
	}
}

// Show colored "OK"
func success(action string) string {
	return action + "\t\t\t\t\t[  OK  ]"
}

// Show colored "FAILED"
func failed(action string) string {
	return action + "\t\t\t\t\t[FAILED]"
}
