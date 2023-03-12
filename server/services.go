package server

import (
	"fmt"
	"scal/basic"
	"scal/gocachedriver"
	"scal/mainbody"
	"strconv"
)

func DefaultWordCount(msg *Message) (err error) {
	zone, filepath := "", ""
	zoneid := ""
	if len(msg.Key) > 0 {
		filepath = msg.Key
		if len(msg.Zone) > 0 {
			//若指定了zone
			zone = msg.Zone
			zoneid, err = getzoneid(zone)
		} else {
			//没指定则就用文件名做zone
			zone = getfilename(filepath)
			zoneid, err = getzoneid(zone)
		}
		resmap := mainbody.CountWords(filepath)
		savewordcounts(resmap, zoneid)
	} else {
		err = fmt.Errorf("filepath is empty")
	}
	return
}
func CustomWordCount(msg *Message) (err error) {
	zone, filepath := "", ""
	usemethodname := ""
	if len(msg.Key) > 0 && len(msg.Value) > 0 {
		zoneid := ""
		filepath = msg.Key
		usemethodname = string(msg.Value)
		if len(msg.Zone) > 0 {
			//若指定了zone
			zone = msg.Zone
			zoneid, err = getzoneid(zone)
		} else {
			//没指定则就用文件名做zone
			zone = getfilename(filepath)
			zoneid, err = getzoneid(zone)
		}
		if usefunc, ok := basic.WordCountFunc[usemethodname]; ok {
			resmap := usefunc(filepath)
			if err == nil {
				savewordcounts(resmap, zoneid)
			}
		} else {
			err = fmt.Errorf("wordcount method %v is not exist", usemethodname)
		}
	} else {
		err = fmt.Errorf("filepath or wordcount method name is empty")
	}
	return
}
func getfilename(filepath string) (filename string) {

	return
}
func savewordcounts(origin_data map[string]int, zoneid string) {
	var err error
	for k, v := range origin_data {
		err = gocachedriver.SetKey(k, strconv.Itoa(v), zoneid)
		if err != nil {
			errorlog.Println(err)
		}
	}
	return
}
func getzoneid(zone string) (zoneid string, err error) {
	var ok bool
start:
	zoneid, ok = gocachedriver.CheckZone(zone)
	if !ok {
		err = gocachedriver.CreateZone(zone)
		if err == nil {
			goto start
		}
	}
	return
}
