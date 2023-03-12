package basic

import "os"

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

var ROOTPATH = os.Getenv("SCALAR_HOME")

const (
	VERSION = "v1.0"
)

type Container struct {
	Origin       string
	originnumber []int
	originfloat  []float64
}
