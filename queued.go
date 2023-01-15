package queued

import "sync"

var queue = make(map[string][]*sync.WaitGroup)
var queueLocker sync.Mutex

func Queued[T any](id string, execute func() (T, error)) (T, error) {
	queueLocker.Lock()
	if groups, ok := queue[id]; ok {
		var wg sync.WaitGroup
		wg.Add(1)
		queue[id] = append(groups, &wg)
		queueLocker.Unlock()
		wg.Wait()
	} else {
		queue[id] = make([]*sync.WaitGroup, 0)
		queueLocker.Unlock()
	}

	defer func() {
		queueLocker.Lock()
		if groups, ok := queue[id]; ok {
			for _, wg := range groups {
				wg.Done()
			}
			delete(queue, id)
		}
		queueLocker.Unlock()
	}()

	return execute()
}
