package gocachedriver

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	driver_tools "github.com/oswaldoooo/gocache-driver/basics"
)

type Cache_Driver interface {
	SetKey(key, val, zoneid string) error
	GetKey(key, zoneid string) (string, error) //error only for network
	GetAllKeys(zoneid string) (map[string][]byte, error)
	DelKey(zoneid string, keys ...string) error //error only for network

}

// this distribute file system extension
// check the node is online
type node struct {
	ipadd   string
	port    int
	alias   string
	authkey string
	db      string
	// con     net.Conn
	con *driver_tools.CacheDB
}

func (s *node) checknode() bool {
	con, err := net.Dial("tcp", s.ipadd+":"+strconv.Itoa(s.port))
	if err == nil {
		con.Close()
		return true
	} else {
		return false
	}
}
func (s *node) connection() (err error) {
	// con, err := net.Dial("tcp", s.ipadd+":"+strconv.Itoa(s.port))
	// if err == nil {
	// 	s.con = con
	// }
	// s.con = driver_tools.New(s.ipadd, s.port, s.authkey, s.db)
	return s.con.Connect()
}
func (s *node) sendcommand() {

}

// insert data into node server
func (s *node) insertdata(key, value, zoneid string) (err error) {
	// err = dbcon.SetKey(key, value)

	return
}
func (s *node) deletedata(key, zoneid string) (err error) {
	err = dbcon.DeleteKeys(key)
	return
}

// test function,to compelete
func (s *node) createdb() (err error) {
	return
}

// test function,to compelete
func (s *node) deletedb() (err error) {
	return
}

func (s *node) modifydata(key, value, zoneid string) (err error) {
	rec, err := dbcon.GetKeys(key)
	if err != nil || len(rec) == 0 {
		//key is not existed
		err = errors.New(key + " is not existed,cant modify it")
	} else {
		//key is existed
		err = dbcon.SetKey(key, value)
	}
	return
}
func (s *node) SetKey(key string, val string, zoneid string) error {
	return s.con.SetKey(zoneid+key, val)
}

func (s *node) GetKey(key string, zoneid string) (res string, err error) {
	if len(zoneid) != 6 {
		err = fmt.Errorf("zoneid's format is not correct")
		return
	}
	key = zoneid + key
	resarr, err := s.con.GetKeys(key)
	if err == nil && len(resarr) > 0 {
		res = resarr[0]
	}
	return
}

func (s *node) GetAllKeys(zoneid string) (res map[string][]byte, err error) {
	if len(zoneid) != 6 {
		err = fmt.Errorf("zoneid's format is not correct")
		return
	}
	res, err = s.con.GetKeysContain(zoneid)
	return
}

func (s *node) DelKey(zoneid string, keys ...string) (err error) {
	if len(zoneid) != 6 {
		err = fmt.Errorf("zoneid's format is not correct")
		return
	}
	if len(keys) == 1 && keys[0] == "*" {
		var resmap map[string][]byte
		resmap, err = s.con.GetKeysContain(zoneid)
		if err == nil {
			keysarr := []string{}
			for k := range resmap {
				keysarr = append(keysarr, k)
			}
			if len(keysarr) > 0 {
				err = s.con.DeleteKeys(keysarr...)
			}
		}
	} else if len(keys) > 0 {
		err = s.con.DeleteKeys(keys...)
	}
	return
}

// cluster interface
/*
notice:
1. cluster value need set version;version is top two byte
*/
const (
	UNUSED  = 0
	OK      = 1
	UNKNOWN = 3
)

type Cluster interface {
	Cache_Driver
	Join(...*ClusterOption) error
	Kick(...*ClusterOption) //kick node from cluster
}
type ClusterOption struct {
	Host     string `json:"host" xml:"host" yaml:"host"`
	Port     int    `json:"port" xml:"port" yaml:"port"`
	Password string `json:"password" xml:"password" yaml:"password"`
	Database string `json:"database" xml:"database" yaml:"database"`
	status   uint8
	db       *driver_tools.CacheDB
}
type ClusterDB struct {
	db_pool     []*ClusterOption
	version_map map[string]map[string]uint32
}

