#import "NSFileManager+DirectoryLocations.h"

enum
{
	DirectoryLocationErrorNoPathFound,
	DirectoryLocationErrorFileExistsAtLocation
};
	
NSString * const DirectoryLocationDomain = @"DirectoryLocationDomain";

//
// applicationSupportDirectory
//
// Returns the path to the applicationSupportDirectory (creating it if it doesn't
// exist).
//
- (NSString *)applicationSupportDirectory
{
	NSString *executableName =
		[[[NSBundle mainBundle] infoDictionary] objectForKey:@"CFBundleExecutable"];
	NSError *error;
	NSString *result =
		[self
			findOrCreateDirectory:NSApplicationSupportDirectory
			inDomain:NSUserDomainMask
			appendPathComponent:executableName
			error:&error];
	if (!result)
	{
		NSLog(@"Unable to find or create application support directory:\n%@", error);
	}
	return result;
}

@end
