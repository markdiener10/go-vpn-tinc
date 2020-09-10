#ifndef LOGGEROSX
#define LOGGEROSX

#import <Cocoa/Cocoa.h>

//External linkage to communicate back into go
extern void gappcb(int gnwhat) ;

@interface AppDelegate : NSObject <NSApplicationDelegate>
@end


#endif
