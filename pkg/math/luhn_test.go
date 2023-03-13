package shared

import (
	"errors"
	"testing"
)

func TestCalculateLuhnSum(t *testing.T) {
	type testCase struct {
		name          string
		number        string
		parity        int
		expectedSum   int64
		expectedError error
	}

	testCases := []testCase{
		{
			name:        "Calculating Luhn sum for 4111111111111111 with parity 0",
			number:      "4111111111111111",
			parity:      0,
			expectedSum: 30,
		},
		{
			name:        "Calculating Luhn sum for 378282246310005 with parity 1",
			number:      "378282246310005",
			parity:      1,
			expectedSum: 60,
		},
		{
			name:        "Calculating Luhn sum for empty number should retutn 0",
			number:      "",
			parity:      0,
			expectedSum: 0,
		},
		{
			name:          "Calculating Luhn sum for empty number should retutn 0",
			number:        "3dfgd",
			parity:        0,
			expectedError: errors.New("invalid digit"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualSum, actualError := calculateLuhnSum(tc.number, tc.parity)
			if actualError != nil {
				if tc.expectedError == nil {
					t.Fatalf("Unexpected error: %v", actualError)
				}
				if actualError.Error() != tc.expectedError.Error() {
					t.Fatalf("Expected error '%v' but got '%v'", tc.expectedError, actualError)
				}
			} else {
				if tc.expectedError != nil {
					t.Fatalf("Expected error '%v' but got nil", tc.expectedError)
				}
				if actualSum != tc.expectedSum {
					t.Fatalf("Expected Luhn sum to be %d, but got %d", tc.expectedSum, actualSum)
				}
			}
		})
	}
}

func TestIsLuhnValid(t *testing.T) {
	type testCase struct {
		name           string
		number         string
		expectedResult bool
	}

	testCases := []testCase{
		{
			name:           "Number 4111111111111111 is valid",
			number:         "4111111111111111",
			expectedResult: true,
		},
		{
			name:           "Number 6011111111111117 is valid",
			number:         "6011111111111117",
			expectedResult: true,
		},
		{
			name:           "Number 1234567890123456 is invalid",
			number:         "1234567890123456",
			expectedResult: false,
		},
		{
			name:           "Number 1234 is invalid",
			number:         "1234",
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualResult := IsLuhnValid(tc.number)
			if actualResult != tc.expectedResult {
				t.Fatalf("Expected Luhn validation to be %v, but got %v", tc.expectedResult, actualResult)
			}
		})
	}
}
