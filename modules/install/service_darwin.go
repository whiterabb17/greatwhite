package install

import (
	"errors"
)

type svcHandler struct {
	main func()
}

// TryServiceInstall attempts to install a Windows Service pointing to Hydra.
func TryServiceInstall() (err error) {
	err = errors.New("unimplemented in Darwin")
	return err
}

// UninstallService attempts to uninstall the Windows Service created by Hydra.
func UninstallService() (err error) {
	err = errors.New("unimplemented in Darwin")
	return err
}

// HandleService starts accepting Service Control Commands from the operating system.
func HandleService(polyfunc func()) {
}

func ServiceCheck() bool {

	return false
}
