package InstructionExecutor

import (
	"fmt"
	"server/ApiProvider/pkg/DataOperator"
	"server/Untils/pkg/GameType"
	"server/Untils/pkg/InstructionType"
	"server/Untils/pkg/MapOperator"
	"server/Untils/pkg/MapType"
)

var data DataOperator.DataSource

func ApplyDataSource(source DataOperator.DataSource) {
	data = source
}

func ExecuteInstruction(gameId GameType.GameId, instruction InstructionType.Instruction) bool {
	var ret bool
	var m *MapType.Map
	switch instruction.(type) {
	case InstructionType.MoveInstruction:
		{
			m = data.GetCurrentGame(gameId).Map
			ret = MapOperator.Move(m, instruction.(InstructionType.MoveInstruction))

		}
	}
	if !ret {
		fmt.Printf("[Warn] Execute instruction failed: %#v \n", instruction)
	}
	return ret
}
