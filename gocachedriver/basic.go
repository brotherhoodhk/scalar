package gocachedriver

type CacheDriver struct {
	hostadd, db string
	port        int
}

func NewDriver(host string, port int, db string) *CacheDriver {
	return &CacheDriver{hostadd: host, port: port, db: db}
}
