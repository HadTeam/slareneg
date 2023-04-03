package GameOperator

import (
	"server/Untils/pkg/DataSource"
)

var data DataSource.TempDataSource

func ApplyDataSource(source interface{}) {
	data = source.(DataSource.TempDataSource)
}
