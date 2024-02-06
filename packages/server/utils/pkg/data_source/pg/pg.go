package pg

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"math/rand"
	"server/game_logic"
	"server/game_logic/game_def"
	"server/game_logic/map"
	data_source "server/utils/pkg/data_source"
	db "server/utils/pkg/pg"
)

var _ data_source.PersistentDataSource = (*Pg)(nil)
var _ data_source.TempDataSource = (*Pg)(nil)

type Pg struct {
}

const sqlCreateGame = "INSERT INTO game(game_id,mode,status,round_num,create_time,map,user_list) VALUES($1,$2,$3,$4,now(),$5,$6)"

func generatorMapJson(m *_map.Map) string {
	type b struct {
		TypeId  uint8  `json:"game_def"`
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
			ob := m.GetBlock(game_def.Position{X: x, Y: y})
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

var sqlQueryGame = "SELECT * FROM game WHERE game_id=$1"

func (p *Pg) CreateGame(mode game_def.Mode) game_logic.Id {
	var gameId game_logic.Id
	for {
		gameId = game_logic.Id(rand.Uint32())
		if ok := db.SqlQueryExist(sqlQueryGame, gameId); !ok && gameId >= 100 { // gameId 1-99 is for debugging usage
			break
		}
	}
	g := game_logic.Game{
		Mode:     mode,
		Id:       gameId,
		RoundNum: 0,
	}
	p.DebugCreateGame(&g)
	return gameId
}

func (p *Pg) DebugCreateGame(g *game_logic.Game) (ok bool) {
	r := db.SqlExec(sqlCreateGame, g.Id, g.Mode.NameStr, game_logic.StatusWaiting, g.RoundNum, generatorMapJson(g.Map), "[]")
	if row, err := r.RowsAffected(); err != nil || row != 1 {
		logrus.Warn("create game filed: ", err)
		return false
	} else {
		return true
	}
}

var sqlQueryGameList = "SELECT game_id FROM game WHERE mode=$1 AND (status=1 OR status=2)"

func (p *Pg) GetGameList(mode game_def.Mode) []game_logic.Game {
	r := db.SqlQuery(sqlQueryGameList, mode)
	var list []game_logic.Id
	for {
		var id game_logic.Id
		r.Next()
		if err := r.Scan(&id); err != nil {
			break
		}
		list = append(list, id)
	}
	ret := make([]game_logic.Game, len(list))
	for i, id := range list {
		ret[i] = *p.GetGameInfo(id)
	}
	return ret
}

func (p *Pg) CancelGame(id game_logic.Id) (ok bool) {
	//TODO implement me
	panic("implement me")
}

func (p *Pg) GetCurrentUserList(id game_logic.Id) []game_def.User {
	//TODO implement me
	panic("implement me")
}

func (p *Pg) GetInstructions(id game_logic.Id, tempId uint16) map[uint16]game_def.Instruction {
	//TODO implement me
	panic("implement me")
}

var sqlQueryGameInfo = "SELECT mode,status,round_num,create_time FROM game WHERE game_id=$1"

func (p *Pg) GetGameInfo(id game_logic.Id) *game_logic.Game {
	r := db.SqlQuery(sqlQueryGameInfo, id)
	defer func() { _ = r.Close() }()

	g := &game_logic.Game{
		Id: id,
	}

	var modeStr string

	r.Next()
	if err := r.Scan(&modeStr, &g.Status, &g.RoundNum, &g.CreateTime); err != nil {
		logrus.Warn("cannot get game info")
		return nil
	}
	if mode, ok := game_def.ModeMap[modeStr]; !ok {
		logrus.Warn("get unknown mode ", modeStr, " when get game info")
	} else {
		g.Mode = mode
	}
	return g
}

func (p *Pg) NewInstructionTemp(id game_logic.Id, tempId uint16) (ok bool) {
	//TODO implement me
	panic("implement me")
}

func (p *Pg) SetGameStatus(id game_logic.Id, status game_logic.Status) (ok bool) {
	//TODO implement me
	panic("implement me")
}

func (p *Pg) SetGameMap(id game_logic.Id, m *_map.Map) (ok bool) {
	//TODO implement me
	panic("implement me")
}

func (p *Pg) SetUserStatus(id game_logic.Id, user game_def.User) (ok bool) {
	//TODO implement me
	panic("implement me")
}

func (p *Pg) SetWinner(id game_logic.Id, teamId uint8) (ok bool) {
	//TODO implement me
	panic("implement me")
}

func (p *Pg) UpdateInstruction(id game_logic.Id, user game_def.User, instruction game_def.Instruction) (ok bool) {
	//TODO implement me
	panic("implement me")
}

func (p *Pg) GetCurrentMap(id game_logic.Id) *_map.Map {
	//TODO implement me
	panic("implement me")
}

func (p *Pg) GetOriginalMap(mapId uint32) *_map.Map {
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
