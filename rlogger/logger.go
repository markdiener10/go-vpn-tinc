package rlogger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Tfnlog func(gmsg ...interface{})

type Tlogger struct {
	gpath     string
	Gberr     bool
	Gserr     string
	groot     string
	gsuffix   string
	gbopen    bool
	gdate     time.Time
	Gbfake    bool
	Gfakedate time.Time
	gbrotate  bool //Allow for daily log rotation
	glock     sync.Mutex
	gfile     *os.File
}

func (glog *Tlogger) Setup(gbrotate bool, gproot string, gpsuffix string, gpath string) {

	var gdnow time.Time = time.Now()

	glog.Gberr = false
	glog.Gserr = ""

	glog.Gbfake = false
	glog.Gfakedate = time.Date(gdnow.Year(), gdnow.Month(), gdnow.Day(), 0, 0, 0, 0, time.UTC)

	glog.gpath = gpath

	glog.gfile = nil
	glog.gbopen = false
	glog.gbrotate = gbrotate
	glog.groot = gproot
	glog.gsuffix = gpsuffix
	glog.gdate = time.Date(gdnow.Year(), gdnow.Month(), gdnow.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)

}

func (glog *Tlogger) Open() bool {

	var gerr error

	if glog.gbopen == true {
		glog.gfile.Close()
		glog.gbopen = false
	}

	glog.gfile, gerr = os.OpenFile(glog.gpath+glog.groot+"."+glog.gsuffix, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if gerr != nil {
		glog.Gberr = true
		glog.Gserr = "Open Err:" + gerr.Error()
		return false
	}
	glog.gbopen = true
	return true
}

func (glog *Tlogger) Lg(gmsg ...interface{}) {

	var gerr error
	var gdyn string

	var gdnow time.Time = time.Now()

	//Today at midnight before 1st second of day
	var gdate time.Time = time.Date(gdnow.Year(), gdnow.Month(), gdnow.Day(), 0, 0, 0, 0, time.UTC)

	//fmt.Println(gdate,glog.gdate)

	if glog.gbrotate == false {
		goto LOGMSG
	}

	if glog.Gbfake == true {
		gdate = glog.Gfakedate
	}

	//If our date has changed to a new day, we want to close the current file
	//and rename it to our original date
	if glog.gdate == gdate {
		goto LOGMSG
	}

	//Try to get exclusive access
	glog.glock.Lock()
	if glog.gdate == gdate {
		glog.glock.Unlock()
		goto LOGMSG
	}

	//If open, close it
	if glog.gbopen == true {
		glog.gfile.Close()
		glog.gbopen = false
	}

	//Rename file to ROOT-DDMMYY-UNIX.SUFFIX
	gdyn = fmt.Sprintf("%02d:%02d:%02d:%03d-%04d%02d%02d", gdnow.Hour(), gdnow.Minute(), gdnow.Second(), gdnow.Nanosecond()/1000000, glog.gdate.Year(), glog.gdate.Month(), glog.gdate.Day())

	gerr = os.Rename(glog.groot+"."+glog.gsuffix, glog.groot+"-"+gdyn+"."+glog.gsuffix)
	if gerr != nil {
		glog.Gserr = "Lg 1:Unable to rotate log file:" + glog.groot + "." + glog.gsuffix + "->" + glog.groot + "-" + gdyn + "." + glog.gsuffix
		glog.Gberr = true
		glog.glock.Unlock()
		goto LOGMSG
	}

	gerr = os.Remove(glog.groot + "." + glog.gsuffix)
	if gerr != nil {
		glog.Gserr = "Lg 2:Unable to remove root log file:" + glog.groot + "." + glog.gsuffix
		glog.Gberr = true
		glog.glock.Unlock()
		goto LOGMSG
	}

	glog.gdate = gdate
	glog.glock.Unlock()

LOGMSG:

	glog.glock.Lock()

	if glog.gbopen == false {

		//Now open the main log file and keep going
		glog.gfile, gerr = os.OpenFile(glog.groot+"."+glog.gsuffix, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
		if gerr != nil {
			glog.Gserr = "Lg 3:" + gerr.Error()
			glog.Gberr = true
			glog.glock.Unlock()
			return
		}
		glog.gbopen = true
	}

	gsmsg := fmt.Sprintf("%d/%02d/%02d %02d:%02d:%02d:%03d:%v\r\n", gdnow.Year(), gdnow.Month(), gdnow.Day(), gdnow.Hour(), gdnow.Minute(), gdnow.Second(), gdnow.Nanosecond()/1000000, gmsg)
	gbmsg := []byte(gsmsg)

	_, gerr = glog.gfile.Write(gbmsg)
	if gerr != nil {
		glog.Gserr = "Lg 4:" + gerr.Error()
		glog.Gberr = true
	}
	glog.glock.Unlock()

}
