package main

import (
	"fmt"
	"scal/mainbody"
)

func main() {
	countres := mainbody.CountWords("/Users/oswaldo/dev/golang/scal/test.txt")
	if countres != nil {
		fmt.Println(countres)
	}
}
