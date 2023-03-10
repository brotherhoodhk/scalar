package basic

import "strconv"

// 找到一句话中所有数字
func ScanNumber(target string) []int {
	targetbytes := []byte(target)
	resbytes := []byte{}
	allnum := []int{}
	start := -1
	var num int
	for k, v := range targetbytes {
		if v > 47 && v < 58 {
			if k != start+1 && len(resbytes) > 0 {
				num, _ = strconv.Atoi(string(resbytes))
				allnum = append(allnum, num)
				resbytes = []byte{}
			}
			start = k
			resbytes = append(resbytes, v)
		}
	}
	if len(resbytes) > 0 {
		num, _ = strconv.Atoi(string(resbytes))
		allnum = append(allnum, num)
	}
	return allnum
}
func ScanFloat(target string) []float64 {
	targetbytes := []byte(target)
	resbytes := []byte{}
	allnum := []float64{}
	start := -1
	breaddown := 0
	var num float64
	for k, v := range targetbytes {
		if v > 47 && v < 58 {
			if breaddown == 0 {
				breaddown = -1
			}
			if k != start+1 && len(resbytes) > 0 {
				num, _ = strconv.ParseFloat(string(resbytes), 10)
				allnum = append(allnum, num)
				resbytes = []byte{}
			}
			start = k
			resbytes = append(resbytes, v)
		} else if v == 46 && breaddown == -1 && len(resbytes) > 0 {
			resbytes = append(resbytes, v)
			breaddown = 1
			start = k
		} else if v == 46 && breaddown == 1 {
			if resbytes[len(resbytes)-1] == 46 {
				num, _ = strconv.ParseFloat(string(resbytes[:len(resbytes)-1]), 10)
				allnum = append(allnum, num)
			} else {
				num, _ = strconv.ParseFloat(string(resbytes), 10)
				allnum = append(allnum, num)
			}
			resbytes = []byte{}
			breaddown = 0
		}
	}
	if len(resbytes) > 0 {
		if resbytes[len(resbytes)-1] == 46 {
			num, _ = strconv.ParseFloat(string(resbytes[:len(resbytes)-1]), 10)
			allnum = append(allnum, num)
		} else {
			num, _ = strconv.ParseFloat(string(resbytes), 10)
			allnum = append(allnum, num)
		}
	}
	return allnum
}
