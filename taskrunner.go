package taskrunner

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

var (
	ErrorTaskManagerIsHalted = errors.New("task manager is shutting down")
)

type runnable func(ctx context.Context, meta Task)

type taskManager struct {
	sync.RWMutex
	count  int
	open   bool
	ledger map[string]Task
}

func (tm *taskManager) addTask(cancelFunc context.CancelFunc, taskid string, fn runnable) Task {
	tm.Lock()
	defer tm.Unlock()
	tm.count++
	tm.ledger[taskid] = newMetric(cancelFunc, fn)
	return tm.ledger[taskid]
}

func (tm *taskManager) done(taskid string) {
	tm.Lock()
	defer tm.Unlock()
	tm.count--
	delete(tm.ledger, taskid)
}

func (tm *taskManager) GO(ctx context.Context, fn runnable) (taskId string, err error) {
	if tm.open {
		taskId = uuid.NewV4().String()
		go func() {
			ctx, cancel := context.WithCancel(ctx)
			fn(ctx, tm.addTask(cancel, taskId, fn))
			tm.done(taskId)
		}()
		return
	}
	return taskId, ErrorTaskManagerIsHalted
}

func (tm *taskManager) IsOpen() bool {
	tm.RLock()
	defer tm.RUnlock()
	return tm.open
}

func (tm *taskManager) FindFromMetaKey(key string) (tasks []Task) {
	tm.RLock()
	defer tm.RUnlock()

	for _, task := range tm.ledger {
		if _, ok := task.GetMeta(key); ok {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

func (tm *taskManager) CancelTaskFromMetaKey(key string) (count int, err error) {
	if !tm.IsOpen() {
		return count, ErrorTaskManagerIsHalted
	}
	tm.Lock()
	defer tm.Unlock()

	for _, task := range tm.ledger {
		if _, ok := task.GetMeta(key); ok {
			task.Cancel()
			count++
		}
	}
	return
}

func (tm *taskManager) RestartTasksFromMetaKey(key string) (count int, err error) {
	if !tm.IsOpen() {
		return count, ErrorTaskManagerIsHalted
	}
	tm.RLock()
	defer tm.RUnlock()
	for _, task := range tm.ledger {
		if _, ok := task.GetMeta(key); ok {
			source := task.getSource()
			task.Cancel()
			tm.GO(context.TODO(), source)
			count++
		}
	}
	return
}

func (tm *taskManager) GetTasksCount() int {
	tm.RLock()
	tm.RUnlock()
	return tm.count
}
