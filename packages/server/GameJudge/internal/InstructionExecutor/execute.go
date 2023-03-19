package InstructionExecutor

import (
	"fmt"
	"server/ApiProvider/pkg/InstructionType"
	"server/GameJudge/internal/DataOperator"
	"server/GameJudge/internal/GameType"
	"server/MapHandler/pkg/MapOperator"
	"server/MapHandler/pkg/MapType"
)

func ExecuteAllInstruction(gameId GameType.GameId) bool {
	instructionList := DataOperator.GetInstruction(gameId)
	ret := true
	for _, instruction := range instructionList {
		if !ExecuteInstruction(gameId, instruction) {
			ret = false
		}
	}
	return ret
}

func ExecuteInstruction(gameId GameType.GameId, instruction InstructionType.Instruction) bool {
	var ret bool
	var m *MapType.Map
	switch instruction.(type) {
	case InstructionType.MoveInstruction:
		{
			m = DataOperator.GetCurrentMap(gameId)
			ret = MapOperator.Move(m, instruction.(InstructionType.MoveInstruction))

		}
	}
	if !ret {
		fmt.Printf("[Warn] Execute instruction failed: %#v \n", instruction)
	}
	return ret
}
