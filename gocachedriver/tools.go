package gocachedriver

import (
	"github.com/oswaldoooo/octools/toolsbox"
)

func CheckZone(zone string) (zoneid string, res bool) {
	if len(zone) < 1 {
		res = false
	} else {
		content, err := toolsbox.ParseList(confmappath + "filemap")
		if id, ok := content[zone]; err == nil && ok {
			zoneid = id
			res = true
			return
		}
	}
	res = false
	return
}
