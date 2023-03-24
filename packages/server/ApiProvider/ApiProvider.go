package ApiProvider

import (
	"server/ApiProvider/pkg/DataOperator"
	"server/ApiProvider/pkg/GameOperator"
)

func ApplyDataSource(source DataOperator.DataSource) {
	GameOperator.ApplyDataSource(source)
}
