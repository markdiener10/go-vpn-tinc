package rjson

import (
	"bytes"
	"encoding/json"
	_ "encoding/xml"
	"errors"
	"fmt"
	_ "io"
	"io/ioutil"
	_ "reflect"
	"strconv"
	"strings"
)

type jsonmode int

const (
	EMPTY jsonmode = iota
	NULL
	BOOLEAN
	NUMBER
	FLOAT
	STRING
	ARRAY
	OBJECT
	TOPOBJECT
	TOPARRAY
	UNKNOWN
)

type Tjson struct {
	Nmode jsonmode
	nodes []Tjson
	Gkey  string
	Gstr  string
	Gbool bool
	Gnum  int64
	Gflt  float64
}

func (gjson *Tjson) Clear(gnval jsonmode) {

	var gchild *Tjson

	for gidx := range gjson.nodes {
		gchild = &gjson.nodes[gidx]
		gchild.Clear(EMPTY)
	}

	gjson.Gbool = false
	gjson.Gkey = ""
	gjson.Gnum = 0
	gjson.Gstr = ""
	gjson.Gflt = 0.0
	gjson.Nmode = gnval
	gjson.nodes = make([]Tjson, 0)
}

func (gjson *Tjson) Loadstr(sjson string) error {
	bjson := []byte(sjson)
	return (gjson.Loadbyte(bjson))
}

func (gjson *Tjson) Loadbyte(gbjson []byte) (gerr error) {

	defer func() {
		gerr := recover()
		if gerr != nil {
			fmt.Println("Panic:", gerr)
		}
	}()

	//var gidx int
	var gval interface{}

	var gbase interface{}

	gdecode := json.NewDecoder(bytes.NewReader(gbjson))
	gdecode.UseNumber()

	//gerr = json.Unmarshal(bjson, &gbase)
	gerr = gdecode.Decode(&gbase)
	if gerr != nil {
		return gerr
	}

	//fmt.Println("Base:", gbase)
	gjson.Clear(EMPTY)

	//Cast the data into a specific map type for the highest level
	switch gbase.(type) {
	case []interface{}:

		gjson.Nmode = TOPARRAY

		gtop := gbase.([]interface{})

		if len(gtop) == 0 {
			//Handle Empty JSON
			return nil
		}

		for _, gval = range gtop {

			//fmt.Println("JSON Array Loop Value:", gidx, gval)
			gerr = gjson.Add("", gval)
			if gerr != nil {
				//fmt.Println("JSON Loop Error:", gkey, " Value:")
				return gerr
			}
			//fmt.Println("JSON Loop Good:", gkey, " Value:", gval)
		}

		//fmt.Println("TOPMAP:", gtop)

		break

	case map[string]interface{}:

		gjson.Nmode = TOPOBJECT

		gtop := gbase.(map[string]interface{})

		if len(gtop) == 0 {
			//Handle Empty JSON
			return nil
		}
		for gkey, gval := range gtop {
			//fmt.Println("JSON Loop Key:", gkey, " Value:", gval)
			gerr = gjson.Add(gkey, gval)
			if gerr != nil {
				//fmt.Println("JSON Loop Error:", gkey, " Value:")
				return gerr
			}
			//fmt.Println("JSON Loop Good:", gkey, " Value:", gval)
		}

		//fmt.Println("TOPMAP:", gtop)

		break
	default:
		return errors.New("Unknown JSON type")
	}

	return nil

}

