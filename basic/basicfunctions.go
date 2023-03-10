package basic

import "fmt"

func (s *Container) Scan() {
	s.originnumber = ScanNumber(s.Origin)
}
func (s *Container) ScanFloat() {
	s.originfloat = ScanFloat(s.Origin)
	fmt.Println(s.originfloat)
}
func (s *Container) Sum() int {
	res := 0
	for _, v := range s.originnumber {
		res += v
	}
	return res
}

// return max,min
func (s *Container) MaxMin() (int, int) {
	if len(s.originnumber) == 0 {
		return 0, 0
	}
	max := s.originnumber[0]
	min := s.originnumber[0]
	for _, v := range s.originnumber {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	return max, min
}