func NewClusterDB(options ...*ClusterOption) (*ClusterDB, error) {
	if len(options) == 0 {
		return nil, fmt.Errorf("not input any cluster option")
	}
	ans := &ClusterDB{db_pool: options, version_map: make(map[string]map[string]uint32)}
	var err, errtwo error
	errlist := []error{}
	for _, ele := range ans.db_pool {
		ele.db = driver_tools.New(ele.Host, ele.Port, ele.Password, ele.Database)
		errtwo = ele.db.Connect()
		if errtwo != nil {
			errlist = append(errlist, errtwo)
			ele.status = UNUSED
			fmt.Printf("gocache %s:%d init failed\n", ele.Host, ele.Port)
		} else {
			fmt.Printf("gocache %s:%d init success\n", ele.Host, ele.Port)
			ele.status = OK
		}
	}
	if len(errlist) > 0 {
		err = errors.Join(errlist...)
	}
	return ans, err
}

type version_info struct {
	version  uint32
	nodeinfo *ClusterOption
}

func (s *ClusterDB) SetKey(key string, val string, zoneid string) error {
	var err, errtwo error
	errlist := []error{}
	for _, node := range s.db_pool { //set value to cluster
		if node.status == OK {
			_, errtwo = node.db.CompareAndSetKey(zoneid+key, val)
			if errtwo != nil {
				errlist = append(errlist, errtwo)
			}
		}
	}
	if len(errlist) > 0 {
		err = errors.Join(errlist...)
	}
	return err
}

func (s *ClusterDB) GetKey(key string, zoneid string) (string, error) {
	var (
		ansarr          []string
		err, errtwo     error
		maxversion_data struct {
			version uint32
			data    string
		}
		version uint32
	)

	errlist := []error{}
	for _, node := range s.db_pool { //get data from all available nodes and return the data which version is max
		if node.status == OK {
			ansarr, errtwo = node.db.GetKeys(zoneid + key)
			if errtwo != nil {
				errlist = append(errlist, errtwo)
			} else if ansarr != nil {
				version = uint32(ansarr[0][0])*256 + uint32(ansarr[0][1])
				if version > maxversion_data.version {
					maxversion_data.data = ansarr[0][2:]
				}
			}
		}
	}
	if len(errlist) > 0 {
		err = errors.Join(errlist...)
	}
	return maxversion_data.data, err
}

func (s *ClusterDB) GetAllKeys(zoneid string) (map[string][]byte, error) {
	var (
		ans         map[string][]byte = make(map[string][]byte)
		ansarr      map[string][]byte = make(map[string][]byte)
		err, errtwo error
		ok          bool
	)
	errlist := []error{}
	for _, node := range s.db_pool {
		if node.status == OK {
			ansarr, errtwo = node.db.GetKeysContain(zoneid)
			if errtwo != nil {
				errlist = append(errlist, errtwo)
			} else if ansarr != nil && len(ansarr) > 0 {
				for key, val := range ansarr {
					if _, ok = ans[key]; !ok {
						ans[key] = val
					}
				}
			}
		}
	}
	if len(errlist) > 0 {
		err = errors.Join(errlist...)
	}
	return ans, err
}

func (s *ClusterDB) DelKey(zoneid string, keys ...string) error {
	panic("not implemented") // TODO: Implement
}
func (s *ClusterDB) Join(options ...*ClusterOption) error {
	if len(options) == 0 {
		return fmt.Errorf("not input any cluster option")
	}
	var err, errtwo error
	errlist := []error{}
	for _, ele := range options {
		ele.db = driver_tools.New(ele.Host, ele.Port, ele.Password, ele.Database)
		errtwo = ele.db.Connect()
		if errtwo != nil {
			errlist = append(errlist, errtwo)
			ele.status = UNUSED
		} else {
			ele.status = OK
		}
	}
	if len(errlist) > 0 {
		err = errors.Join(errlist...)
	}
	s.db_pool = append(s.db_pool, options...)
	return err
}

func (s *ClusterDB) Kick(options ...*ClusterOption) {

}
