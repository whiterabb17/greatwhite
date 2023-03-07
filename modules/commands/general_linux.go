package commands

import (
	gophersocket "github.com/whiterabb17/gopher-socket"
	"github.com/whiterabb17/greatwhite/modules/models"
	"os/exec"
)

func shell(command string, c *gophersocket.Channel, selfTag string) {
	cmd := exec.Command("/bin/bash", "-c", command)
	b, err := cmd.CombinedOutput()
	out := string(b)
	if err != nil {
		out = err.Error() + "\n" + out
	}
	if out == "" {
		out = "<success>"
	}
	c.Emit("repl", models.Resp{"Result: " + out, selfTag})
}
