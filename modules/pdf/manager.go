package coalpdf

import "sync"

var cpdfm *CPDFManager

func CPDFManagerFactory() *CPDFManager {

	if cpdfm == nil {
		cpdfm = &CPDFManager{
			Mutex: &sync.RWMutex{},
			MaxObjects: 50,
			Honeys:  make(map[string][]string),
		}
	}
	return cpdfm
}

