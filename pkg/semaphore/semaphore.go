package semaphore

import (
	"sync"
)

// Semaphore is
type Semaphore struct {
	condition bool
	mutex     sync.RWMutex
	value     int
}

// Grant is
type Grant struct {
	sem  *Semaphore
	have bool
}
