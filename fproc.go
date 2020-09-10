package main

import (
	_ "fmt"
	"os"
	_ "strings"
	"syscall"
	"time"

	"../../lib/rjson"
	"../../lib/rstream"
)

func signalcb(gosig os.Signal) {

	if gosig != syscall.SIGTERM {
		if gosig != syscall.SIGINT {
			return
		}
	}
	gglog.Lg("Callback Shutdown:", time.Now().UTC().String())
	ggstop = true
}

func srvrx(gstream *rstream.Tstream, gconn *rstream.Tstreamconn) bool {

	//fmt.Println("SRV RX Callbacka:", grxcnt, gconn.Gnetconn.RemoteAddr(), gconn.Glen)
	//fmt.Println("SRV RX Callback:", gconn.Glen, string(gconn.Gbuf))
	//var gstr string
	//gstr := `{"TAGG":` + strconv.FormatInt(int64(grxcnt), 10) + `}`

	var gjson rjson.Tjson
	var gerr error

	gjson.Clear()
	gerr = gjson.Loadbyte(gconn.Gbuf)
	if gerr != nil {
		gconn.Setbuf(`{"GCMD":"ERROR","GVAL":"Invalid JSON Format supplied"}`)
	} else {
		msgjson(&gjson)
		gconn.Gbuf = gjson.Dump()
		gconn.Glen = len(gconn.Gbuf)
	}
	return (gstream.Tx(gconn))
}

func srvlog(gmsg ...interface{}) {
	gglog.Lg(gmsg)
}
