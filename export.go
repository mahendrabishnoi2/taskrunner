package taskrunner

import "context"

func NewTaskManager() *TaskManager {
	return &TaskManager{
		ledger: make(map[string]*Task),
	}
}

func newMetric(cancelFunc context.CancelFunc, fn runnable) *Task {
	return &Task{
		cancelFunc: cancelFunc,
		meta:       make(map[string]interface{}),
		source:     fn,
	}
}
