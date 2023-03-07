package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	//"strconv"

	"github.com/whiterabb17/greatwhite/modules/models"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	gophersocket "github.com/whiterabb17/gopher-socket"
	"github.com/whiterabb17/gopher-socket/transport"
)

var ddebug bool

// Vars injected via ldflags by bundler
// Application Vars
var (
	fs = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	W  *astilectron.Window
)

var (
	AppName            string
	BuiltAt            string
	VersionAstilectron string
	VersionElectron    string
)

var ServerClient *gophersocket.Client

func init() {
	fs.BoolVar(&ddebug, "d", true, "enables the debug mode")
}

func ConnectToTeamServer(addr string, port string, usr, pwd string, local bool) bool {
	//go relayServer()
	if !local {
		_socket, err := gophersocket.Dial(
			gophersocket.GetUrl(addr, 55556, false, "&_u="+usr+"&_p="+pwd), //_tags
			transport.GetDefaultWebsocketTransport())
		if err != nil {
			log.Println("error")
		}
		ServerClient = _socket
	}
	/*
		go func() {
			handleTeamServer(ServerClient)
		}()
	*/
	return true
}

func SendToFrontEnd(window *astilectron.Window, event string, message string) {
	rep := models.Packet{Event: event, Message: message}
	if err := bootstrap.SendMessage(window, event, rep, func(m *bootstrap.MessageIn) {
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
func SendToFrontend(window *astilectron.Window, event string, payload interface{}) {
	if err := bootstrap.SendMessage(window, event, payload, func(m *bootstrap.MessageIn) {
		// Unmarshal payload
		// models.Logger.Println(m)
		// var s string
		// if err := json.Unmarshal(m.Payload, &s); err != nil {
		// 	models.Logger.Println(fmt.Errorf("unmarshaling payload failed: %w", err))
		// 	return
		// }
		// models.Logger.Printf("About modal has been displayed and payload is %s!\n", s)
	}); err != nil {
		log.Println(fmt.Errorf("sending about event failed: %w", err))
	}
}
func SendToRelay(event string, payload interface{}) {
	RelayChannel.Emit(event, payload)
}
func handleTeamServer(sock *gophersocket.Client) {
	for {
		err := sock.On(gophersocket.OnDisconnection, func(h *gophersocket.Channel) {
			log.Println("Disconnected\nWill try to reconnect")
			NewNotification("Disconnected from TeamServer, will try to reconnect", "Disconnected")
			// SendToFrontend(models.MainWindow, "Disconnected from TeamServer", "Will attempt to reconnect to the TeamServer automatically")
		Connect:
			if ConnectToTeamServer("127.0.0.1", "55556", "admin", "admin", false) {
				return
			} else {
				time.Sleep(10 * time.Second)
				goto Connect
			}
		})
		if err != nil {
			log.Println(err)
		}
		err = sock.On("listener", func(c *gophersocket.Channel, packet models.Packet) {
			//c.Emit("log", models.Packet{Event: "Started", Message: packet.Message})
			if packet.Message == "reserved" {
				if err := bootstrap.SendMessage(models.ListenerWindow, "reserved", "reserved", func(m *bootstrap.MessageIn) {
					// Unmarshal payload
					// var s1 string
					// if err := json.Unmarshal(m.Payload, &s1); err != nil {
					// 	log.Println(fmt.Errorf("unmarshaling payload failed: %w", err))
					// 	return
					// }
					//log.Printf("About modal has been displayed and payload is %s!\n", s1)
				}); err != nil {
					log.Println(fmt.Errorf("sending about event failed: %w", err))
				}
				SendToFrontend(models.ListenerWindow, "reserved", "reserved")
			} else if packet.Message == "started" {
				if err := bootstrap.SendMessage(models.ListenerWindow, "started", "okay", func(m *bootstrap.MessageIn) {
					// Unmarshal payload
					var s1 string
					if err := json.Unmarshal(m.Payload, &s1); err != nil {
						log.Println(fmt.Errorf("unmarshaling payload failed: %w", err))
						return
					}
					log.Printf("About modal has been displayed and payload is %s!\n", s1)
				}); err != nil {
					log.Println(fmt.Errorf("sending about event failed: %w", err))
				}
				SendToFrontend(models.ListenerWindow, "started", "started")
			}
		})
		if err != nil {
			log.Println(err)
		}
		err = sock.On("reserve", func(c *gophersocket.Channel, packet models.Packet) {
			log.Println("Reserve Command Recieved!")
			SendToFrontend(models.ListenerWindow, "reserve", "Successfully reserved port: "+packet.Message)
			c.Emit("log", models.Packet{Event: packet.Event, Message: packet.Message})
		})
		if err != nil {
			log.Println(err)
		}
		err = sock.On("startport", func(c *gophersocket.Channel, packet models.Packet) {
			log.Println("StartPort Command Recieved!")
			c.Emit("log", models.Packet{Event: packet.Event, Message: packet.Message})
			SendToFrontend(models.ListenerWindow, "started", packet.Message)
		})
		if err != nil {
			log.Println(err)
		}
		err = sock.On("portmap", func(c *gophersocket.Channel, portList models.Packet) {
			log.Println("PortMap Command Recieved!")
			models.Logger.Printf(fmt.Sprintf("Recieved List: %s", portList.Message))

			SendToFrontend(W, "alert", "Recieving Port list for Team Server")
			SendToFrontend(models.ListenerWindow, "portmap", portList.Message)
			WriteLog(fmt.Sprintf("Tag: %s  Data: %s", "PortMap", portList.Message))
			//	c.Emit("log", models.Packet{Event: "listener", Message: "event recieved"})
		})
		if err != nil {
			log.Println(err)
		}
		err = sock.On("notify", func(c *gophersocket.Channel, data models.Resp) {
			log.Println(data)
			SendToFrontend(W, "notify", data.Resp)
		})
		if err != nil {
			log.Println(err)
		}
		err = sock.On("newClient", func(c *gophersocket.Channel, clientInfo models.Resp) {
			log.Println(clientInfo)
			//data := strings.Split(clientInfo.Resp, "#")
			//SendToRelay("newclient", clientInfo)

			SendToFrontend(W, "cList", strings.Split(clientInfo.Resp, "#")[0])
			SendToFrontend(W, "newClient", clientInfo.Resp)
			WriteLog(fmt.Sprintf("Tag: %s  Data: %s", "Data for ClientSelect", clientInfo.Tag))
			//	SendToFrontend(models.LogWindow, "addLog", "New Client\n Tag: "+data[0]+"\n OS: "+data[1]+" "+"\n BuildLang: Golang")
		})
		if err != nil {
			log.Println(err)
		}
		err = sock.On("alert", func(c *gophersocket.Channel, data models.Resp) {
			log.Println("Alert Command Recieved!")
			SendToFrontend(W, "discon", "Client Disconnected: "+data.Tag+"::"+fmt.Sprintf("Reponse: <i>%s</i>", data.Resp))
			models.RemoveFromKeyMap(data.Tag)
		})
		if err != nil {
			log.Println(err)
		}
		err = sock.On("repl", func(c *gophersocket.Channel, data models.Resp) {
			log.Println("Reply Recieved!")
			SendToFrontend(W, "notify", fmt.Sprintf("Response Recieved\nData: %s\nFrom: %s", data.Resp, data.Tag))
		})
		if err != nil {
			log.Println(err)
		}
		err = sock.On("reg", func(c *gophersocket.Channel, data models.Resp) {
			log.Println(data)
			//SendToRelay("reg", data)
			ntf := strings.Split(data.Resp, "#")
			msg := fmt.Sprintf("ClientID: %s\nOS: %s\nLanguage: %s", ntf[0], ntf[1], ntf[2])
			NewNotification(msg, "New Client")
			//models.RemoveFromKeyMap(data.Tag)
		})
		if err != nil {
			log.Println(err)
		}
		err = sock.On("register", func(c *gophersocket.Channel, clientInfo models.ClientInfo) {
			log.Println(clientInfo)
			//	SendToRelay("register", clientInfo)
			var infoStr string
			infoStr = "SocketID: " + clientInfo.SocketID + "\nPrivilege: " + clientInfo.Priv + "\nVersion: " + clientInfo.Version + "\nHostname: " + clientInfo.Hostname + "\nUser: " + clientInfo.User + "\nIPAddress: " + clientInfo.IPAddr + "\nOperating System: " + clientInfo.OS + " " + clientInfo.Arch + "\nCPU: " + clientInfo.CPU + "\nGPU: " + clientInfo.GPU + "\nMemory: " + clientInfo.Memory + "\nAntiVirus: " + clientInfo.AntiVirus + "\nBuildLanguage: " + clientInfo.Lang
			log.Println(infoStr)
			WriteLog(fmt.Sprintf("Tag: %s  Data: %s", clientInfo.User, infoStr))
			//	SendToRelay("registerStr", models.Packet{Event: "infoStr", Message: infoStr})
			models.ClientList = append(models.ClientList, clientInfo)
			SendToFrontend(W, "registration", clientInfo)
			NewNotification(fmt.Sprintf("Client %s's info was added to the DB successfully", clientInfo.User), "Client Registered")
			//	models.RemoveFromKeyMap(clientInfo.Tag)
		})
		if err != nil {
			log.Println(err)
		}
		err = sock.On("discon", func(c *gophersocket.Channel, data models.Resp) {
			WriteLog(fmt.Sprintf("Tag: %s  Data: %s", data.Tag, data.Resp))
			SendToFrontend(models.MainWindow, "discon", "Client Disconnected: "+data.Tag+"::"+fmt.Sprintf("Reponse: <i>%s</i>", data.Resp))
			models.RemoveFromKeyMap(data.Tag)
			NewNotification(fmt.Sprintf("Client %s was disconnected", data.Tag), "Client Disconnected")
		})
		if err != nil {
			log.Println(err)
		}
	}
}
func WriteLog(logmsg string) {
	_, e := os.ReadDir("Logs")
	if e != nil {
		os.Mkdir("Logs", 0644)
	}
	f, err := os.OpenFile("Logs/operations.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	f.WriteString(logmsg + "\n")
}
func getOptions(windows []*bootstrap.Window) bootstrap.Options {
	return bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
			SingleInstance:     true,
			VersionAstilectron: VersionAstilectron,
			VersionElectron:    VersionElectron,
		},
		Debug:  ddebug, //*debug,
		Logger: models.Logger,
		TrayMenuOptions: []*astilectron.MenuItemOptions{
			{
				Label: astikit.StrPtr("Dashboard"),
				OnClick: func(e astilectron.Event) (deleteListener bool) {
					if models.MainWinCreated {
						models.MainWindow.Focus()
						return
					} else {
						if !models.MainWinCreated {
							models.WindowOpts = nil
							models.WindowOpts = returnMainWindow()
							varr, err := models.CreateWindow(models.WindowOpts)
							if err != nil {
								log.Println(err)
							}
							models.MainWinCreated = true
							models.MainWindow = varr[0]
							return
						} else {
							models.MainWindow.Focus()
							return
						}
					}
				},
			},
			{
				Label: astikit.StrPtr("Listeners"),
				OnClick: func(e astilectron.Event) (deleteListener bool) {
					if !models.ListenerWinCreated {
						models.WindowOpts = nil
						models.WindowOpts = returnListenWindow()
						varr, err := models.CreateWindow(models.WindowOpts)
						if err != nil {
							log.Println(err)
						}
						models.ListenerWinCreated = true
						models.ListenerWindow = varr[0]
						time.Sleep(5 * time.Second)
						filebuffer, err := ioutil.ReadFile("./.listeners")
						if err != nil {
							log.Fatal(err)
						}
						var inputdata string = string(filebuffer)
						SendToFrontend(models.ListenerWindow, "ports", inputdata)
						return
					} else {
						models.ListenerWindow.Focus()
						return
					}
				},
			},
			{
				Label: astikit.StrPtr("Builder"),
				OnClick: func(e astilectron.Event) (buildComplete bool) {
					models.WindowOpts = nil
					models.WindowOpts = returnBuildWindow()
					varr, err := models.CreateWindow(models.WindowOpts)
					if err != nil {
						log.Println(err)
					}
					models.BuildWinCreated = true
					models.BuildWindow = varr[0]
					return
				},
			},
		},
		MenuOptions: []*astilectron.MenuItemOptions{
			{
				Label: astikit.StrPtr("Server"),
				SubMenu: []*astilectron.MenuItemOptions{
					{
						Label: astikit.StrPtr("Listeners"),
						OnClick: func(e astilectron.Event) (deleteListener bool) {
							models.WindowOpts = nil
							models.WindowOpts = returnListenWindow()
							varr, err := models.CreateWindow(models.WindowOpts)
							if err != nil {
								log.Println(err)
							}
							models.ListenerWinCreated = true
							models.ListenerWindow = varr[0]
							time.Sleep(5 * time.Second)
							filebuffer, err := ioutil.ReadFile("./.listeners")
							if err != nil {
								log.Fatal(err)
							}
							var inputdata string = string(filebuffer)
							SendToFrontend(models.ListenerWindow, "ports", inputdata)
							return
						},
					},
					{
						Label: astikit.StrPtr("Builder"),
						OnClick: func(e astilectron.Event) (buildComplete bool) {
							models.WindowOpts = nil
							models.WindowOpts = returnBuildWindow()
							varr, err := models.CreateWindow(models.WindowOpts)
							if err != nil {
								log.Println(err)
							}
							models.BuildWinCreated = true
							models.BuildWindow = varr[0]
							return
						},
					},
					{
						Label: astikit.StrPtr("Exit"),
						OnClick: func(e astilectron.Event) (buildComplete bool) {
							os.Exit(0)
							return
						},
					},
					//{Role: astilectron.MenuItemRoleClose},
				},
			},
			{
				Label: astikit.StrPtr("Logs"),
				OnClick: func(e astilectron.Event) (deleteListener bool) {
					models.WindowOpts = nil
					models.WindowOpts = returnLogWindow()
					varr, err := models.CreateWindow(models.WindowOpts)
					if err != nil {
						log.Println(err)
					}
					models.LogWinCreated = true
					models.LogWindow = varr[0]
					return
				},
			},
			{
				Label: astikit.StrPtr("Help"),
				SubMenu: []*astilectron.MenuItemOptions{
					{
						Label: astikit.StrPtr("About"),
						OnClick: func(e astilectron.Event) (deleteListener bool) {
							if err := bootstrap.SendMessage(models.MainWindow, "about", "", func(m *bootstrap.MessageIn) {
								// Unmarshal payload
								var s string
								if err := json.Unmarshal(m.Payload, &s); err != nil {
									log.Println(fmt.Errorf("unmarshaling payload failed: %w", err))
									return
								}
								models.Logger.Printf("About modal has been displayed and payload is %s!\n", s)
							}); err != nil {
								log.Println(fmt.Errorf("sending about event failed: %w", err))
							}
							return
						},
					},
					//{Role: astilectron.MenuItemRoleClose},
				},
			},
		},
		OnWait: func(app *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			models.MainWindow = ws[1]
			models.MainWinCreated = true
			W = ws[1]
			models.MainWindow.Hide()
			models.LoginWindow = ws[0]
			models.LoginWinCreated = true
			models.AppPtr = app
			models.App = *app
			return nil
		},
		RestoreAssets: RestoreAssets,
		Windows:       windows,
	}
}

func returnMainWindow() []*bootstrap.Window {
	return []*bootstrap.Window{
		{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astikit.StrPtr("#333"),
				Center:          astikit.BoolPtr(true),
				Height:          astikit.IntPtr(670),
				Width:           astikit.IntPtr(920),
			},
		},
	}
}

func returnLoginWindow() []*bootstrap.Window {
	return []*bootstrap.Window{
		{
			Homepage:       "login.html",
			MessageHandler: handleServer,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astikit.StrPtr("#333"),
				Center:          astikit.BoolPtr(true),
				Height:          astikit.IntPtr(620),
				Width:           astikit.IntPtr(380),
			},
		},
	}
}

