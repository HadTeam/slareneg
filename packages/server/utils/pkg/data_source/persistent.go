package data_source

import (
	"server/game_logic/map"
)

type PersistentDataSource interface {
	GetOriginalMap(mapId uint32) *_map.Map
}
