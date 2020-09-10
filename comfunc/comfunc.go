package comfunc

import (
	"encoding/binary"
	"math/rand"
	"os"
	"strconv"
	"time"

	"crypto/sha256"
	"encoding/base64"
)

func F64tob(gparm uint64) []byte {

	glen := binary.Size(gparm)
	gbyte := make([]byte, glen)
	binary.BigEndian.PutUint64(gbyte, gparm)
	return gbyte
}

func F32tob(gparm uint32) []byte {

	glen := binary.Size(gparm)
	gbyte := make([]byte, glen)
	binary.BigEndian.PutUint32(gbyte, gparm)
	return gbyte
}

func Fbto64(gbyte []byte) uint64 {
	gval := binary.BigEndian.Uint64(gbyte)
	return gval
}

func Fbto32(gbyte []byte) uint32 {
	gval := binary.BigEndian.Uint32(gbyte)
	return gval
}

func Strpad(gstr string, glen int, gpad string) string {

	gn := len(gstr)
	if gn >= glen {
		return gstr
	}
	for {
		gstr = gpad + gstr
		if len(gstr) < glen {
			continue
		}
		break
	}
	return gstr
}

func Strstrip(gstr string, gpad string) string {

	var gres string = ""
	var glen int = len(gstr)
	var gchar string

	for gcnt := 0; gcnt < glen; gcnt++ {
		gchar = string(gstr[gcnt])
		if gchar == gpad {
			continue
		}
		gres += gchar
	}
	return gres
}

func Strkey(gval int64, glen int, gpad string) string {
	gs := strconv.FormatInt(gval, 36)
	gs = Strpad(gs, glen, gpad)
	return gs
}

func Filevalid(gfiledir string) bool {

	_, gerr := os.Stat(gfiledir)
	if gerr == nil {
		return true
	}
	return false
}

//Milliseconds since 1970
func Msec() int64 {
	return (time.Now().UnixNano() / 1000000)
}

//Seconds since 1970 (Class unix stamp)
func Mtik() int64 {
	return time.Now().Unix()
}

//Guaranteed not to be zero
func Mrand(nhi int64) int64 {
	return (rand.Int63n(nhi-1) + 1)
}

func Tagloc(xml string, tag string, noffset int) int {

	gntag := len(tag)
	if gntag == 0 {
		return 0
	}

	gnxml := len(xml)
	if gnxml < gntag {
		return -1
	}

	if gnxml == gntag {
		if xml == tag {
			return 0
		}
	}

	if noffset < 0 {
		noffset = 0
	}

	var gbfnd bool = false

	for ga := noffset; ga < gnxml-gntag; ga++ {

		if xml[ga] != tag[0] {
			continue
		}

		gbfnd = false

		for gb := 1; gb < gntag; gb++ {

			if xml[ga+gb] != tag[gb] {
				gbfnd = true
				break
			}
		}

		if gbfnd == true {
			continue
		}

		return ga

	}
	return -1
}

func Tag(xml string, taga string, tagb string, noffset int) (int, string) {

	gna := len(taga)
	gnb := len(tagb)

	if gna == 0 {
		return -1, ""
	}
	if gnb == 0 {
		return -2, ""
	}

	if noffset < 0 {
		noffset = 0
	}

	ga := Tagloc(xml, taga, noffset)
	if ga == -1 {
		return -3, ""
	}
	gb := Tagloc(xml, tagb, ga+1)
	if gb == -1 {
		return -4, ""
	}

	gtag := string(xml[(ga + gna):gb])

	return (ga + gna), gtag

}

func Diff(ga int64, gb int64) int64 {
	if ga > gb {
		return (ga - gb)
	}
	return (gb - ga)
}

//Thread safe way to copy string to byte
func Strtobyte(gssrc string, gbdst []byte, gnlen int) {

	var gb []byte = []byte(gssrc)
	var gblen int = len(gb)

	if gblen == 0 {
		if gnlen > 0 {
			gbdst[0] = 0
		}
		return
	}

	if gblen < gnlen {
		gbdst[gblen] = 0
		gnlen = gblen
	}

	for gncnt := 0; gncnt < gnlen; gncnt++ {
		gbdst[gncnt] = gb[gncnt]
	}

}

func Hashb(ginput []byte) string {
	var ghash [32]byte = sha256.Sum256(ginput)
	return base64.StdEncoding.EncodeToString(ghash[:])
}

func Hashs(ginput string) string {
	var ghash [32]byte = sha256.Sum256([]byte(ginput))
	return base64.StdEncoding.EncodeToString(ghash[:])
}

func B64S(ginput []byte) string {
	return base64.StdEncoding.EncodeToString(ginput[:])
}

func SB64(ginput string) (error, []byte) {
	var gerr error
	var gbres []byte
	gbres, gerr = base64.StdEncoding.DecodeString(ginput)
	return gerr, gbres
}
