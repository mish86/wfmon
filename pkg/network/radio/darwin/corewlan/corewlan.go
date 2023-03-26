//go:build darwin

package corewlan

// Note: preamble should be immediately folowed by import of "C"

/*
#cgo CFLAGS: -x objective-c -Wno-incompatible-pointer-types
#cgo LDFLAGS: -lobjc -framework Cocoa -framework Foundation -framework CoreWLAN
#import <Cocoa/Cocoa.h>
#import <Foundation/Foundation.h>
#import <CoreWLAN/CoreWLAN.h>
#import <string.h>
#import <stdlib.h>

typedef struct _NetworkData {
    NSString *bssid;
    NSString *ssid;
	CWChannel *channel;
} NetworkData;

typedef struct _Result {
	bool status;
	NSError *err;
} Result;

static inline
NetworkData * cwnetworkdata(CWInterface *iface)
{
    NetworkData *d = malloc(sizeof(NetworkData));
    d->bssid = [iface bssid];
    d->ssid = [iface ssid];
	d->channel = [iface wlanChannel];
    return d;
}

static inline
int channelNumber(CWChannel *channel)
{
	if (channel == NULL) { return 0; }
	return [channel channelNumber];
}

static inline
int resultErrorCode(Result *r)
{
	if (r == NULL) { return -1; }

	NSError *err = r->err;
	if (err == NULL) { return 0; }

	return [err code];
}

static inline
const char * resultErrorDomain(Result *r)
{
	if (r == NULL) { return ""; }

	NSError *err = r->err;
	if (err == NULL) { return 0; }

	NSString *domain = [err domain];
	char *out = strdup([domain UTF8String]);

	return out;
}

static inline
const char * getDefaultInterface()
{
	CWWiFiClient *client = [CWWiFiClient sharedWiFiClient];
    CWInterface *iface = [client interface];
	if (!iface) {
		return "";
	}

	NSString *ifaceName = [iface interfaceName];
	char *out = strdup([ifaceName UTF8String]);

	return out;
}

static inline
NetworkData * getAssociatedNetwork(char *cIFaceName)
{
	NSString *ifaceName = [[NSString alloc] initWithUTF8String:cIFaceName];

	CWWiFiClient *client = [CWWiFiClient sharedWiFiClient];
    CWInterface *iface = [client interfaceWithName:ifaceName];
	if (!iface) {
		return NULL;
	}

	return cwnetworkdata(iface);
}

static inline
void disassociateNetwork(char *cIFaceName)
{
	NSString *ifaceName = [[NSString alloc] initWithUTF8String:cIFaceName];

	CWWiFiClient *client = [CWWiFiClient sharedWiFiClient];
    CWInterface *iface = [client interfaceWithName:ifaceName];
	if (!iface) {
		return;
	}

	[iface disassociate];
}

static inline
NSArray * getSupportedChannels(char *cIFaceName)
{
	NSString *ifaceName = [[NSString alloc] initWithUTF8String:cIFaceName];

	CWWiFiClient *client = [CWWiFiClient sharedWiFiClient];
    CWInterface *iface = [client interfaceWithName:ifaceName];
	if (!iface) {
		return @[];
	}

	NSSet<CWChannel *> *channels = [iface supportedWLANChannels];
	NSMutableArray *channelsArr = [[NSMutableArray alloc] initWithCapacity:[channels count]];
	for(CWChannel *channel in channels) {
		[channelsArr addObject:[NSNumber numberWithInt:[channel channelNumber]]];
	}

	return channelsArr;
}

static inline
Result* setInterfaceChannel(char *cIFaceName, int c)
{
	Result *r = malloc(sizeof(Result));
    r->status = false;
    r->err = NULL;

	NSString *ifaceName = [[NSString alloc] initWithUTF8String:cIFaceName];

	CWWiFiClient *client = [CWWiFiClient sharedWiFiClient];
    CWInterface *iface = [client interfaceWithName:ifaceName];
	if (!iface) {
		return r;
	}

	NSSet<CWChannel *> *channels = [iface supportedWLANChannels];
	NSPredicate *predicate = [NSPredicate predicateWithFormat:@"SELF.channelNumber == %d", c];
	NSSet<CWChannel *> *filteredChannels = [channels filteredSetUsingPredicate:predicate];
	CWChannel *channel = [filteredChannels anyObject];
	if (!channel) {
		return r;
	}

	NSError *anyerror;
	r->status = [iface setWLANChannel:channel error: &anyerror];
	r->err = anyerror;

	return r;
}

static inline
const char* nsstring2cstring(NSString *s) {
    if (s == NULL) { return NULL; }
    const char *cstr = [s UTF8String];
    return cstr;
}

static inline
int nsnumber2int(NSNumber *i) {
    if (i == NULL) { return 0; }
    return i.intValue;
}


static inline
unsigned long nsarraycount(NSArray *arr) {
    if (arr == NULL) { return 0; }
    return [arr count];
}

static inline
void* nsarrayitem(NSArray *arr, unsigned long i) {
    if (arr == NULL) { return NULL; }
    return [arr objectAtIndex:i];
}

*/
import "C"
import (
	"fmt"
	"sort"
	"unsafe"
	"wfmon/pkg/network"
)

