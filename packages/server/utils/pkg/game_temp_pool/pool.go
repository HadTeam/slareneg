package game_temp_pool

import (
	"server/game_logic"
	"sync"
)

// similar with sync.Pool, but only gc when the game end

var pool = &sync.Map{}

func Create(id game_logic.Id) {
	if _, ok := pool.Load(id); !ok {
		pool.Store(id, &sync.Map{})
	}
}

func Delete(id game_logic.Id) {
	pool.Delete(id)
}

func Get(id game_logic.Id, key string) (any, bool) {
	v, _ := pool.Load(id)
	return v.(*sync.Map).Load(key)
}

func Put(id game_logic.Id, key string, value any) {
	v, _ := pool.Load(id)
	v.(*sync.Map).Store(key, value)
}
