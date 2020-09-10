package main

import (
	"strings"

	"../../lib/rjson"
)

func msgjson(gjson *rjson.Tjson) {

	if gjson == nil {
		return
	}

	var gjcmd *rjson.Tjson
	var gjname *rjson.Tjson

	//func (gjson *Tjson) Find(gkey string, grecurse, gcase, gsubstr bool) *Tjson {
	gjcmd = gjson.Find("GCMD", false, false, false)
	if gjcmd == nil {
		gjson.Loadstr(`{"GCMD":"ERROR","GVAL":"Missing GCMD Tag"}`)
		return
	}

	gjname = gjson.Find("GNAME", false, false, false)
	if gjname == nil {
		gjson.Loadstr(`{"GCMD":"ERROR","GVAL":"Missing GNAME Tag"}`)
		return
	}

	if strings.Index("START,STOP,STATS", gjcmd.Gstr) == -1 {
		gjson.Loadstr(`{"GCMD":"ERROR","GVAL":"GCMD invalid value"}`)
		return
	}

	if strings.Index("TINC,NFT,HIA,FPM,RETHINK,ELASTIC,NETSTAT,CACHE,LOCKS,DUPE,VERIFY,CONTPROXY,DOORBELL", gjname.Gstr) == -1 {
		gjson.Loadstr(`{"GCMD":"ERROR","GVAL":"GNAME invalid value"}`)
		return
	}

	gjson.Loadstr(`{"GCMD":` + gjcmd.Gstr + `","GVAL":"Command received"}`)

}