func (gjson *Tjson) Add(gkey string, gval interface{}) error {

	//fmt.Println("JSON Add Key:", gkey, gval)
	var gerr error

	//var gchild *Tjson
	switch gval.(type) {
	case nil:

		//fmt.Println(gkey, "is nil ", gval)
		gchild := Tjson{}
		gchild.Clear(EMPTY)
		gchild.Gkey = gkey
		gchild.Nmode = NULL
		gjson.nodes = append(gjson.nodes, gchild)

	case bool:

		//fmt.Println(gkey, "is boolean ", gval)

		gchild := Tjson{}
		gchild.Clear(EMPTY)
		gchild.Nmode = BOOLEAN
		gchild.Gkey = gkey
		gchild.Gbool = gval.(bool)
		gjson.nodes = append(gjson.nodes, gchild)

	case string:

		//fmt.Println(gkey, "is string ", gval)

		gchild := Tjson{}
		gchild.Clear(EMPTY)
		gchild.Nmode = STRING
		gchild.Gkey = gkey
		gchild.Gstr = gval.(string)
		gjson.nodes = append(gjson.nodes, gchild)

	case json.Number:

		//fmt.Println(gkey, "is Number", gval)

		gchild := Tjson{}
		gchild.Clear(EMPTY)
		gchild.Gkey = gkey
		gchild.Gstr = gval.(json.Number).String()
		if strings.Contains(gchild.Gstr, ".") {
			gchild.Gflt, gerr = gval.(json.Number).Float64()
			gchild.Nmode = FLOAT
		} else {
			gchild.Gnum, gerr = gval.(json.Number).Int64()
			gchild.Nmode = NUMBER
		}
		if gerr != nil {
			return gerr
		}
		gjson.nodes = append(gjson.nodes, gchild)

	case []interface{}:

		//fmt.Println(gkey, "is an array: ", gval)

		gchild := Tjson{}
		gchild.Clear(EMPTY)
		gchild.Nmode = ARRAY
		gchild.Gkey = gkey

		for _, gaval := range gval.([]interface{}) {

			//fmt.Println("ARRAY Loop Key:", gkey, " Value:", gaval, "---------------")

			//Array does not have a "KEY"
			gerr := gchild.Add("", gaval)
			if gerr != nil {
				return gerr
			}
		}

		gjson.nodes = append(gjson.nodes, gchild)

	case map[string]interface{}:

		//fmt.Println(gkey, "is an object:", gval)

		gchild := Tjson{}
		gchild.Clear(EMPTY)
		gchild.Nmode = OBJECT
		gchild.Gkey = gkey

		for gokey, goval := range gval.(map[string]interface{}) {

			//fmt.Println("OBJECT Loop Key:", gokey, " Value:", goval, "---------------")
			gerr := gchild.Add(gokey, goval)
			if gerr != nil {
				return gerr
			}
			//fmt.Println("OBJECT Loop Key Done---------", gokey, " Value:", goval, "---------------")
		}

		gjson.nodes = append(gjson.nodes, gchild)

	default:

		//fmt.Println(gkey, "is unknown: ", gval)

		gchild := Tjson{}
		gchild.Clear(EMPTY)
		gchild.Nmode = UNKNOWN
		gchild.Gkey = gkey
		//gchild.Gstr = gval
		gjson.nodes = append(gjson.nodes, gchild)

	}

	return nil

}

//Quicker to dump bytes than string
func (gjson *Tjson) Dump() []byte {

	var gresu bytes.Buffer
	gresu.Reset()

	if gjson == nil {
		return gresu.Bytes()
	}

	if gjson.Nmode == TOPARRAY {
		gresu.WriteString("[")
	}
	if gjson.Nmode == TOPOBJECT {
		gresu.WriteString("{")
	}
	var gcomma string = ""
	var gpre string = ""

	for gidx := range gjson.nodes {

		gobj := &gjson.nodes[gidx]

		if len(gobj.Gkey) > 0 {
			gpre = gcomma + `"` + gobj.Gkey + `":`
		} else {
			gpre = gcomma
		}

		if gcomma == "" {
			gcomma = ","
		}

		gresu.WriteString(gpre)

		switch gobj.Nmode {
		case TOPOBJECT:
			continue
		case TOPARRAY:
			continue
		case NULL:
			gresu.WriteString(`null`)
		case BOOLEAN:
			if gobj.Gbool == true {
				gresu.WriteString("true")
				break
			}
			gresu.WriteString("false")
		case NUMBER:
			gresu.WriteString(strconv.FormatInt(gobj.Gnum, 10))
		case FLOAT:
			gresu.WriteString(strconv.FormatFloat(gobj.Gflt, 'g', 6, 64))
		case STRING:
			gresu.WriteString(`"` + gobj.Gstr + `"`)
		case ARRAY:
			gresu.WriteString(`[` + gobj.Dumpstr() + `]`)
		case OBJECT:
			gresu.WriteString(`{` + gobj.Dumpstr() + `}`)
		case UNKNOWN:
			continue
		}
	}

	if gjson.Nmode == TOPARRAY {
		gresu.WriteString("]")
	}
	if gjson.Nmode == TOPOBJECT {
		gresu.WriteString("}")
	}
	return gresu.Bytes()
}

