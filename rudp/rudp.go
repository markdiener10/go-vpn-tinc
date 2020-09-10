package rudp

import (
	"bytes"
	_ "fmt"
	"net"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"../rjson"
	"../rlogger"
)

type Tfnudprx func(gudp *Tudp, gconn *net.UDPConn, gadd *net.UDPAddr, grxlen int, grxbuf []byte)
type Tfnudprxj func(gudp *Tudp, gconn *net.UDPConn, gadd *net.UDPAddr, gcmd string, gjson *rjson.Tjson)

type Tudp struct {
	grun   int
	gstop  int
	gconn  *net.UDPConn
	gip    string
	gnlo   int64
	gnhi   int64
	Gnbind int64
	gfnlog rlogger.Tfnlog
	gfnrx  Tfnudprx
	gfnrxj Tfnudprxj
}

func (gudp *Tudp) Setup(gip string, gnlo int64, gnhi int64) {
	gudp.grun = 0
	gudp.gstop = 0
	gudp.gconn = nil
	gudp.gip = gip
	gudp.gnlo = gnlo
	gudp.gnhi = gnhi
	gudp.Gnbind = 0
	gudp.gfnrx = nil
	gudp.gfnrxj = nil
	gudp.gfnlog = nil
}

func (gudp *Tudp) Setrx(gfnrx Tfnudprx) {
	gudp.gfnrxj = nil
	gudp.gfnrx = gfnrx
}

func (gudp *Tudp) Setrxj(gfnrxj Tfnudprxj) {
	gudp.gfnrx = nil
	gudp.gfnrxj = gfnrxj
}

func (gudp *Tudp) Setlog(gfnlog rlogger.Tfnlog) {
	gudp.gfnlog = gfnlog
}

func (gudp *Tudp) Lg(gmsg ...interface{}) {
	if gudp.gfnlog == nil {
		return
	}
	gudp.gfnlog(gmsg)
}

func (gudp *Tudp) Start(gntimeout int) bool {

	if gudp.grun == 1 {
		return true
	}
	if gudp.grun == 2 {
		if gudp.gstop > 0 {
			gudp.Lg("Running UDP scheduled to stop.")
			return false
		}
		return true
	}

	gudp.gstop = 0
	gudp.grun = 1

	var gsrvadd net.UDPAddr
	var gerr error
	var gcnt int64
	var gbfnd bool = false

	for gcnt = gudp.gnlo; gcnt <= gudp.gnhi; gcnt++ {

		gsrvadd = net.UDPAddr{IP: net.ParseIP(gudp.gip), Port: int(gcnt)}
		gudp.gconn, gerr = net.ListenUDP("udp4", &gsrvadd)
		if gerr == nil {
			gudp.Gnbind = gcnt
			gbfnd = true
			break
		}

	}
	if gbfnd == false {
		gudp.Lg("Unable to bind on address:", gudp.gip, gudp.gnlo, gudp.gnhi, gerr)
		return false
	}

	gudp.Lg("Listen:", gsrvadd)

	//Start a thread and let it run
	go gudp.exec()

	for gudp.grun < 2 {
		if gudp.gstop > 0 {
			gudp.Lg("UDP Listen terminated by stop signal before startup OK.")
			return false
		}
		time.Sleep(20 * time.Millisecond)
		gntimeout -= 20
		if gntimeout > 0 {
			continue
		}
		gudp.gstop = 1
		gudp.Lg("Rudp Listen timed out before startup OK.")
		return false
	}
	//Good startup
	return true
}

func (gudp *Tudp) Stop() {

	if gudp.grun != 2 {
		return
	}
	if gudp.gstop > 0 {
		return
	}
	gudp.gstop = 1
}

func (gudp *Tudp) Stopchk() bool {

	if gudp.grun < 3 {
		return false
	}
	return true
}

