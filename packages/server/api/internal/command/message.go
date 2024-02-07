package command

import (
	"encoding/json"
	"server/game"
	"server/game/block"
	"server/game/map"
	"server/utils/pkg/data_source"
	game_temp_pool "server/utils/pkg/game_temp_pool"
)

var data data_source.TempDataSource

func ApplyDataSource(source any) {
	data = source.(data_source.TempDataSource)
}

type visibilityArr [][]bool

func getVisibility(id game.Id, userId uint16) *visibilityArr {
	m := data.GetCurrentMap(id)
	ul := data.GetCurrentUserList(id)

	var teamId uint8
	var teamUsers []uint16
	for _, v := range ul {
		if v.UserId == userId {
			teamId = v.TeamId
		}
	}
	for _, v := range ul {
		if v.TeamId == teamId {
			teamUsers = append(teamUsers, v.UserId)
		}
	}

	var ret *visibilityArr
	if r, ok := game_temp_pool.Get(id, "visibility"); ok {
		ret = r.(*visibilityArr)
	} else {
		t := make(visibilityArr, m.Size().H)
		ret = &t
		game_temp_pool.Put(id, "visibility", &t)
	}

	for rowNum := uint8(0); rowNum <= m.Size().H-1; rowNum++ {
		(*ret)[rowNum] = make([]bool, m.Size().W)
	}

	light := func(x int, y int) {
		lightOffset := []struct {
			x, y int
		}{{0, 1}, {0, -1}, {-1, 0}, {1, 0}}
		for _, r := range lightOffset {
			ly := y + r.y
			lx := x + r.x
			if 0 <= ly && ly <= int(m.Size().H-1) && 0 <= lx && lx <= int(m.Size().W-1) {
				(*ret)[ly][lx] = true
			}
		}
	}

	for rowNum := uint8(0); rowNum <= m.Size().H-1; rowNum++ {
		for colNum := uint8(0); colNum <= m.Size().W-1; colNum++ {
			o := m.GetBlock(block.Position{X: colNum + 1, Y: rowNum + 1}).OwnerId()
			if o == 0 {
				continue
			}
			for _, u := range teamUsers {
				if u == o {
					light(int(colNum), int(rowNum))
				}
			}
		}
	}

	return ret
}

type processedMap [][][]uint16

func getProcessedMap(id game.Id, userId uint16, m *_map.Map) *processedMap {
	vis := getVisibility(id, userId)
	var ret *processedMap
	if r, ok := game_temp_pool.Get(id, "processedMap"); ok {
		ret = r.(*processedMap)
	} else {
		t := make(processedMap, m.Size().H)
		ret = &t
		game_temp_pool.Put(id, "processedMap", &t)
	}

	for rowNum := uint8(0); rowNum <= m.Size().H-1; rowNum++ {
		(*ret)[rowNum] = make([][]uint16, m.Size().W)
		for colNum := uint8(0); colNum <= m.Size().W-1; colNum++ {
			b := m.GetBlock(block.Position{X: colNum + 1, Y: rowNum + 1})
			if (*vis)[rowNum][colNum] {
				(*ret)[rowNum][colNum] = []uint16{uint16(b.Meta().BlockId), b.OwnerId(), b.Number()}
			} else {
				const noOwner = uint16(0)
				(*ret)[rowNum][colNum] = []uint16{uint16(b.Meta().VisitFallBackType), noOwner, 0}
			}
		}
	}
	return ret
}

type playerInfo struct {
	Name       string `json:"name"`
	Id         uint16 `json:"id"`
	ForceStart bool   `json:"forceStart"`
	TeamId     uint8  `json:"teamId"`
	Status     string `json:"status"`
}

func getUserList(id game.Id) []playerInfo {
	l := data.GetCurrentUserList(id)
	ret := make([]playerInfo, len(l))
	var status string
	for i, u := range l {
		if u.Status == game.UserStatusConnected {
			status = "connected"
		} else {
			status = "disconnect"
		}
		ret[i] = playerInfo{
			Name:       u.Name,
			Id:         u.UserId,
			ForceStart: u.ForceStartStatus,
			TeamId:     u.TeamId,
			Status:     status,
		}
	}
	return ret
}

func GenerateMessage(_type string, id game.Id, userId uint16) string {
	switch _type {
	case "start":
		{
			m := data.GetCurrentMap(id)
			res := struct {
				Action    string       `json:"action"`
				MapWidth  uint8        `json:"mapWidth"`
				MapHeight uint8        `json:"mapHeight"`
				Map       [][][]uint16 `json:"map"`
			}{"start", m.Size().W, m.Size().H, *getProcessedMap(id, userId, m)}
			ret, _ := json.Marshal(res)
			return string(ret)
		}
	case "wait":
		{
			res := struct {
				Action string `json:"action"`
			}{"wait"}
			ret, _ := json.Marshal(res)
			return string(ret)
		}
	case "end":
		{
			g := data.GetGameInfo(id)
			res := struct {
				Action string `json:"action"`
				Winner uint8  `json:"winnerTeam"`
			}{"end", g.Winner}
			ret, _ := json.Marshal(res)
			return string(ret)
		}
	case "info":
		{
			g := data.GetGameInfo(id)
			res := struct {
				Action  string       `json:"action"`
				Players []playerInfo `json:"players"`
				Mode    game.Mode    `json:"mode"`
			}{"info", getUserList(id), g.Mode}
			ret, _ := json.Marshal(res)
			return string(ret)
		}
	case "newTurn":
		{
			g := data.GetGameInfo(id)
			m := data.GetCurrentMap(id)
			res := struct {
				Action     string       `json:"action"`
				TurnNumber uint16       `json:"turnNumber"`
				Map        [][][]uint16 `json:"map"`
			}{"newTurn", g.RoundNum, *getProcessedMap(id, userId, m)}
			ret, _ := json.Marshal(res)
			return string(ret)
		}
	default:
		{
			return "" // TODO
		}
	}
}
