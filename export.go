package taskrunner

import "context"

func NewTaskManager() TaskManager {
	return &taskManager{
		ledger: make(map[string]Task),
		open:   true,
	}
}

type TaskManager interface {
	GO(ctx context.Context, fn runnable) (taskId string, err error)
	RestartTasksFromMetaKey(key string) (count int, err error)
	CancelTaskFromMetaKey(key string) (count int, err error)
	FindFromMetaKey(key string) (tasks []Task)
	IsOpen() bool
	GetTasksCount() int
}

func newMetric(cancelFunc context.CancelFunc, fn runnable) Task {
	return &task{
		cancelFunc: cancelFunc,
		meta:       make(map[string]interface{}),
		source:     fn,
	}
}

type Task interface {
	RemoveMeta(key string) error
	AddMeta(key string, val interface{}) error
	GetMeta(key string) (val interface{}, ok bool)
	getSource() runnable
	Cancel()
}