func (gudp *Tudp) exec() {

	var gerr error

	var gremadd *net.UDPAddr
	var grxbuf [1500]byte
	var grxlen int = 0

	gudp.grun = 2

	for gudp.gstop == 0 {

		gudp.gconn.SetReadDeadline(time.Now().Add(time.Duration(1000) * time.Millisecond))

		grxlen, gremadd, gerr = gudp.gconn.ReadFromUDP(grxbuf[:])
		if gerr != nil {
			if strings.Contains(gerr.Error(), "i/o timeout") == true {
				continue
			}
			gudp.Lg("Read Error:", gerr)
			break
		}

		//fmt.Println("SRV Read:", grxlen, gremadd, grxlen, grxbuf)
		if grxlen < 1 {
			continue
		}
		go gudp.rxsafe(gudp.gconn, gremadd, grxlen, grxbuf[0:grxlen])
	}
	gudp.grun = 3
	gudp.gstop = 0
	gudp.gconn.Close()
}

func (gudp *Tudp) rxsafe(gconn *net.UDPConn, gadd *net.UDPAddr, glen int, gbuf []byte) {

	//Prevent memory exceptions from killing our server process
	defer func() {
		gerr := recover()
		if gerr != nil {
			gudp.Lg("UDP Rxsafe PANIC:", gerr, string(debug.Stack()))
		}
		return
	}()

	var gnhed int = 0
	var gndata int = 0
	var gnlen int = 0
	var gnlen64 int64 = 0
	var gshed string
	var gerr error
	var gjson rjson.Tjson
	var gj *rjson.Tjson

	if glen < 1 {
		return
	}

	//ONLY Need 2 digits to identify UDP length always less than 36^2
	//       01234
	//		 123456
	//Search [HEXX
	gnhed = bytes.Index(gbuf[0:int(glen)], []byte("[HE"))
	if gnhed < 0 {
		return
	}

	//Our header must be before the end of the received buffer
	if gnhed+5 >= glen {
		return
	}

	gshed = string(gbuf[gnhed+3]) + string(gbuf[gnhed+4])
	gnlen64, gerr = strconv.ParseInt(gshed, 36, 64)
	if gerr != nil {
		return
	}

	gnlen = int(gnlen64)

	if gnhed+5+gnlen > glen {
		return
	}
	gndata = gnhed + 5

	if gudp.gfnrxj == nil {
		gudp.gfnrx(gudp, gconn, gadd, gnlen, gbuf[gndata:gndata+gnlen])
		return
	}

	gerr = gjson.Loadbyte(gbuf[gndata : gndata+gnlen])
	if gerr != nil {
		if gnlen > 20 {
			gudp.Lg("Invalid JSON Format:", string(gbuf[0:20]))
		} else {
			gudp.Lg("Invalid JSON Format:", string(gbuf[0:gnlen]))
		}
		return
	}

	gj = gjson.Qfind("GCMD")
	if gj == nil {
		gj = gjson.Qfind("GERR")
	}
	if gj == nil {
		if gnlen > 20 {
			gudp.Lg("RxSafe Missing GCMD/GERR Tag:", string(gbuf[0:20]))
		} else {
			gudp.Lg("Missing GCMD/GERR Tag:", string(gbuf[0:gnlen]))
		}
		return
	}
	gudp.gfnrxj(gudp, gconn, gadd, gj.Gstr, &gjson)
	return
}

func (gudp *Tudp) Txj(gconn *net.UDPConn, gadd *net.UDPAddr, gjson *rjson.Tjson) {
	var gj *rjson.Tjson = gjson.Qfind("GERR")
	gtxbuf := gjson.Dump()
	gtxlen := len(gtxbuf)
	if gj == nil {
		if gtxlen > 20 {
			gudp.Lg("Txj() Missing GERR Tag:", string(gtxbuf[0:20]))
		} else {
			gudp.Lg("Txj() Missing GERR Tag:", string(gtxbuf[0:gtxlen]))
		}
		return
	}
	gudp.Txb(gconn, gadd, gtxlen, gtxbuf)
}

func (gudp *Tudp) Txs(gconn *net.UDPConn, gadd *net.UDPAddr, gstx string) {
	gtxbuf := []byte(gstx)
	gtxlen := len(gtxbuf)
	gudp.Txb(gconn, gadd, gtxlen, gtxbuf)
}

