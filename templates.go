// +build windows linux darwin freebsd

package daemon

const (
	// nolint:gochecknoglobals
	defaultTemplateMacOSPorpertyList = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>KeepAlive</key>
	<true/>
	<key>Label</key>
	<string>{{.Name}}</string>
	<key>ProgramArguments</key>
	<array>
	    <string>{{.Path}}</string>
		{{range .Args}}<string>{{.}}</string>
		{{end}}
	</array>
	<key>RunAtLoad</key>
	<true/>
    <key>WorkingDirectory</key>
    <string>/usr/local/var</string>
    <key>StandardErrorPath</key>
    <string>/usr/local/var/log/{{.Name}}.err</string>
    <key>StandardOutPath</key>
    <string>/usr/local/var/log/{{.Name}}.log</string>
</dict>
</plist>
`

	// nolint:gochecknoglobals
	defaultTemplateFreeBSDSystemV = `#!/bin/sh
#
# PROVIDE: {{.Name}}
# REQUIRE: networking syslog
# KEYWORD:
# Add the following lines to /etc/rc.conf to enable the {{.Name}}:
#
# {{.Name}}_enable="YES"
#
. /etc/rc.subr
name="{{.Name}}"
rcvar="{{.Name}}_enable"
command="{{.Path}}"
pidfile="{{.PIDFile}}"
start_cmd="/usr/sbin/daemon -p $pidfile -f $command {{.Args}}"
load_rc_config $name
run_rc_command "$1"
`

	// nolint:gochecknoglobals
	defaultTemplateLinuxSystemD = `[Unit]
Description={{.Description}}
Requires={{.Dependencies}}
After={{.Dependencies}}
[Service]
ExecStart={{.Path}} {{.Args}}
Restart=on-failure
[Install]
WantedBy=multi-user.target
`

	// nolint:gochecknoglobals
	defaultTemplateLinuxSystemV = `#! /bin/sh
#
#       /etc/rc.d/init.d/{{.Name}}
#
#       Starts {{.Name}} as a daemon
#
# chkconfig: 2345 87 17
# description: Starts and stops a single {{.Name}} instance on this system
### BEGIN INIT INFO
# Provides: {{.Name}} 
# Required-Start: $network $named
# Required-Stop: $network $named
# Default-Start: {{.StartRunLevels}}
# Default-Stop: {{.StopRunLevels}}
# Short-Description: This service manages the {{.Description}}.
# Description: {{.Description}}
### END INIT INFO
#
# Source function library.
#
if [ -f /etc/rc.d/init.d/functions ]; then
    . /etc/rc.d/init.d/functions
fi
exec="{{.Path}}"
servname="{{.Description}}"
proc="{{.Name}}"
pidfile="{{.PIDFile}}"
lockfile="/var/lock/subsys/$proc"
stdoutlog="/var/log/$proc.log"
stderrlog="/var/log/$proc.err"
[ -d $(dirname $lockfile) ] || mkdir -p $(dirname $lockfile)
[ -e /etc/sysconfig/$proc ] && . /etc/sysconfig/$proc
start() {
    [ -x $exec ] || exit 5
    if [ -f $pidfile ]; then
        if ! [ -d "/proc/$(cat $pidfile)" ]; then
            rm $pidfile
            if [ -f $lockfile ]; then
                rm $lockfile
            fi
        fi
    fi
    if ! [ -f $pidfile ]; then
        printf "Starting $servname:\t"
        echo "$(date)" >> $stdoutlog
        $exec {{.Args}} >> $stdoutlog 2>> $stderrlog &
        echo $! > $pidfile
        touch $lockfile
        success
        echo
    else
        # failure
        echo
        printf "$pidfile still exists...\n"
        exit 7
    fi
}
stop() {
    echo -n $"Stopping $servname: "
    killproc -p $pidfile $proc
    retval=$?
    echo
    [ $retval -eq 0 ] && rm -f $lockfile
    return $retval
}
restart() {
    stop
    start
}
reload() {
    echo -n $"Reloading $servname: "
    killproc -p ${pidfile} $proc -HUP
    retval=$?
    if [ $retval -eq 7 ]; then
        failure $"httpd shutdown"
    fi
    echo
}
rh_status() {
    status -p $pidfile $proc
}
rh_status_q() {
    rh_status >/dev/null 2>&1
}
case "$1" in
    start)
        rh_status_q && exit 0
        $1
        ;;
    stop)
        rh_status_q || exit 0
        $1
        ;;
    restart)
        $1
        ;;
    reload)
        rh_status_q || exit 0
        $1
        ;;
    status)
        rh_status
        ;;
    *)
        echo $"Usage: $0 {start|stop|status|restart|reload}"
        exit 2
esac
exit $?
`

	// nolint:gochecknoglobals
	defaultTemplateLinuxUpstart = `# {{.Name}} {{.Description}}
description     "{{.Description}}"
author          "QuickQ <support@quickq.ru>"
start on runlevel [{{.StartRunLevels}}]
stop on runlevel [{{.StopRunLevels}}]
respawn
#kill timeout 5
exec {{.Path}} {{.Args}} >> /var/log/{{.Name}}.log 2>> /var/log/{{.Name}}.err
`
)
