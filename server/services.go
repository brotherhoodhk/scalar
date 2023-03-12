package server

import (
	"fmt"
	"scal/basic"
	"scal/gocachedriver"
	"scal/mainbody"
	"strconv"
	"strings"
)

// use default method count words
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

// use custom method count words
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

// get filename from filepath.example: filepath /home/brotherhoodhk/jake filename jake
func getfilename(filepath string) (filename string) {
	if strings.ContainsRune(filepath, '/') {
		patternarr := strings.Split(filepath, "/")
		if len(patternarr) >= 2 {
			if len(patternarr[len(patternarr)-1]) > 0 {
				filename = patternarr[len(patternarr)-1]
			} else {
				filename = patternarr[len(patternarr)-2]
			}
		} else if len(patternarr) == 1 {
			filename = patternarr[0]
		}
	} else {
		filename = filepath
	}
	return
}

// save wordcount result to database
func savewordcounts(origin_data map[string]int, zoneid string) {
	var err error
	for k, v := range origin_data {
		err = gocachedriver.SetKey(k, strconv.Itoa(v), zoneid)
		if err != nil {
			errorlog.Println(err)
		}
	}
	gocachedriver.ForceSave()
	return
}

// get the target zone's zoneid.if dont exist ,it will auto create
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
