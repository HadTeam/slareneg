package InstructionExecutor

import (
	"fmt"
	"server/Untils/pkg/DataSource"
	"server/Untils/pkg/GameType"
	"server/Untils/pkg/InstructionType"
	"server/Untils/pkg/MapType"
)

var data DataSource.TempDataSource

func ApplyDataSource(source interface{}) {
	data = source.(DataSource.TempDataSource)
}

func ExecuteInstruction(id GameType.GameId, instruction InstructionType.Instruction) bool {
	var ret bool
	var m *MapType.Map
	switch instruction.(type) {
	case InstructionType.Move:
		{
			m = data.GetCurrentMap(id)
			ret = m.Move(instruction.(InstructionType.Move))

		}
	}
	if !ret {
		fmt.Printf("[Warn] Execute instruction failed: %#v \n", instruction)
	}
	return ret
}
