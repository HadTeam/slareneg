package DataOperator

import (
	"fmt"
	"server/MapHandler/pkg/MapType"
)

func PutMap(newMap MapType.Map) bool {
	// TODO: Update the map in the db
	// TODO: Start to receiving new instruction
	fmt.Print(newMap)
	return true
}
