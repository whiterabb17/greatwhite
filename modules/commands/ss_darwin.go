package commands

import (
	"sync"

	gophersocket "github.com/whiterabb17/gopher-socket"
	"github.com/whiterabb17/greatwhite/modules/models"
)

func Shoot(wg *sync.WaitGroup, c *gophersocket.Channel, selfTag string) {
	c.Emit("error", models.Resp{Resp: "Reimplementing this Feature on Darwin", Tag: selfTag})
	wg.Done()
}
