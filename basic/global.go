package basic

var FuzzyMatchFunc = make(map[string]func(target, tocompare string) bool)
var WordCountFunc = make(map[string]func(filepath string) map[string]int)
var ValueCalFunc = make(map[string]*VC)
var Common_Path = ""
var Extension_Func = make(map[int]func(argsone, argstwo string, argsthree []byte) (content []byte, typeinfo string, err error)) //扩展功能

type VC struct {
	RankFunc func(map[string][]byte) map[string][]byte
	GetRank  func(map[string][]byte) []struct {
		Key   string
		Value []byte
	}
}
