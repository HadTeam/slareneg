package JudgePool

import (
	"server/ApiProvider/pkg/DataOperator"
	"server/JudgePool/internal/InstructionExecutor"
	"server/JudgePool/internal/Judge"
)

type Pool struct {
}

var data DataOperator.DataSource

func ApplyDataSource(source DataOperator.DataSource) {
	data = source
	Judge.ApplyDataSource(source)
	InstructionExecutor.ApplyDataSource(source)
}

func CreatePool() {

}
