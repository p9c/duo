package check
import (
	"sync"
)
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
type Control struct {
	queue *Queue
	done  bool
}
