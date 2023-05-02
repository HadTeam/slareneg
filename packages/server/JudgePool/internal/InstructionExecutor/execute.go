package InstructionExecutor

import (
	"log"
	"server/Utils/pkg/DataSource"
	"server/Utils/pkg/GameType"
	"server/Utils/pkg/InstructionType"
	"server/Utils/pkg/MapType"
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
		log.Printf("[Warn] Execute instruction failed: %#v \n", instruction)
	}
	return ret
}