func returnBuildWindow() []*bootstrap.Window {
	return []*bootstrap.Window{
		{
			Homepage:       "build.html",
			MessageHandler: buildHandler,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astikit.StrPtr("#333"),
				Center:          astikit.BoolPtr(true),
				Height:          astikit.IntPtr(550),
				Width:           astikit.IntPtr(400),
			},
		},
	}
}

func returnLogWindow() []*bootstrap.Window {
	return []*bootstrap.Window{
		{
			Homepage:       "logs.html",
			MessageHandler: handleLogs,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astikit.StrPtr("#333"),
				Center:          astikit.BoolPtr(true),
				Height:          astikit.IntPtr(700),
				Width:           astikit.IntPtr(700),
			},
		},
	}
}

func returnListenWindow() []*bootstrap.Window {
	return []*bootstrap.Window{
		{
			Homepage:       "listeners.html",
			MessageHandler: handleListeners,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astikit.StrPtr("#333"),
				Center:          astikit.BoolPtr(true),
				Height:          astikit.IntPtr(580),
				Width:           astikit.IntPtr(850),
			},
		},
	}
}

func returnPassWindow() []*bootstrap.Window {
	return []*bootstrap.Window{
		{
			Homepage:       "passwords.html",
			MessageHandler: handleListeners,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astikit.StrPtr("#333"),
				Center:          astikit.BoolPtr(true),
				Height:          astikit.IntPtr(580),
				Width:           astikit.IntPtr(850),
			},
		},
	}
}

