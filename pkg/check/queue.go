package check

import (
	"sync"
)

// Queue is a queue of items waiting to be worked on
type Queue struct {
	mutex                            sync.RWMutex
	conditionWorker, conditionMaster bool
	queue                            []bool
	idle, total                      int
	allOK                            bool
	todo                             uint
	quit                             bool
	batchSize                        uint
}

// Control is the data tracking a collection of tasks to be done
type Control struct {
	queue *Queue
	done  bool
}
