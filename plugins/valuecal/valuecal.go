package main

// 价值计算法,算出结果不做操作
func ValueCal(origin_data map[string][]byte) (res []struct {
	Key   string
	Value []byte
}) {
	//按照key长度进行赋值，长度越长，值越小
	mapbuff := make(map[string]struct{})
	max := struct {
		Key   string
		Value int
	}{Key: "", Value: 0}
	for i := 0; i < len(origin_data); i++ {
		for k, _ := range origin_data {
			if _, ok := mapbuff[k]; len(k) > max.Value && !ok {
				max = struct {
					Key   string
					Value int
				}{k, len(k)}
			}
		}
		res = append(res, struct {
			Key   string
			Value []byte
		}{Key: max.Key, Value: origin_data[max.Key]})
		delete(origin_data, max.Key)
	}
	return
}

// 算出结果并执行操作
func GetFinalLeavet(origin_data map[string][]byte, drop_rate float64) (res map[string][]byte) {
	resrank := ValueCal(origin_data)
	drop_num := int(float64(len(resrank)) * drop_rate)
	resrank = resrank[drop_num:]
	for _, ve := range resrank {
		res[ve.Key] = ve.Value
	}
	return
}