func createProfile(path string, conf models.Profile) (err error) {
	f, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return err
	}
	defer f.Close()
	proStr := "Username:" + conf.Username + "\nPassword:" + conf.Password + "\nC2Addr:" + conf.C2Addr + "\nC2Port:" + conf.C2Port
	_, err = io.WriteString(f, proStr)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

var (
	Uname string
	Pword string
	C2A   string
	C2P   string
)

func main() {
	func() {
		// Create logger
		l := log.New(log.Writer(), log.Prefix(), log.Flags())
		models.Logger = l
		// Parse flags
		fs.Parse(os.Args[1:])
		_d, _ := os.UserConfigDir()
		err := os.Mkdir(_d+"\\Necromancer", 0755)
		if err != nil {
			log.Println(err)
		}
		_, notok := os.Stat(_d + "\\profiles.ts")
		if notok == nil {
			f, err := os.Open(_d + "\\profiles.ts")
			if err != nil {
				log.Println(err)
			}
			defer f.Close()

			b, err := ioutil.ReadAll(f)
			if err != nil {
				log.Println(err)
			}
			cfg := strings.Split(string(b), "\n")
			Uname = strings.Split(cfg[0], ":")[1]
			Pword = strings.Split(cfg[1], ":")[1]
			C2A = strings.Split(cfg[2], ":")[1]
			C2P = strings.Split(cfg[3], ":")[1]
		}
		// Run bootstrap
		initWindows := []*bootstrap.Window{
			{
				Homepage:       "login.html",
				MessageHandler: handleServer,
				Options: &astilectron.WindowOptions{
					BackgroundColor: astikit.StrPtr("#333"),
					Center:          astikit.BoolPtr(true),
					Height:          astikit.IntPtr(600),
					Width:           astikit.IntPtr(380),
				},
			},
			{
				Homepage:       "index.html",
				MessageHandler: handleMessages,
				Options: &astilectron.WindowOptions{
					BackgroundColor: astikit.StrPtr("#333"),
					Center:          astikit.BoolPtr(true),
					Height:          astikit.IntPtr(670),
					Width:           astikit.IntPtr(920),
				},
			},
			/*
				{
					Homepage:       "index.html",
					MessageHandler: handleMessages,
					Options: &astilectron.WindowOptions{
						BackgroundColor: astikit.StrPtr("#333"),
						Center:          astikit.BoolPtr(true),
						Height:          astikit.IntPtr(700),
						Width:           astikit.IntPtr(700),
					},
				},
			*/
		}
		models.Logger.Printf("Running app built at %s\n", BuiltAt)
		models.AppOpts = getOptions(initWindows)
		if err := bootstrap.Run(models.AppOpts); err != nil {
			l.Println(fmt.Errorf("running bootstrap failed: %w", err))
		}
	}()
}
