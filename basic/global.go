package basic

var WordCountFunc = make(map[string]func(filepath string) map[string]int)
var ValueCalFunc = make(map[string]*VC)

type VC struct {
	RankFunc func(map[string][]byte) map[string][]byte
	GetRank  func(map[string][]byte) []struct {
		Key   string
		Value []byte
	}
}
