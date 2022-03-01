package legacy

import (
	"sync"
)

// Singleton - private type, that wraps library
type Singleton struct {
	IsAllocated bool
}

// variable for storing instance of Singleton
var instance *Singleton = nil

// once is ann object of type sync.Once that guarantee that once.Do
// call underline function only once.
var once sync.Once
