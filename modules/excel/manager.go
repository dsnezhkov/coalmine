package excel

import "sync"

type CXLSManager struct {
	Mutex *sync.RWMutex
	MaxObjects int64
	Honeys map[string][]string
}

var cxlsm *CXLSManager

func CXLSManagerFactory() *CXLSManager {

	if cxlsm == nil {
		cxlsm = &CXLSManager{
			Mutex: &sync.RWMutex{},
			Honeys:  make(map[string][]string),
		}
	}
	return cxlsm
}

func (cxlsm *CXLSManager) GetHoneys() map[string][]string {
	return cxlsm.Honeys
}

