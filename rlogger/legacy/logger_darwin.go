package rlogger


#import <Cocoa/Cocoa.h>

//External linkage to communicate back into go
extern void gappcb(int gnwhat) ;



func (glog *Tlogger) Setpath() {
	glog.gpath = "/syslog/"
}