func (gudp *Tudp) Txb(gconn *net.UDPConn, gadd *net.UDPAddr, gtxlen int, gtxbuf []byte) {

	var gtxnew = 0
	var gerr error

	if gudp.gstop > 0 {
		return
	}

	//1296 is maximum size for UDP packet, better at less than 1024
	if gtxlen > 1296 {
		gudp.Lg("TX Size too large", gtxlen)
	}
	if gtxlen < 1 {
		gudp.Lg("TX Size too small,", gtxlen)
	}

	//Now add a [HEXX header to our output to end the transaction
	var gshed string

	if gtxlen < 36 {
		gshed = "[HE0" + strconv.FormatInt(int64(gtxlen), 36)
	} else {
		gshed = "[HE" + strconv.FormatInt(int64(gtxlen), 36)
	}

	gbhed := []byte(gshed)
	gtxnew = len(gbhed) + gtxlen
	gbout := append(gbhed, gtxbuf...)

	gconn.SetWriteDeadline(time.Now().Add(time.Duration(1000) * time.Millisecond))

	glen, _, gerr := gconn.WriteMsgUDP(gbout[:], nil, gadd)
	if gerr != nil {
		gudp.Lg("TXB Error:", gerr)
		return
	}
	if glen < gtxnew {
		gudp.Lg("TXB Failed to transmit:", glen, gtxnew)
	}
	//gudp.Lg("TXB:", gadd, glen, gtxnew, gbout[:])
	return
}

func (gudp *Tudp) Rtxj(gip string, gnport int64, gjson *rjson.Tjson) {

	var gj *rjson.Tjson = gjson.Qfind("GCMD")
	gtxbuf := gjson.Dump()
	gtxlen := len(gtxbuf)
	if gj == nil {
		if gtxlen > 20 {
			gudp.Lg("Rtxj() Missing GCMD Tag:", string(gtxbuf[0:20]))
		} else {
			gudp.Lg("Rtxj() Missing GCMD Tag:", string(gtxbuf[0:gtxlen]))
		}
		return
	}
	gudp.Rtxb(gip, gnport, gtxbuf)
}

func (gudp *Tudp) Rtxs(gip string, gnport int64, gmsg string) {
	gtxbuf := []byte(gmsg)
	gudp.Rtxb(gip, gnport, gtxbuf)
}

func (gudp *Tudp) Rtxb(gip string, gnport int64, gb []byte) {

	var gerr error
	var gconn *net.UDPConn = nil
	var gremadd net.UDPAddr
	var gnlen int = len(gb)

	gremadd = net.UDPAddr{IP: net.ParseIP(gip), Port: int(gnport)}

	if gudp.gconn != nil {
		gudp.Txb(gudp.gconn, &gremadd, gnlen, gb)
		return
	}

	gconn, gerr = net.DialUDP("udp4", nil, &gremadd)
	if gerr != nil {
		gudp.Lg("TX Error:", gerr)
		return
	}
	gudp.Txb(gconn, nil, gnlen, gb)
	gconn.Close()

}

//Useful for broadcasting a log of JSON to many different ports
func (gudp *Tudp) Txrange(gip string, gnlo int64, gnhi int64, gmsg string) {

	var gerr error
	var gconn *net.UDPConn = nil
	var gremadd net.UDPAddr
	var gcnt int64

	for gcnt = gnlo; gcnt <= gnhi; gcnt++ {

		gremadd = net.UDPAddr{IP: net.ParseIP(gip), Port: int(gcnt)}

		if gudp.gconn != nil {
			gudp.Txs(gudp.gconn, &gremadd, gmsg)
			return
		}

		gconn, gerr = net.DialUDP("udp4", nil, &gremadd)
		if gerr != nil {
			gudp.Lg("TX Error:", gerr)
			return
		}
		gudp.Txs(gconn, nil, gmsg)
		gconn.Close()

	}

}
