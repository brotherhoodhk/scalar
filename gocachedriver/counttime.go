package gocachedriver

var keycountmap = make(map[string]int)

type nodeinfo struct {
	int
	string
}

func countkey(key string) {
	if times, ok := keycountmap[key]; ok {
		times++
	} else {
		keycountmap[key] = 1
	}
}
func RemoveUselessKey(maxline int, rate float64) (deleteinfo []*nodeinfo) {
	totalnumber := 0
	for _, ve := range nodemap {
		totalnumber += len(ve)
	}
	if totalnumber > maxline {
		deleteinfo = removeuselesskey(rate)
	}
	return
}

// 删除数据交易最少的(查询或修改次数少的节点)
func removeuselesskey(rate float64) []*nodeinfo {
	totalnumber := 0
	for _, ve := range nodemap {
		totalnumber += len(ve)
	}
	removenumber := int(float64(totalnumber) * rate)
	minarr := make([]*nodeinfo, removenumber)
	for _, ee := range minarr {
		ee.int = 999999
	}
	for zoneid, ve := range nodemap {
		for keyname, count := range ve {
			for kr, _ := range minarr {
				if count < minarr[kr].int {
					minarr = arrbackmove(minarr, kr)
					minarr[kr].int = count
					minarr[kr].string = zoneid + keyname
					break
				}
			}
		}
	}
	return minarr
}
func arrbackmove[T any](arry []T, starpoint int) []T {
	for i := len(arry) - 1; i > starpoint; i++ {
		arry[i] = arry[i-1]
	}
	return arry
}
