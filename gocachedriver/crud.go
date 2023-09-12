package gocachedriver

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	driver_tools "github.com/oswaldoooo/gocache-driver/basics"
)

var (
	dbcon    *driver_tools.CacheDB
	MaxLine  = 1000
	Rate     = 0.3
	Conn_Err = driver_tools.CON_ERR
)

func SetKey(key, value string, zoneid string) (err error) {
	if len(zoneid) != 6 {
		err = fmt.Errorf("zoneid's format is not correct")
		return
	}
	key = zoneid + key
	err = dbcon.SetKey(key, value)
	if err == nil {
		NodePlus(key)
	}
	return
}
func GetKey(key string, zoneid string) (res string, err error) {
	if len(zoneid) != 6 {
		err = fmt.Errorf("zoneid's format is not correct")
		return
	}
	key = zoneid + key
	resarr, err := dbcon.GetKeys(key)
	if err == nil {
		NodePlus(key)
	}
	if len(resarr) > 0 {
		res = resarr[0]
	}
	return
}
func GetZoneAllKeys(zoneid string) (resmap map[string][]byte, err error) {
	resmap, err = dbcon.GetKeysContain(zoneid)
	return
}
func RemoveKey(key string, zoneid string) (err error) {
	if len(zoneid) != 6 {
		err = fmt.Errorf("zoneid's format is not correct")
		return
	}
	key = zoneid + key
	if key == "*" {
		//delete all zone keys
		resmap, err := dbcon.GetKeysContain(zoneid)
		if err == nil {
			keysarr := []string{}
			for k, _ := range resmap {
				keysarr = append(keysarr, k)
			}
			if len(keysarr) > 0 {
				err = dbcon.DeleteKeys(keysarr...)
			}
		}
	} else {
		err = dbcon.DeleteKeys(key)
		if err == nil {
			NodeDelete(key)
		}
	}
	return
}
func autoclean() {
	toremovelist := RemoveUselessKey(MaxLine, Rate)
	for _, v := range toremovelist {
		if _, ok := nodemap[v.string[:6]][v.string[6:]]; ok {
			delete(nodemap[v.string[:6]], v.string[6:])
		}
	}
	go SaveInfo()

}

// 这里的key是完整的key，分区号+键名
func NodePlus(key string) {
	if _, ok := nodemap[key]; ok {
		nodemap[key[:6]][key[6:]]++
	} else {
		nodemap[key[:6]][key[6:]] = 1
	}
	go SaveInfo()
}

// key 是完整的key
func NodeDelete(key string) {
	zoneid := key[:6]
	realkey := key[6:]
	if _, ok := nodemap[zoneid][realkey]; ok {
		delete(nodemap[zoneid], realkey)
	}
	SaveInfo()
}

var rwlock sync.RWMutex

// 采用随机数来决定是否储存到硬盘
func SaveInfo() {
	rand.Seed(time.Now().UnixNano())
	var realpath string
	parentdir := ROOTPATH + "/dblink/"
	var err error
	if rand.Intn(3) == 1 {
		for zoneid, content := range nodemap {
			realpath = parentdir + zoneid
			rwlock.Lock()
			fe, err := os.OpenFile(realpath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
			if err == nil {
				defer fe.Close()
				resbytes, err := json.Marshal(&content)
				if err == nil {
					_, err = fe.Write(resbytes)
				}
			}
			rwlock.Unlock()
		}
	}
	if err != nil {
		errorlog.Println(err)
	}
}

// 手动保存
func ForceSave() {
	go SaveInfo()
	err := dbcon.Save()
	if err != nil {
		errorlog.Println(err)
	}
}
func getallkeys() (keysmap map[string][]byte, err error) {
	keysmap, err = dbcon.GetAllKeys()
	return
}
