package dini

import "sync"

type CDINIManager struct {
	Mutex *sync.RWMutex
	MaxObjects int64
	Honeys map[string][]string
}

var cdinim *CDINIManager

func CDINIManagerFactory() *CDINIManager {

	if cdinim == nil {
		cdinim = &CDINIManager{
			Mutex: &sync.RWMutex{},
			Honeys:  make(map[string][]string),
		}
	}
	return cdinim
}

func (cdinim *CDINIManager) GetHoneys() map[string][]string {
	return cdinim.Honeys
}

