package dataSource

import "server/utils/pkg/map"

type PersistentDataSource interface {
	GetOriginalMap(mapId uint32) *_map.Map
}
