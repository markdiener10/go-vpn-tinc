package main

import (
	_ "fmt"
	_ "strings"
	_ "syscall"
	_ "time"

	"../../lib/comfunc"
	"../../lib/rlogger"
	"../../lib/rudp"
)

const (
	VPNERRBUFLEN = 100
)

type Tvpn struct {
	Gnloc int64
	Gnerr int64
	Gberr [VPNERRBUFLEN]byte

	gfnlog rlogger.Tfnlog

	gudpcont rudp.Tudp
}

func (gvpn *Tvpn) lg(gmsg ...interface{}) {
	if gvpn.gfnlog == nil {
		return
	}
	gvpn.gfnlog(gmsg)
}

func (gvpn *Tvpn) seterr(gnloc int64, gnerr int64, gserr string) {
	gvpn.Gnloc = gnloc
	gvpn.Gnerr = gnerr
	comfunc.Strtobyte(gserr, gvpn.Gberr[:], VPNERRBUFLEN)
	gvpn.lg("VPNerror:", gnloc, gnerr, gserr)
}
