package pdf

import "sync"

type CPDFManager struct {
	Mutex *sync.RWMutex
	MaxObjects int64
	Honeys map[string][]string
}

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

func (cpdfm *CPDFManager) GetHoneys() map[string][]string {
	return cpdfm.Honeys
}
