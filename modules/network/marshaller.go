package network

import (
	"encoding/json"

	bootstrap "github.com/asticode/go-astilectron-bootstrap"
)

// var packet models.Packet
// err := GetJsonPacket(m, &packet)
func GetJson(m *bootstrap.MessageIn, parseTo interface{}) (err error) {
	if err = json.Unmarshal(m.Payload, &parseTo); err != nil {
		return
	}
	return
}
