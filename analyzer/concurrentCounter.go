package analyzer

import (
	"sync"
)

type concurrentCounter struct {
	count  int
	locker sync.Mutex
}

func (c *concurrentCounter) getCount() int {
	return c.count
}

func (c *concurrentCounter) addAndGetCount() int {
	c.locker.Lock()
	c.count++
	temp := c.count
	c.locker.Unlock()
	return temp
}
