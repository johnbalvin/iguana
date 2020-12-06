package iguana

import (
	"sync"
)

var lockHTML sync.RWMutex
var lockStatic sync.RWMutex
