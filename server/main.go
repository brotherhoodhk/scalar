package server

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"scal/basic"
	"scal/gocachedriver"
	"strings"
	"syscall"

	"github.com/oswaldoooo/octools/toolsbox"
)

// accept msg
type Message struct {
	Zone  string `json:"zone"`
	Key   string `json:"key"`
	Value []byte `json:"value"`
	Act   int    `json:"act"`
}

// status replay
type ReplayStatus struct {
	Content    []byte `json:"content"`
	StatusCode int    `json:"code"`
}

var errorlog = toolsbox.LogInit("error", os.Getenv("SCALAR_HOME"+"/logs/error.log"))

func Start() {
	listener, err := net.Listen("unix", os.Getenv("SCALAR_HOME")+"/tmp/scalar.socket")
	if err == nil {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
		go func() {
			<-c
			os.Remove(os.Getenv("SCALAR_HOME") + "/tmp/scalar.socket")
			os.Exit(1)
		}()
		for {
			ac, err := listener.Accept()
			if err == nil {
				go Process(ac)
			}
		}
	} else {
		os.Remove(os.Getenv("SCALAR_HOME") + "/tmp/scalar.socket")
	}
}
func Process(con net.Conn) {
	var buff = make([]byte, 10*basic.MB)
	var lang int
	var err error
	msg := new(Message)
	rpy := new(ReplayStatus)
	for {
		var senddata = make([]byte, 10*basic.MB)
		lang, err = con.Read(buff)
		if err == nil {
			fmt.Println("accept mess\n", string(buff[:lang]))
			err = json.Unmarshal(buff[:lang], msg)
			if err == nil {
				if msg.Act != 10 {
					if zoneid, ok := gocachedriver.CheckZone(msg.Zone); ok {
						switch msg.Act {
						case 1:
							err = gocachedriver.SetKey(msg.Key, string(msg.Value), zoneid)
						case 2:
							res, err := gocachedriver.GetKey(msg.Key, zoneid)
							if err == nil {
								senddata = []byte(res)
							}
						case 3:
							err = gocachedriver.Delete(msg.Key, string(msg.Value))
						default:
							err = fmt.Errorf("unknown command :-(")
						}
					}
				} else {
					if len(msg.Zone) > 0 && !strings.ContainsRune(msg.Zone, ' ') {
						err = gocachedriver.CreateZone(msg.Zone)
					} else {
						err = fmt.Errorf("zone format is wrong")
					}
				}
			} else {
				rpy.StatusCode = 500
				goto sendtocli
			}
		}
		if err == nil {
			rpy.StatusCode = 200
		} else {
			rpy.StatusCode = 400
			senddata = []byte(err.Error())
		}
		rpy.Content = senddata
	sendtocli:
		resbytes, _ := json.Marshal(rpy)
		_, err = con.Write(resbytes)
		if err != nil {
			errorlog.Println(err)
		}
	}
}
