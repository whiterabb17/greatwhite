package util

import (
	"errors"
)

func antiInfo() string {
	return "Uninplemented on Darwin"
}

func parseAntisByClass(class string) []*anti {
	var antis []*anti
	return antis
}
func softwareInfo() string {
	return "Uninplemented on Darwin"
}
func addDefenderExclusion(path string) error {
	return errors.New("Unimplemented on Darwin")
}

func removeDefenderExclusion(path string) error {
	return errors.New("Unimplemented on Darwin")
}