// NSString -> C string
func cstring(s *C.NSString) *C.char { return C.nsstring2cstring(s) }

// NSString -> Go string
func gostring(s *C.NSString) string { return C.GoString(cstring(s)) }

// NSNumber -> Go int
func goint(i *C.NSNumber) int { return int(C.nsnumber2int(i)) }

// NSArray count
func nsarraycount(arr *C.NSArray) uint { return uint(C.nsarraycount(arr)) }

// NSArray item
func nsarrayitem(arr *C.NSArray, i uint) unsafe.Pointer {
	return C.nsarrayitem(arr, C.ulong(i))
}

//nolint:nonamedreturns // ignore
func GetDefaultWiFiInterface() (iface string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to get default interface, %v", r)
		}
	}()

	cIFaceName := C.getDefaultInterface()
	defer C.free(unsafe.Pointer(cIFaceName))

	iface = C.GoString(cIFaceName)
	return iface, err
}

// Location should be active for the app. Otherwise BSSID is null.
//
//nolint:nonamedreturns // ignore
func GetAssociatedNetwork(ifaceName string) (net network.Network, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to get associated network for interface %s, %v", ifaceName, r)
		}
	}()

	cIFaceName := C.getDefaultInterface()
	defer C.free(unsafe.Pointer(cIFaceName))

	ptrNetworkData := C.getAssociatedNetwork(cIFaceName)
	defer C.free(unsafe.Pointer(ptrNetworkData))

	networkData := *ptrNetworkData
	bssid := gostring(networkData.bssid)
	ssid := gostring(networkData.ssid)
	channel := C.channelNumber(networkData.channel)

	net = network.Network{
		BSSID:   bssid,
		SSID:    ssid,
		Channel: uint8(channel),
	}
	return net, err
}

// Does not associate interface to network back after program exit.
func DisassociateFromNetwork(ifaceName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to disassociated interface %s from network, %v", ifaceName, r)
		}
	}()

	cIFaceName := C.getDefaultInterface()
	defer C.free(unsafe.Pointer(cIFaceName))

	C.disassociateNetwork(cIFaceName)

	return err
}

//nolint:nonamedreturns // ignore
func GetSupportedChannels(ifaceName string) (channels []int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to get supported channels for interface %s, %v", ifaceName, r)
		}
	}()

	cIFaceName := C.CString(ifaceName)
	defer C.free(unsafe.Pointer(cIFaceName))

	nsArrChannels := C.getSupportedChannels(cIFaceName)
	defer C.free(unsafe.Pointer(nsArrChannels))

	cnt := nsarraycount(nsArrChannels)
	chSet := make(map[int]struct{}, cnt)

	for i := uint(0); i < cnt; i++ {
		nsNumChannel := (*C.NSNumber)(nsarrayitem(nsArrChannels, i))
		chSet[goint(nsNumChannel)] = struct{}{}
	}

	channels = make([]int, len(chSet))
	{
		var i int
		for ch := range chSet {
			channels[i] = ch
			i++
		}
	}
	sort.Ints(channels)

	return channels, err
}

func SetInterfaceChannel(ifaceName string, channel int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to set channel %d for interface %s, %v", channel, ifaceName, r)
		}
	}()

	cIFaceName := C.CString(ifaceName)
	defer C.free(unsafe.Pointer(cIFaceName))

	ptrResult := C.setInterfaceChannel(cIFaceName, C.int(channel))
	defer C.free(unsafe.Pointer(ptrResult))

	if ptrResult.status {
		return err
	}

	if ptrResult.err == nil {
		err = fmt.Errorf("unknown error while setting interface %s channel %d", ifaceName, channel)
		return err
	}

	code := C.resultErrorCode(ptrResult)
	cDomain := C.resultErrorDomain(ptrResult)
	domain := C.GoString(cDomain)

	err = fmt.Errorf(
		"error (code: %d, domain: %s) while setting interface %s channel %d",
		code,
		domain,
		ifaceName,
		channel,
	)

	return err
}
