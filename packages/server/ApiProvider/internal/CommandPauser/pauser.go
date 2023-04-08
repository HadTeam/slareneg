package CommandPauser

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"server/Untils/pkg/InstructionType"
	"strconv"
	"strings"
)

func PauseCommandStr(userId uint8, str string) (InstructionType.Instruction, error) {
	var err error = nil
	ret := InstructionType.Instruction(nil)
	args := strings.Split(str, " ")
	fmt.Printf("%#v\n", args)
	v := validator.New()
	switch args[0] {
	case "Move":
		{
			if len(args)-1 == 4 {
				c := moveCommand{
					X:       args[1],
					Y:       args[2],
					Towards: args[3],
					Number:  args[4],
				}
				err = v.Struct(c)
				if err == nil {
					n, _ := strconv.Atoi(c.Number)
					x, _ := strconv.Atoi(c.X)
					y, _ := strconv.Atoi(c.Y)
					ret = &InstructionType.Move{
						UserId:   userId,
						Position: InstructionType.BlockPosition{X: uint8(x), Y: uint8(y)},
						Towards:  InstructionType.MoveTowardsType(c.Towards),
						Number:   uint8(n),
					}
				} else {
					err = fmt.Errorf("illegally arguments\n")
				}
			} else {
				err = fmt.Errorf("argument number not right\n")
			}
			return ret, err
		}
	case "ForceStart":
		{
			if len(args)-1 == 1 {
				c := forceStartCommand{Status: args[1]}
				err = v.Struct(c)
				if err == nil {
					s := false
					if c.Status == "true" {
						s = true
					}
					ret = &InstructionType.ForceStart{
						UserId: userId,
						Status: s,
					}
				} else {
					err = fmt.Errorf("illegally arguments\n")
				}
			} else {
				err = fmt.Errorf("argument number not right")
			}
			return ret, err
		}
	case "Surrender":
		{
			if len(args)-1 == 0 {
				ret = &InstructionType.Surrender{UserId: userId}
			} else {
				err = fmt.Errorf("argument number not right")
			}
			return ret, err
		}
	default:
		return nil, fmt.Errorf("instruction not found\n")
	}
}
