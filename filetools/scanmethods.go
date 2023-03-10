package filetools

import (
	"bufio"
	"os"
	"scal/basic"
)

var wrbuffsize = 100 * basic.MB

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
