package chart

import (
	"strconv"
	"sync"
)

var (
	itoaCache = make(map[int]string)
	itoaMutex = &sync.Mutex{}

	ftoaCache = make(map[float64]string)
	ftoaMutex = &sync.Mutex{}
)

// Implement a caching itoa function.
func itoa(i int) string {
	itoaMutex.Lock()
	defer itoaMutex.Unlock()
	if s, ok := itoaCache[i]; ok {
		return s
	}
	s := strconv.FormatInt(int64(i), 10)
	itoaCache[i] = s
	return s
}

// Implement a caching ftoa function.
func ftoa(f float64, prec int) string {
	ftoaMutex.Lock()
	defer ftoaMutex.Unlock()
	if s, ok := ftoaCache[f]; ok {
		return s
	}
	s := strconv.FormatFloat(f, 'f', prec, 64)
	ftoaCache[f] = s
	return s
}
