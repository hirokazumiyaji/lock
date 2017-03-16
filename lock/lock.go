package lock

import "sync"

var (
	object map[string]bool
	m      *sync.Mutex
)

func Initialize(cap int) {
	object = make(map[string]bool, cap)
	m = new(sync.Mutex)
}

func lock(key string) bool {
	m.Lock()
	defer m.Unlock()

	if _, ok := object[key]; ok {
		return false
	}
	object[key] = true
	return true
}

func unlock(key string) bool {
	m.Lock()
	defer m.Unlock()

	delete(object, key)

	return true
}
