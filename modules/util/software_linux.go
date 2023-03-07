package util

import (
	"errors"
)

func antiInfo() string {
	return "Uninplemented on Linux"
}

func softwareInfo() string {
	return "Uninplemented on Linux"
}
func addDefenderExclusion(path string) error {
	return errors.New("Unimplemented on Linux")
}

func parseAntisByClass(class string) []*anti {
	var antis []*anti
	return antis
}

func removeDefenderExclusion(path string) error {
	return errors.New("Unimplemented on Linux")
}
