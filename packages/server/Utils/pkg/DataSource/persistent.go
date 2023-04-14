package DataSource

import "server/Utils/pkg/MapType"

type PersistentDataSource interface {
	GetOriginalMap(mapId uint32) *MapType.Map
}
