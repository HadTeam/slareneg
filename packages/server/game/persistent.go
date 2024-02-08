package game

import (
	"server/game/map"
)

type PersistentDataSource interface {
	GetOriginalMap(mapId uint32) *_map.Map
}
