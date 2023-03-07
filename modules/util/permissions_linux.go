package util

import "errors"

// RunningAsAdmin returns whether the current process has administrative privileges.
func runningAsAdmin() bool {
	return false
}

// RunAsAdmin attempts to execute a command with admin rights.
func runAsAdmin(command, arguments string) error {
	return errors.New("[!] Reimplementing this Feature on Linux")
}

// ElevateNormal attempts to relaunch Hydra with admin rights.
// It displays a common UAC prompt to the user, with the name of the executable.
func elevateNormal() error {
	return errors.New("[!] Reimplementing this Feature on Linux")
}

// ElevateDisguised attempts to relaunch Hydra with admin rights.
// It displays a common UAC prompt to the user, containing the details of an
// executable signed by Microsoft, namely powershell.exe.
func elevateDisguised() error {
	return errors.New("[!] Reimplementing this Feature on Linux")
}

func elevateLogic() error {
	return errors.New("[!] Reimplementing this Feature on Linux")
}
