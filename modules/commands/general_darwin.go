package commands

import (
	"os/exec"

	gophersocket "github.com/whiterabb17/gopher-socket"
	"github.com/whiterabb17/greatwhite/modules/models"
)

func shell(command string, c *gophersocket.Channel, selfTag string) {
	cmd := exec.Command("/bin/zsh", "-c", command)
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
