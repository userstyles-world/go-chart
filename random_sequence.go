package chart

import (
	"math"
	"math/rand"
	"time"
)

var (
	_ Sequence = (*RandomSeq)(nil)
)

// RandomValues returns an array of random values.
func RandomValues(count int) []float64 {
	return Seq{NewRandomSequence().WithLen(count)}.Values()
}

// RandomValuesWithMax returns an array of random values with a given average.
func RandomValuesWithMax(count int, max float64) []float64 {
	return Seq{NewRandomSequence().WithMax(max).WithLen(count)}.Values()
}

// NewRandomSequence creates a new random seq.
func NewRandomSequence() *RandomSeq {
	return &RandomSeq{
		rnd: rand.NewSource(time.Now().UnixNano()),
	}
}

func randFloat64(r rand.Source) float64 {
again:
	f := float64(r.Int63()) / (1 << 63)
	if f == 1 {
		goto again // resample; this branch is taken O(never)
	}
	return f
}

// RandomSeq is a random number seq generator.
type RandomSeq struct {
	rnd rand.Source
	max *float64
	min *float64
	len *int
}

// Len returns the number of elements that will be generated.
func (r *RandomSeq) Len() int {
	if r.len != nil {
		return *r.len
	}
	return math.MaxInt32
}

// GetValue returns the value.
func (r *RandomSeq) GetValue(_ int) float64 {
	switch {
	case r.min != nil && r.max != nil:
		var delta float64

		if *r.max > *r.min {
			delta = *r.max - *r.min
		} else {
			delta = *r.min - *r.max
		}

		return *r.min + (randFloat64(r.rnd) * delta)
	case r.max != nil:
		return randFloat64(r.rnd) * *r.max
	case r.min != nil:
		return *r.min + randFloat64(r.rnd)
	default:
		return randFloat64(r.rnd)
	}
}

// WithLen sets a maximum len
func (r *RandomSeq) WithLen(length int) *RandomSeq {
	r.len = &length
	return r
}

// Min returns the minimum value.
func (r RandomSeq) Min() *float64 {
	return r.min
}

// WithMin sets the scale and returns the Random.
func (r *RandomSeq) WithMin(min float64) *RandomSeq {
	r.min = &min
	return r
}

// Max returns the maximum value.
func (r RandomSeq) Max() *float64 {
	return r.max
}

// WithMax sets the average and returns the Random.
func (r *RandomSeq) WithMax(max float64) *RandomSeq {
	r.max = &max
	return r
}
