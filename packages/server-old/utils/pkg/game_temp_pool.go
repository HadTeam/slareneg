package pkg

import (
	"sync"
)

// similar with sync.Pool, but only gc when the game end

var pool = &sync.Map{}

func TempPoolCreate(id interface{}) {
	if _, ok := pool.Load(id); !ok {
		pool.Store(id, &sync.Map{})
	}
}

func TempPoolDelete(id interface{}) {
	pool.Delete(id)
}

func TempPoolGet(id interface{}, key string) (any, bool) {
	v, _ := pool.Load(id)
	if v == nil {
		return nil, false
	}
	return v.(*sync.Map).Load(key)
}

func TempPoolPut(id interface{}, key string, value any) {
	v, _ := pool.Load(id)
	if v != nil {
		v.(*sync.Map).Store(key, value)
	}
}
