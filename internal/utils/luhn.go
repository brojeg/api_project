package utils

import (
	"errors"
	"io"
)

const (
	asciiZero = 48
	asciiTen  = 57
)

func IsLuhnValid(number string) bool {
	p := len(number) % 2
	sum, err := calculateLuhnSum(number, p)
	if err != nil {
		return false
	}

	// If the total modulo 10 is not equal to 0, then the number is invalid.
	if sum%10 != 0 {
		return false
	}

	return true
}
func calculateLuhnSum(number string, parity int) (int64, error) {
	var sum int64
	for i, d := range number {
		if d < asciiZero || d > asciiTen {
			return 0, errors.New("invalid digit")
		}

		d = d - asciiZero
		// Double the value of every second digit.
		if i%2 == parity {
			d *= 2
			// If the result of this doubling operation is greater than 9.
			if d > 9 {
				// The same final result can be found by subtracting 9 from that result.
				d -= 9
			}
		}

		// Take the sum of all the digits.
		sum += int64(d)
	}

	return sum, nil
}

// func StringToIntSlice(n int64) []int64 {
// 	var ret []int64
// 	// i, err := strconv.ParseInt(n, 10, 64)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	for n != 0 {
// 		ret = append(ret, n%10)
// 		n /= 10
// 	}

// 	reverseInt(ret)

// 	return ret
// }

func GetRawOrderNumber(body io.Reader) (string, error) {
	b, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