func (gjson *Tjson) Dumpstr() string {
	return (string(gjson.Dump()))
}

func (gjson *Tjson) Find(gkey string, grecurse, gcase, gsubstr bool) *Tjson {

	if gcase == true {
		gkey = strings.ToUpper(gkey)
	}

	var gkeyb string = ""
	var gobj *Tjson = nil
	var gcobj *Tjson = nil
	var gidx int = 0

	for gidx = range gjson.nodes {

		gobj = &gjson.nodes[gidx]

		if gobj.Nmode == EMPTY {
			continue
		}

		gkeyb = gobj.Gkey
		if gcase == true {
			gkeyb = strings.ToUpper(gkeyb)
		}

		if gsubstr == false {
			if gkeyb == gkey {
				return gobj
			}
			continue
		}

		if strings.Contains(gkeyb, gkey) == true {
			return gobj
		}
	}

	if grecurse == false {
		return nil
	}

	for gidx = range gjson.nodes {

		gobj = &gjson.nodes[gidx]

		gcobj = gobj.Find(gkey, grecurse, gcase, gsubstr)
		if gcobj != nil {
			return gcobj
		}
	}
	return nil
}

func (gjson *Tjson) Qfind(gkey string) *Tjson {
	return gjson.Find(gkey, false, true, false)
}

func (gjson *Tjson) Qfindb(gkey string, gnerr *int, gval *bool) {

	var gjp *Tjson = gjson.Qfind(gkey)
	if gjp == nil {
		*gnerr += 1
		return
	}
	*gval = gjp.Gbool
}

func (gjson *Tjson) Qfinds(gkey string, gnerr *int, gval *string) {

	var gjp *Tjson = gjson.Qfind(gkey)
	if gjp == nil {
		*gnerr += 1
		return
	}
	*gval = gjp.Gstr
}

func (gjson *Tjson) Qfindi(gkey string, gnerr *int, gval *int) {

	var gjp *Tjson = gjson.Qfind(gkey)
	if gjp == nil {
		*gnerr += 1
		return
	}
	*gval = int(gjp.Gnum)
}

func (gjson *Tjson) Qfindi64(gkey string, gnerr *int, gval *int64) {

	var gjp *Tjson = gjson.Qfind(gkey)
	if gjp == nil {
		*gnerr += 1
		return
	}
	*gval = gjp.Gnum
}

func (gjson *Tjson) Count() int64 {
	return int64(len(gjson.nodes))
}

func (gjson *Tjson) Idx(gidx int) *Tjson {

	for gcnt, gobj := range gjson.nodes {

		if gcnt != gidx {
			continue
		}
		return &gobj
	}
	return nil
}

func (gjson *Tjson) Str() string {

	switch gjson.Nmode {
	case EMPTY, NULL:
		return ""
	case BOOLEAN:
		if gjson.Gbool == true {
			return "TRUE"
		}
		return "FALSE"
	case NUMBER:
		return strconv.FormatInt(gjson.Gnum, 10)
	case FLOAT:
		return strconv.FormatFloat(gjson.Gflt, 'g', 6, 64)
	case STRING:
		return gjson.Gstr
	case ARRAY:
		return "ARRAY"
	case OBJECT:
		return "OBJECT"
	default:
		return "????"
	}
}

