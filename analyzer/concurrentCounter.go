package analyzer

import "sync/atomic"

// concurrentCounter A struct that has a thread safe counter
type concurrentCounter struct {
	count int64
}

// Returns the current count
func (c *concurrentCounter) getCount() int {
	return int(c.count)
}

// increments the counter and returns the new value
func (c *concurrentCounter) addAndGetCount() int {
	return int(atomic.AddInt64(&c.count, 1))
}
