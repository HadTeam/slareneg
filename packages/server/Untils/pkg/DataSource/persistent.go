package DataSource

import "server/Untils/pkg/MapType"

type PersistentDataSource interface {
	GetOriginalMap(mapId uint32) *MapType.Map
}
