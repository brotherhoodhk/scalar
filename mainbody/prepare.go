package mainbody

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"plugin"
	"scal/basic"
	"strings"

	"github.com/oswaldoooo/octools/toolsbox"
)

type basictype interface {
	int | float64 | string | []string | []int | []float64
}
type confinfo struct {
	XMLName xml.Name `xml:"scalar"`
	Plugins []struct {
		XMLName  xml.Name `xml:"plugin_info"`
		Class    string   `xml:"classname"`
		FileName string   `xml:"filename"`
	} `xml:"plugins"`
	PathInfo pathinfo `xml:"paths"`
}
type pathinfo struct {
	XMLName     xml.Name `xml:"paths"`
	Common_Path string   `xml:"common_path"`
}

func init() {
	//识别plugin信息
	siteconfpath := basic.ROOTPATH + "/conf/site.xml"
	conf := new(confinfo)
	fmt.Println("start init plugins info")
	content, err := ioutil.ReadFile(siteconfpath)
	if err == nil {
		err = xml.Unmarshal(content, conf)
		if err == nil {
			for _, plugin_info := range conf.Plugins {
				pluginer, err := loadplugins(plugin_info.FileName)
				if err == nil {
					switch strings.ToLower(plugin_info.Class) {
					//plugin分类，查看是否符合分类
					case "wordcount":
						resfun, err := loadwordcountfunc(pluginer)
						if err == nil {
							realname := ""
							if strings.Contains(plugin_info.FileName, ".so") {
								realname = strings.Replace(plugin_info.FileName, ".so", "", 1)
							} else {
								realname = plugin_info.FileName
							}
							basic.WordCountFunc[realname] = resfun
						}
					case "value calculation":
						resfunone, resfuntwo, err := loadvalcal(pluginer)
						realname := ""
						if err == nil {
							if strings.Contains(plugin_info.FileName, ".so") {
								realname = strings.Replace(plugin_info.FileName, ".so", "", 1)
							} else {
								realname = plugin_info.FileName
							}
							basic.ValueCalFunc[realname] = &basic.VC{RankFunc: resfuntwo, GetRank: resfunone}
						}
					default:
						err = fmt.Errorf("dont support class %v", plugin_info.Class)
					}
				}
				if err != nil {
					fmt.Println(err)
				}
			}
			if _, err = os.Stat(conf.PathInfo.Common_Path); len(conf.PathInfo.Common_Path) > 0 && err == nil {
				//通用插件目录存在

			}
		}
	}
}

// 加载plugin
func loadplugins(filename string) (newplugin *plugin.Plugin, err error) {
	if strings.Contains(filename, ".so") {
		newplugin, err = plugin.Open(basic.ROOTPATH + "/plugins/" + filename)
	} else {
		newplugin, err = plugin.Open(basic.ROOTPATH + "/plugins/" + filename + ".so")
	}
	return
}

// 加载wordcount类插件
func loadwordcountfunc(pluginer *plugin.Plugin) (resfunc func(filename string) map[string]int, err error) {
	sym, err := pluginer.Lookup("WordCount")
	if err == nil {
		resfunc = sym.(func(filename string) map[string]int)
	}
	return
}

// 加载缓冲池权值计算插件
func loadvalcal(pluginer *plugin.Plugin) (resfuncone func(origin_data map[string][]byte) []struct {
	Key   string
	Value []byte
}, resfunctwo func(origin_data map[string][]byte) map[string][]byte, err error) {
	sym, err := pluginer.Lookup("ValueCal")
	if err == nil {
		resfuncone = sym.(func(origin_data map[string][]byte) []struct {
			Key   string
			Value []byte
		})
	} else {
		return
	}
	sym, err = pluginer.Lookup("GetFinalLeavet")
	if err == nil {
		resfunctwo = sym.(func(origin_data map[string][]byte) map[string][]byte)
	}
	return
}

// 读取common_plugins directory
func readcommonpluginsdir(dirpath string) (err error) {
	fearr, err := ioutil.ReadDir(dirpath)
	if err == nil {
		pluginer := new(plugin.Plugin)
		for _, ve := range fearr {
			if strings.Contains(ve.Name(), ".so") {
				pluginer, err = toolsbox.ScanPluginByName(ve.Name(), dirpath)
				if err == nil {
					sym, err := pluginer.Lookup("classname")
					if err == nil {
						switch *sym.(*string) {
						case "fuzzy match":
							sym, err = pluginer.Lookup("FuzzyMatch")
							if err == nil {
								//将插件添加到模糊查询池中
								basic.FuzzyMatchFunc[strings.Replace(ve.Name(), ".so", "", 1)] = sym.(func(target, tocompare string) bool)
							}
						}
					}
				}
			}
		}
	}
	return
}
