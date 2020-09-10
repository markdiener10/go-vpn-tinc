package main

import (
	"../../lib/rlogger"
	"../../lib/rstream"
	"../../libsrv/rsigproc"
)

var ggsig rsigproc.Tsigproc
var ggstop bool = false
var gglog rlogger.Tlogger
var ggsrv rstream.Tstream

var ggproc [13 + 1]Tproc

//Processes to control
//1-TINC
//2-NFT
//3-HIA (Hiawatha)
//4-FPM (PHP-FPM)
//5-RETHINK
//6-ELASTIC
//7-NETSTAT
//8-CACHE
//9-LOCKS
//10-DUPE
//11-VERIFY
//12-DOORBELL
//13-CONTESTPROXY
