package dsp

import (
	"math"
	"testing"
)

func TestBiquadFilter(t *testing.T) {
	// A simple moving average filter: y[n] = 0.5*x[n] + 0.5*x[n-1]
	// b0 = 0.5, b1 = 0.5, b2 = 0, a1 = 0, a2 = 0
	filter := NewBiquadFilter(0.5, 0.5, 0.0, 0.0, 0.0)
	defer filter.Close()

	input := []float32{1.0, 1.0, 1.0, 1.0}
	// Expected outputs:
	// y[0] = 0.5*1.0 + 0.5*0.0 = 0.5
	// y[1] = 0.5*1.0 + 0.5*1.0 = 1.0
	// y[2] = 0.5*1.0 + 0.5*1.0 = 1.0
	// y[3] = 0.5*1.0 + 0.5*1.0 = 1.0
	expected := []float32{0.5, 1.0, 1.0, 1.0}

	filter.Process(input)

	for i := range input {
		if math.Abs(float64(input[i]-expected[i])) > 1e-6 {
			t.Errorf("Mismatch at index %d: got %f, expected %f", i, input[i], expected[i])
		}
	}
}
