package utils

import (
	"math"
	"strconv"
)

// Abs returns the absolute value of x
func Abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// ParseUint parses a string to uint with a default value
func ParseUint(s string, defaultValue uint) (uint, error) {
	if s == "" {
		return defaultValue, nil
	}

	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(val), nil
}

// RoundTo4 rounds a float64 to 4 decimal places
func RoundTo4(val float64) float64 {
	return math.Round(val*1e4) / 1e4
}
