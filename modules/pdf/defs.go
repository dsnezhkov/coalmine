package coalpdf

import (
	"regexp"
	"sync"
)

type CPDFManager struct {
	Mutex *sync.RWMutex
	MaxObjects int64
	Honeys map[string][]string
}
type Detector struct {
	CanaryOrgTokens []regexp.Regexp
}



