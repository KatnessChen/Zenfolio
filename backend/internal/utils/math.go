package utils

import (
	"math"
	"strconv"
	"time"
)

// XIRR calculates the internal rate of return for irregular cash flows using Newton-Raphson method.
func XIRR(cashFlows []struct {
	Amount float64
	Date   time.Time
}) float64 {
	if len(cashFlows) < 2 {
		return 0
	}

	rate := 0.1
	tolerance := 1e-6
	maxIterations := 1000
	baseDate := cashFlows[0].Date

	for i := 0; i < maxIterations; i++ {
		npv := 0.0
		npvDerivative := 0.0

		for _, cf := range cashFlows {
			days := cf.Date.Sub(baseDate).Hours() / 24
			years := days / 365.25
			denominator := math.Pow(1+rate, years)
			npv += cf.Amount / denominator
			npvDerivative -= cf.Amount * years / (denominator * (1 + rate))
		}

		if math.Abs(npv) < tolerance {
			return rate
		}
		if math.Abs(npvDerivative) < tolerance {
			break
		}
		newRate := rate - npv/npvDerivative
		if newRate < -0.99 || newRate > 10.0 {
			break
		}
		rate = newRate
	}
	return 0
}

// StandardDeviation calculates the standard deviation of a slice of float64.
func StandardDeviation(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	mean := sum / float64(len(values))
	sumSquaredDiff := 0.0
	for _, value := range values {
		diff := value - mean
		sumSquaredDiff += diff * diff
	}
	variance := sumSquaredDiff / float64(len(values))
	return math.Sqrt(variance)
}

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
