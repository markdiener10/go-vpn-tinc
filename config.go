package main

import (
	_ "fmt"
	"io/ioutil"
	"os"
	"time"

	"../../lib/comfunc"
	"../../lib/rjson"
)

func configtag(gstag string, gjson *rjson.Tjson, gproc *Tproc) bool {

	var gjproc *rjson.Tjson
	var gjrun *rjson.Tjson
	var gjparm *rjson.Tjson
	var gbytes []byte
	var gerr error

	//func (gjson *Tjson) Find(gkey string, grecurse, gcase, gsubstr bool) *Tjson {
	gjproc = gjson.Find(gstag, false, false, false)
	if gjproc == nil {
		return false
	}

	gjrun = gjproc.Find("RUN", false, false, false)
	if gjrun == nil {
		return false
	}

	gproc.gname = gstag
	gproc.gjparms.Clear()
	gproc.gnrun = gjrun.Gnum
	if gproc.gnrun == 0 {
		return true
	}

	gjparm = gjproc.Find("PARMS", false, false, false)
	if gjparm == nil {
		return true
	}

	gbytes = gjparm.Dump()

	gerr = gproc.gjparms.Loadbyte(gbytes)
	if gerr != nil {
		return false
	}
	return true

}

func configload() (int64, string) {

	var gsfile string
	var gjson rjson.Tjson
	var gbytes []byte
	var gerr error

	if len(os.Args) < 2 {
		gsfile = "/sysrpz/config/procman/configjson.txt"
	} else {
		gsfile = os.Args[1]
	}

	if comfunc.Filevalid(gsfile) == false {
		return -1, "Config file is not valid:" + gsfile
	}

	gbytes, gerr = ioutil.ReadFile(gsfile)
	if gerr != nil {
		return -2, gerr.Error()
	}

	gerr = gjson.Loadbyte(gbytes)
	if gerr != nil {
		return -3, "Invalid json format in config file"
	}

	if configtag("VPN", &gjson, &ggproc[1]) == false {
		return -4, "Invalid VPN tag in config file:" + gsfile
	}

	if configtag("NFT", &gjson, &ggproc[2]) == false {
		return -5, "Invalid VPN tag in config file:" + gsfile
	}

	if configtag("HIA", &gjson, &ggproc[3]) == false {
		return -6, "Invalid VPN tag in config file:" + gsfile
	}

	if configtag("FPM", &gjson, &ggproc[4]) == false {
		return -7, "Invalid VPN tag in config file:" + gsfile
	}

	if configtag("RETHINK", &gjson, &ggproc[5]) == false {
		return -8, "Invalid VPN tag in config file:" + gsfile
	}

	if configtag("ELASTIC", &gjson, &ggproc[6]) == false {
		return -9, "Invalid VPN tag in config file:" + gsfile
	}

	if configtag("NETSTAT", &gjson, &ggproc[7]) == false {
		return -10, "Invalid VPN tag in config file:" + gsfile
	}

	if configtag("CACHE", &gjson, &ggproc[8]) == false {
		return -11, "Invalid VPN tag in config file:" + gsfile
	}

	if configtag("LOCKS", &gjson, &ggproc[9]) == false {
		return -12, "Invalid VPN tag in config file:" + gsfile
	}

	if configtag("DUPE", &gjson, &ggproc[10]) == false {
		return -13, "Invalid VPN tag in config file:" + gsfile
	}

	if configtag("VERIFY", &gjson, &ggproc[11]) == false {
		return -14, "Invalid VPN tag in config file:" + gsfile
	}

	if configtag("DOORBELL", &gjson, &ggproc[12]) == false {
		return -15, "Invalid VPN tag in config file:" + gsfile
	}

	if configtag("CONTPROXY", &gjson, &ggproc[13]) == false {
		return -16, "Invalid VPN tag in config file:" + gsfile
	}

	return 1, "Success"

}

func configwait(gpproc *Tproc, gntime int64) int64 {

	if gpproc == nil {
		return -1
	}

	var gnwait int64 = 0

	for ggstop == false {

		if gpproc.gnexec < 1 {
			time.Sleep(time.Millisecond * 100)
			gnwait += 100
			if gnwait > gntime {
				return -2
			}
		}
		return 1
	}
	return -2

}

func configrun() (int64, string) {

	if ggproc[1].gnrun > 0 {
		go watchvpn(&ggproc[1])

		if configwait(&ggproc[1], 2000) < 1 {
			return -1, "Watch VPN failed to start on time"
		}

	}
	if ggproc[2].gnrun > 0 {
		go watchnft(&ggproc[2])

		if configwait(&ggproc[2], 2000) < 1 {
			return -1, "Watch NFT failed to start on time"
		}

	}

	return 1, ""

}
