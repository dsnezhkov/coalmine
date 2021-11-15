package word

import "sync"

type CDOCManager struct {
	Mutex *sync.RWMutex
	MaxObjects int64
	Honeys map[string][]string
}

var cdocm *CDOCManager

func CDOCManagerFactory() *CDOCManager {

	if cdocm == nil {
		cdocm = &CDOCManager{
			Mutex: &sync.RWMutex{},
			Honeys:  make(map[string][]string),
		}
	}
	return cdocm
}

