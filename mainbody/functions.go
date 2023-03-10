package mainbody

import (
	"fmt"
	"scal/filetools"
	"strings"
)

func CountWords(filepath string) map[string]int {
	content, err := filetools.ScanFile(filepath)
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
