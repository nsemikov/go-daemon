// +build linux darwin freebsd

package daemon

import (
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	statusNotInstalled = "Service not installed"
)

// Show colored "OK"
func success(action string) string {
	return action + "\t\t\t\t\t[  \033[32mOK\033[0m  ]"
}

// Show colored "FAILED"
func failed(action string) string {
	return action + "\t\t\t\t\t[\033[31mFAILED\033[0m]"
}

// Check root rights to use system service
func checkPrivileges() (bool, error) {
	if output, err := exec.Command("id", "-g").Output(); err == nil {
		if gid, parseErr := strconv.ParseUint(strings.TrimSpace(string(output)), 10, 32); parseErr == nil {
			if gid == 0 {
				return true, nil
			}
			return false, ErrRootPrivileges
		}
	}
	return false, ErrUnsupportedSystem
}

// Check service running or not
func checkInstalled(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Check service running or not
func checkRunning(matchRegexStr, pidRegexStr string, cmd string, args ...string) (string, bool) {
	output, err := exec.Command(cmd, args...).Output()
	if err == nil {
		if matched, err := regexp.MatchString(matchRegexStr, string(output)); err == nil && matched {
			reg := regexp.MustCompile(pidRegexStr)
			data := reg.FindStringSubmatch(string(output))
			if len(data) > 1 {
				return "Service (pid  " + data[1] + ") is running...", true
			}
			return "Service is running...", true
		}
	}
	return "Service is stopped", false
}

// Lookup path for executable file
func executablePath(name string) (string, error) {
	if path, err := exec.LookPath(name); err == nil {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return os.Executable()
}

func runLevels(levels []int) []string {
	var result []string
	for _, lvl := range levels {
		result = append(result, strconv.FormatInt(int64(lvl), 10))
	}
	return result
}