func (gjson *Tjson) Num() int64 {

	var gerr error
	var gb bool
	var gi int64

	switch gjson.Nmode {
	case EMPTY:
		return 0
	case NULL:
		return 0
	case BOOLEAN:

		//fmt.Println("Bool:", gjson.Nmode)

		gb, gerr = strconv.ParseBool(gjson.Gstr)
		if gerr != nil {
			return 0
		}
		if gb == true {
			return 1
		}
		return 0

	case NUMBER:

		return gjson.Gnum

	case FLOAT:

		//fmt.Println("FloatY:", gjson.Nmode, "-", gjson.Gstr, "=")
		return int64(gjson.Gflt)

	case STRING:

		gi, gerr = strconv.ParseInt(gjson.Gstr, 10, 64)
		if gerr != nil {
			return 0
		}
		return gi

	case ARRAY:
		return 0
	case OBJECT:
		return 0
	default:
		return 0
	}
}

func (gjson *Tjson) Addj(gkey string, gval *Tjson) {

	gchild := Tjson{}
	gchild.Clear(EMPTY)
	gchild.Nmode = gval.Nmode
	gchild.Gkey = gval.Gkey
	gchild.Gstr = gval.Gstr
	gchild.Gnum = gval.Gnum
	gchild.Gbool = gval.Gbool
	gchild.Gflt = gval.Gflt

	for _, gval := range gval.nodes {
		gchild.nodes = append(gchild.nodes, gval)
	}

	gjson.nodes = append(gjson.nodes, gchild)

}

func (gjson *Tjson) Addb(gkey string, gval bool) {

	gchild := Tjson{}
	gchild.Clear(EMPTY)
	gchild.Nmode = BOOLEAN
	if gjson.Nmode != TOPARRAY && gjson.Nmode != ARRAY {
		gchild.Gkey = gkey
	}
	gchild.Gbool = gval
	if gval == true {
		gchild.Gnum = 1
		gchild.Gflt = 1
	}
	gjson.nodes = append(gjson.nodes, gchild)
}

func (gjson *Tjson) Adds(gkey, gval string) {
	gchild := Tjson{}
	gchild.Clear(EMPTY)
	gchild.Nmode = STRING
	if gjson.Nmode != TOPARRAY && gjson.Nmode != ARRAY {
		gchild.Gkey = gkey
	}
	gchild.Gstr = gval
	gjson.nodes = append(gjson.nodes, gchild)
}

func (gjson *Tjson) Addi(gkey string, gval int) {
	gchild := Tjson{}
	gchild.Clear(EMPTY)
	gchild.Nmode = NUMBER

	if gjson.Nmode != TOPARRAY && gjson.Nmode != ARRAY {
		gchild.Gkey = gkey
	}

	gchild.Gnum = int64(gval)
	gchild.Gflt = float64(gval)
	gjson.nodes = append(gjson.nodes, gchild)
}

func (gjson *Tjson) Addi64(gkey string, gval int64) {
	gchild := Tjson{}
	gchild.Clear(EMPTY)
	gchild.Nmode = NUMBER

	if gjson.Nmode != TOPARRAY && gjson.Nmode != ARRAY {
		gchild.Gkey = gkey
	}

	gchild.Gnum = int64(gval)
	gchild.Gflt = float64(gval)
	gjson.nodes = append(gjson.nodes, gchild)
}

func (gjson *Tjson) Addf(gkey string, gval float64) {
	gchild := Tjson{}
	gchild.Clear(EMPTY)
	gchild.Nmode = FLOAT

	if gjson.Nmode != TOPARRAY && gjson.Nmode != ARRAY {
		gchild.Gkey = gkey
	}

	gchild.Gflt = gval
	gchild.Gnum = int64(gval)
	gjson.nodes = append(gjson.nodes, gchild)
}

func (gjson *Tjson) Loadfile(gfile string) error {

	gbytes, gerr := ioutil.ReadFile(gfile)
	if gerr != nil {
		return gerr
	}
	gerr = gjson.Loadbyte(gbytes)
	if gerr != nil {
		return gerr
	}
	return nil
}
