package test

import (
	"math"
	"testing"
	"time"

	"github.com/transaction-tracker/backend/internal/utils"
)

func TestStandardDeviation(t *testing.T) {
	cases := []struct {
		values   []float64
		expected float64
	}{
		{[]float64{}, 0},
		{[]float64{1, 1, 1, 1}, 0},
		{[]float64{1, 2, 3, 4, 5}, math.Sqrt(2)},
		{[]float64{2, 4, 4, 4, 5, 5, 7, 9}, 2},
	}
	for _, c := range cases {
		result := utils.StandardDeviation(c.values)
		if math.Abs(result-c.expected) > 1e-6 {
			t.Errorf("StandardDeviation(%v) = %v, want %v", c.values, result, c.expected)
		}
	}
}

func TestXIRR(t *testing.T) {
	// Example: invest -1000, after 1 year get 1100, IRR should be about 10%
	cashFlows := []struct {
		Amount float64
		Date   time.Time
	}{
		{Amount: -1000, Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Amount: 1100, Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	r := utils.XIRR(cashFlows)
	if math.Abs(r-0.1) > 1e-3 {
		t.Errorf("XIRR = %v, want ~0.1", r)
	}

	// Edge: only one cash flow
	cashFlows = []struct {
		Amount float64
		Date   time.Time
	}{
		{Amount: -1000, Date: time.Now()},
	}
	if utils.XIRR(cashFlows) != 0 {
		t.Error("XIRR with single cash flow should be 0")
	}

	// Edge: no cash flows
	if utils.XIRR([]struct {
		Amount float64
		Date   time.Time
	}{}) != 0 {
		t.Error("XIRR with empty cash flows should be 0")
	}
}
