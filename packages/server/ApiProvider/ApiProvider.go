package ApiProvider

import (
	"server/ApiProvider/pkg/GameOperator"
)

func ApplyDataSource(source interface{}) {
	GameOperator.ApplyDataSource(source)
}
