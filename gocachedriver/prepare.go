package gocachedriver

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	driver_tools "github.com/oswaldoooo/gocache-driver/basics"
	"github.com/oswaldoooo/octools/toolsbox"
)

var ROOTPATH = os.Getenv("SCALAR_HOME")
var confmappath = ROOTPATH + "/conf/"
var isinit = false
var nodemap = make(map[string]map[string]int) //记录每个节点操作次数,zoneid>>{keyname;use count}
type siteconfig struct {
	XMLName   xml.Name    `xml:"scalar"`
	Cacheinfo gocacheinfo `xml:"gocache"`
}
type gocacheinfo struct {
	XMLName    xml.Name `xml:"gocache"`
	Host       string   `xml:"hostadd"`
	Port       int      `xml:"port"`
	Default_DB string   `xml:"default_db"`
}

var errorlog = toolsbox.LogInit("error", ROOTPATH+"/logs/error.log")

func init() {
	buff, err := ioutil.ReadFile(confmappath + "site.xml")
	if err == nil {
		siteconf := new(siteconfig)
		err = xml.Unmarshal(buff, siteconf)
		if err == nil && len(siteconf.Cacheinfo.Host) > 0 && len(siteconf.Cacheinfo.Default_DB) > 0 && siteconf.Cacheinfo.Port > 0 {
			dbcon = driver_tools.New(siteconf.Cacheinfo.Host, siteconf.Cacheinfo.Port, "", siteconf.Cacheinfo.Default_DB)
			err = dbcon.Connect()
			if err == nil {
				isinit = true
			}
		}
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	//解析filemap to nodemap
	filemap, err := toolsbox.ParseList(confmappath + "filemap")
	if err == nil {
		fearr, err := ioutil.ReadDir(ROOTPATH + "/dblink")
		if err == nil {
			namemap := make(map[string]int64)
			for _, ve := range fearr {
				namemap[ve.Name()] = ve.Size()
			}
			for _, zoneid := range filemap {
				if _, ok := namemap[zoneid]; !ok {
					nodemap[zoneid] = make(map[string]int)
				} else {
					if namemap[zoneid] > 6 {
						//如果有数据记录
						be, err := ioutil.ReadFile(ROOTPATH + "/dblink/" + zoneid)
						tmap := make(map[string]int)
						if err == nil {
							err = json.Unmarshal(be, &tmap)
						}
						nodemap[zoneid] = tmap
					} else {
						nodemap[zoneid] = make(map[string]int)
					}
				}
			}
		}
	}
	keymap, err := dbcon.GetAllKeys()
	if err == nil {
		comparedatastore(keymap)
	}
	if err != nil {
		fmt.Println(err)
	}
}
func comparedatastore(origin_key map[string][]byte) {
	for k, vearr := range nodemap {
		for kname, _ := range vearr {
			name := k + kname
			if _, ok := origin_key[name]; !ok {
				//在gocache中不存在,移除本地记录
				delete(vearr, kname)
			}
		}
	}
}
