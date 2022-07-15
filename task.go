package taskrunner

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrorKeyAlreadyExistInMeta = errors.New("key already exist in meta")
	ErrorKeyDidNotExistInMeta  = errors.New("key is not present in meta")
)

type task struct {
	sync.RWMutex
	cancelFunc context.CancelFunc
	source     runnable
	meta       map[string]interface{}
}

func (t *task) AddMeta(key string, val interface{}) error {
	t.Lock()
	defer t.Unlock()
	if _, ok := t.meta[key]; ok {
		return ErrorKeyAlreadyExistInMeta
	}
	t.meta[key] = val

	return nil
}

func (t *task) RemoveMeta(key string) error {
	t.Lock()
	defer t.Unlock()
	if _, ok := t.meta[key]; !ok {
		return ErrorKeyDidNotExistInMeta
	}
	delete(t.meta, key)
	return nil
}

func (t *task) Cancel() {
	t.Lock()
	defer t.Unlock()
	t.cancelFunc()
}

func (t *task) GetMeta(key string) (val interface{}, ok bool) {
	t.RLock()
	defer t.RUnlock()
	val, ok = t.meta[key]
	return
}

func (t *task) getSource() runnable {
	t.RLock()
	defer t.RUnlock()
	return t.source
}
