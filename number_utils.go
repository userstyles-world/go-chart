package chart

import (
	"strconv"
	"sync"
)

var (
	itoaCache = make(map[int]string)
	itoaMutex = &sync.Mutex{}

	ftoa1Cache = make(map[float64]string)
	ftoa1Mutex = &sync.Mutex{}

	ftoa2Cache = make(map[float64]string)
	ftoa2Mutex = &sync.Mutex{}
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
func ftoa1(f float64) string {
	ftoa1Mutex.Lock()
	defer ftoa1Mutex.Unlock()
	if s, ok := ftoa1Cache[f]; ok {
		return s
	}
	s := strconv.FormatFloat(f, 'f', 1, 64)
	ftoa1Cache[f] = s
	return s
}

func ftoa2(f float64) string {
	ftoa2Mutex.Lock()
	defer ftoa2Mutex.Unlock()
	if s, ok := ftoa2Cache[f]; ok {
		return s
	}
	s := strconv.FormatFloat(f, 'f', 2, 64)
	ftoa2Cache[f] = s
	return s
}
