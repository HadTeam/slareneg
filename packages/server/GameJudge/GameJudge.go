package main

import (
	"fmt"
	"server/GameJudge/internal/DataOperator"
)

func main() {
	fmt.Print(DataOperator.GetOriginalGameMap(0))
}
