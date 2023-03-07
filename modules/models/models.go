package models

import (
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	gophersocket "github.com/whiterabb17/gopher-socket"
)

var (
	// Database Variables
	ClientIdList []string
	Status       bool
	ClientList   []ClientInfo
	ClientMap    map[string]ClientInfo
	PortStatus   map[int]bool
	Keys         map[string]string
	DBname       = "necro.db"
)

type PortMap struct {
	Name  string       `mapstructure:"name"`
	Ports map[int]bool `mapstructure:"ports"`
}

type PortMapPacket struct {
	Name  string `json:"name"`
	Ports string `json:"ports"` // Split port map into <port>|<port>|<port> for reconstruction on other end of node
}

var (
	// Server Variables
	Server  *gophersocket.Server
	Network []*Listener
)

var (
	App     astilectron.Astilectron
	AppPtr  *astilectron.Astilectron
	Logger  astikit.StdLogger
	AppOpts bootstrap.Options
)

var (
	// Windows Variables
	WindowOpts         []*bootstrap.Window
	MainWindow         *astilectron.Window = nil
	MainWinCreated     bool                = false
	PassWindow         *astilectron.Window = nil
	PassWinCreated     bool                = false
	LogWindow          *astilectron.Window = nil
	LogWinCreated      bool                = false
	ListenerWindow     *astilectron.Window = nil
	ListenerWinCreated bool                = false
	BuildWindow        *astilectron.Window = nil
	BuildWinCreated    bool                = false
	LoginWindow        *astilectron.Window = nil
	LoginWinCreated    bool                = false
)

type Listener struct {
	Name   string
	Port   int
	Server *gophersocket.Server
}

var ClientDBInfo ClientInfo

type ClientInfo struct {
	SocketID  string `json:"socketid"`
	Priv      string `json:"priv"`
	Version   string `json:"version"`
	IPAddr    string `json:"ipaddr"`
	Hostname  string `json:"hostname"`
	User      string `json:"username"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	CPU       string `json:"cpu"`
	GPU       string `json:"gpu"`
	Memory    string `json:"memory"`
	AntiVirus string `json:"antivirus"`
	Lang      string `json:"lang"`
}

type Spell struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

type Dispatch struct {
	Args   string `json:"args"`
	ArgCnt string `json:"argcnt"`
	Tag    string `json:"tag"`
}

type ToClient struct {
	Tag     string `json:"tag"`
	Command string `json:"command"`
	Data    string `json:"data"`
}

type Mail struct {
	Buffer string `json:"img"`
	Uid    string `json:"uid"`
}
type Resp struct {
	Resp string `json:"resp"`
	Tag  string `json:"tag"`
}

type Agent struct {
	Name    string `json:"name"`
	Info    *ClientInfo
	Channel *gophersocket.Channel
}

var Agents []*Agent

type Room struct {
	Name   string  `json:"name"`
	Count  string  `json:"count"`
	Agents []Agent `mapstructure:"agents"`
}

type NodeStats struct {
	Channels      []*gophersocket.Channel
	RoomOccupants map[*gophersocket.Channel]int // The Channel and its client count
}

type Nodes struct {
	Name  string
	Port  string
	Stats []*NodeStats
}

type Packet struct {
	Event   string `json:"event"`
	Message string `json:"message"`
}

type Profile struct {
	Username string `json:"username"`
	Password string `json:"password"`
	C2Addr   string `json:"c2addr"`
	C2Port   string `json:"c2port"`
}

type Message struct {
	Cmd string `json:"cmd"`
	Uid string `json:"uid"`
}
