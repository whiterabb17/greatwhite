package install

import (
	"errors"
	"log"
)

func install() {
	log.Println("Reimplementing this Feature on Darwin")
}

func uninstall() {
	Removal()
}

func readInstallInfo() (err error) {
	err = errors.New("unimplemented on Darwin")
	return
}
