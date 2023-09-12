package main

import (
	"net/http"
	_ "net/http/pprof"
	"scal/cmd"
)

func main() {
	go http.ListenAndServe(":9988", nil)
	cmd.Execute()
}
