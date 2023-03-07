package install

import (
	"errors"
	"log"
)

func install() {
	log.Println("Reimplementing this Feature on Linux")
}

func uninstall() {
	Removal()
}
func readInstallInfo() (err error) {
	err = errors.New("unimplemented on Linux")
	return
}
