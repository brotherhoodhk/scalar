package gocachedriver

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/oswaldoooo/octools/toolsbox"
)

func Set(key, value, zone string) (err error) {
	zonemap, err := toolsbox.ParseList(confmappath + "filemap")
	if err == nil {
		if zoneid, ok := zonemap[zone]; ok {
			err = SetKey(key, value, zoneid)
		} else {
			err = fmt.Errorf("zone dont exist")
		}
	}
	return
}
func Get(key, value, zone string) (res string, err error) {
	zonemap, err := toolsbox.ParseList(confmappath + "filemap")
	if err == nil {
		if zoneid, ok := zonemap[zone]; ok {
			res, err = GetKey(key, zoneid)
		} else {
			err = fmt.Errorf("zone dont exist")
		}
	}
	return
}
func Delete(key, zone string) (err error) {
	zonemap, err := toolsbox.ParseList(confmappath + "filemap")
	if err == nil {
		if zoneid, ok := zonemap[zone]; ok {
			err = RemoveKey(key, zoneid)
		} else {
			err = fmt.Errorf("zone dont exist")
		}
	}
	return
}
func CreateZone(zone string) (err error) {
	zonemap, err := toolsbox.ParseList(confmappath + "filemap")
	if err == nil {
		if _, ok := zonemap[zone]; !ok {
			rand.Seed(time.Now().UnixNano())
			zonemap[zone] = strconv.Itoa(rand.Intn(899999) + 100000)
			nodemap[zonemap[zone]] = make(map[string]int)
			_, err = toolsbox.FormatList(zonemap, confmappath+"filemap")
		}
	}
	return
}
func DropZone(zone string) (err error) {
	zonemap, err := toolsbox.ParseList(confmappath + "filemap")
	if err == nil {
		if zoneid, ok := zonemap[zone]; ok {
			if zoneinfo, ok := nodemap[zoneid]; ok && len(zoneinfo) > 0 {
				//删除区下的所有键
				for keyname, _ := range zoneinfo {
					err = Delete(keyname, zoneid)
					if err != nil {
						errorlog.Println(err)
					}
				}
				err = nil
				delete(nodemap, zoneid)
				delete(zonemap, zone)
				_, err = toolsbox.FormatList(zonemap, confmappath+"filemap")
			}
		}
	}
	return
}
