package main

import (
	_ "fmt"
	"net"
	_ "strconv"
	"time"

	_ "../../lib/comfunc"
	"../../lib/rjson"
	"../../lib/rudp"
)

/*
	if ggdevs.Startcontrol(2000) == false {
		gglog.Lg("VPN Control Start Error:", ggdevs.Gnerr, string(ggdevs.Gberr[:]))
		return false
	}
*/

func (gvpn *Tvpn) Startcontrol(gntimeout int, gsip string) bool {

	var gerr error

	gvpn.gudpcont.Setup(gsip, 8520, 8520, gvpn.controlrx)
	gvpn.gudpcont.Setlog(gvpn.gfnlog)

	gerr = gvpn.gudpcont.Start(gntimeout)
	if gerr != nil {
		gvpn.seterr(210, 0, "Start Up error:"+gerr.Error())
		return false
	}
	return true
}

func (gvpn *Tvpn) Stopcontrol(gntimeout int) bool {

	gvpn.gudpcont.Shutdown()

	var grep int = 0

	for gvpn.gudpcont.Shutcheck() == false {
		time.Sleep(100 * time.Millisecond)
		grep += 100
		if grep > gntimeout {
			gvpn.seterr(220, 0, "Signal Callback Term Stop Failure")
			return false
		}
	}
	return true
}

func (gvpn *Tvpn) controlrx(gsrv *rudp.Tudp, gconn *net.UDPConn, gadd *net.UDPAddr, grxlen int, grxbuf []byte) {

	var gerr error
	var gj rjson.Tjson

	gerr = gj.Loadbyte(grxbuf)
	if gerr != nil {
		gvpn.seterr(230, 0, "Control Rx JSON format error:"+gerr.Error())
		return
	}

	//grxcnt++
	//fmt.Println("SRV RX Callbackb:", grxcnt, gadd, grxlen, string(grxbuf[0:grxlen]))

	/*
		if grxcnt%1000 == 1 {
			fmt.Println("SRV RX Callbackb:", grxcnt, grxlen, string(grxbuf[0:grxlen]))
		}
	*/

	//var gstr string
	//gstr := `{"TAGG":` + strconv.FormatInt(int64(grxcnt), 10) + `}`
	//gsrv.Txs(gconn, gadd, gstr)
	return
}
