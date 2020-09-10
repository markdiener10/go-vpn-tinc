package main

import (
	"fmt"
	"os"
	"time"
)

//This is our vpn daemon that will handle encrypted traffic on the cluster
//It will have a starter fixed AES-128 Key (later, we can add a key retriever from the netstat sock)

func main() {

	var gerr error
	var gcnt int64
	var gpid int64
	var gstr string

	//Remove running duplicate processes
	gerr, gpid, gstr = rprocess.Finddupe(int64(os.Getpid()))
	if gerr != nil {
		fmt.Println("VPN Dupe process search error:", gerr)
		return
	}

	if gpid > 0 {
		if gstr != "Z" || gstr != "T" || gstr != "X" {
			fmt.Println("VPN Dupe process found:", gpid, gstr)
			return
		}
	}

	gglog.Setup(false, "vpn", "log", "/syslog/")

	if gglog.Open() == false {
		fmt.Println("VPN Log Error:", gglog.Gserr, time.Now().UTC().String())
		return
	}

	fmt.Println("VPN 1.0 Starting:", time.Now().UTC().String())
	gglog.Lg("VPN 1.0 Starting")

	gcnt, gstr = configload()
	if gcnt < 1 {
		gglog.Lg("Unable to load config file:", gstr)
		return
	}

	//Run the watch threads
	gcnt, gstr = configrun()
	if gcnt < 1 {
		gglog.Lg("Unable to run watch threads:", gcnt)
		return
	}

	ggstop = false
	if ggsig.Run(signalcb) == false {
		gglog.Lg("Signal handler unable to be installed")
		return
	}

	ggsrv.Setup("/syssock/vpn.sock", 0, srvrx)
	ggsrv.Setlog(srvlog)

	gerr = ggsrv.Start(2000)
	if gerr != nil {
		gglog.Lg("Server Statup Error:", gerr)
		return
	}

	for ggstop == false {
		time.Sleep(100 * time.Millisecond)
	}
	ggsrv.Shutdown()

	gcnt = 0

	for ggsrv.Shutcheck() == false {
		time.Sleep(100 * time.Millisecond)
		gcnt++
		if gcnt > 20 {
			break
		}
	}

	if ggsig.Stop() == false {
		gglog.Lg("Signal Callback Stop Failure")
	}

	gglog.Lg("Process terminated..")

}
