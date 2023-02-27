package utils

import (
	"io"
	"io/ioutil"
	"strconv"
)

func IsLuhnValid(value []int64) bool {

	sum := computeCheckSum(value)
	return sum%10 == 0
}
func computeCheckSum(data []int64) int64 {
	var sum int64
	double := false
	for _, n := range data {
		if double {
			n = (n * 2)
			if n > 9 {
				n = (n - 9)
			}
		}
		sum += n
		double = !double
	}
	return sum
}

func reverseInt(s []int64) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func StringToIntSlice(n int64) []int64 {
	var ret []int64
	// i, err := strconv.ParseInt(n, 10, 64)
	// if err != nil {
	// 	panic(err)
	// }
	for n != 0 {
		ret = append(ret, n%10)
		n /= 10
	}

	reverseInt(ret)

	return ret
}

func GetRawOrderNumber(body io.Reader) int64 {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		panic(err)
	}

	rawOrderNumber, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		panic(err)
	}
	return rawOrderNumber
}
