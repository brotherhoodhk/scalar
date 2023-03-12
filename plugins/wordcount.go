package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

var wrbuffsize = 10 * MB

func CountWords(filepath string) map[string]int {
	content, err := ScanFile(filepath)
	wordcount := make(map[string]int)
	if err != nil {
		fmt.Println(err)
		return nil
	} else {
		origin_words := string(content)
		origin_words_arr := strings.Split(origin_words, "\n")
		for _, v := range origin_words_arr {
			if len(v) > 0 {
				vearr := removepunctuate(v)
				if len(vearr) > 0 {
					for _, ve := range vearr {
						if _, ok := wordcount[ve]; ok {
							wordcount[ve] = wordcount[ve] + 1
						} else {
							wordcount[ve] = 1
						}
					}
				}
			}
		}
	}
	return wordcount
}

// 移除掉标点符号等杂质，筛选出单词
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

// 扫描文档进内存
func ScanFile(filepath string) ([]byte, error) {
	fe, err := os.OpenFile(filepath, os.O_RDONLY, 0700)
	if err != nil {
		return nil, err
	}
	defer fe.Close()
	buff := make([]byte, wrbuffsize)
	read := bufio.NewReader(fe)
	lang, err := read.Read(buff)
	if err != nil {
		return nil, err
	}
	return buff[:lang], nil
}
