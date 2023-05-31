package pg

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	data_source "server/utils/pkg/datasource"
	"server/utils/pkg/game"
	"server/utils/pkg/instruction"
	_map "server/utils/pkg/map"
	"server/utils/pkg/map/block"
	db "server/utils/pkg/pg"
)

var _ data_source.PersistentDataSource = (*Pg)(nil)
var _ data_source.TempDataSource = (*Pg)(nil)

type Pg struct {
}

const sqlCreateGame = "INSERT INTO game(game_id,mode,status,round_num,map,user_list) VALUES($1,$2,$3,$4,$5,$6)"

func generatorMapJsonb(m _map.Map) string {
	type b struct {
		TypeId  uint8  `json:"type"`
		OwnerId uint16 `json:"owner"`
		Number  uint16 `json:"num"`
	}
	type jm struct {
		Id   uint32 `json:"id"`
		Size struct {
			Height uint8 `json:"height"`
			Width  uint8 `json:"width"`
		}
		Blocks [][]b `json:"blocks"`
	}
	ret := jm{
		Id: m.Id(),
		Size: struct {
			Height uint8 `json:"height"`
			Width  uint8 `json:"width"`
		}{m.Size().H, m.Size().W},
		Blocks: make([][]b, m.Size().H),
	}
	for y := uint8(1); y <= m.Size().H; y++ {
		ret.Blocks[y] = make([]b, m.Size().W)
		for x := uint8(1); x <= m.Size().W; x++ {
			ob := m.GetBlock(block.Position{X: x, Y: y})
			ret.Blocks[y][x] = b{
				TypeId:  ob.Meta().BlockId,
				OwnerId: ob.OwnerId(),
				Number:  ob.Number(),
			}
		}
	}

	str, err := json.Marshal(ret)
	if err != nil {
		logrus.Panic("Failed to marshal object to json: ", err)
	}
	return string(str)
}

func (p Pg) CreateGame(mode game.Mode) game.Id {
	//TODO implement me
	panic("implement me")
}

func (p Pg) DebugCreateGame(game *game.Game) (ok bool) {
	//TODO implement me
	panic("implement me")
}

func (p Pg) GetGameList(mode game.Mode) []game.Game {
	//TODO implement me
	panic("implement me")
}

func (p Pg) CancelGame(id game.Id) (ok bool) {
	//TODO implement me
	panic("implement me")
}

func (p Pg) GetCurrentUserList(id game.Id) []game.User {
	//TODO implement me
	panic("implement me")
}

func (p Pg) GetInstructions(id game.Id, tempId uint16) []instruction.Instruction {
	//TODO implement me
	panic("implement me")
}

func (p Pg) GetGameInfo(id game.Id) *game.Game {
	//TODO implement me
	panic("implement me")
}

func (p Pg) NewInstructionTemp(id game.Id, tempId uint16) (ok bool) {
	//TODO implement me
	panic("implement me")
}

func (p Pg) SetGameStatus(id game.Id, status game.Status) (ok bool) {
	//TODO implement me
	panic("implement me")
}

func (p Pg) SetGameMap(id game.Id, m *_map.Map) (ok bool) {
	//TODO implement me
	panic("implement me")
}

func (p Pg) SetUserStatus(id game.Id, user game.User) (ok bool) {
	//TODO implement me
	panic("implement me")
}

func (p Pg) SetWinner(id game.Id, teamId uint8) (ok bool) {
	//TODO implement me
	panic("implement me")
}

func (p Pg) UpdateInstruction(id game.Id, user game.User, instruction instruction.Instruction) (ok bool) {
	//TODO implement me
	panic("implement me")
}

func (p Pg) GetCurrentMap(id game.Id) *_map.Map {
	//TODO implement me
	panic("implement me")
}

func (p Pg) GetOriginalMap(mapId uint32) *_map.Map {
	sql := "SELECT map_str FROM original_map WHERE map_id=$1"
	r := db.SqlQuery(sql, mapId)
	defer func() {
		_ = r.Close()
	}()
	r.Next()
	var str string
	if err := r.Scan(&str); err != nil {
		return nil
	}
	return _map.Str2GameMap(mapId, str)
}
