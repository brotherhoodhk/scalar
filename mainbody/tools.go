package mainbody

func removepunctuate(target string) []string {
	bytearr := []byte(target)
	resstr := []string{}
	var buff = []byte{}
	for _, v := range bytearr {
		if v == 44 || v == 46 || v == 33 || v == 63 || v == 32 {
			//为标点符号
			if len(buff) > 0 {
				resstr = append(resstr, string(buff))
				buff = []byte{}
			}
		} else {
			buff = append(buff, v)
		}
	}
	if len(buff) > 0 {
		resstr = append(resstr, string(buff))
	}
	return resstr
}
