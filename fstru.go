package main

import (
	"../../lib/rjson"
)

const TPROCERRLEN = 50

type Tproc struct {
	gname    string
	gnrun    int64
	gchklast int64
	gnexec   int64
	gnerr    int64
	gberr    [TPROCERRLEN]byte //Thread safe array

	//Extra config values specific to the process
	gjparms rjson.Tjson
}
