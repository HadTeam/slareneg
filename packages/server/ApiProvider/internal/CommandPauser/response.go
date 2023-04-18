package CommandPauser

import (
	"encoding/json"
	"server/Utils/pkg/DataSource"
	"server/Utils/pkg/GameType"
	"server/Utils/pkg/MapType"
)

var data DataSource.TempDataSource

func ApplyDataSource(source any) {
	data = source.(DataSource.TempDataSource)
}

func getVisibility(id GameType.GameId, userId uint16) [][]bool {
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

	ret := make([][]bool, m.Size.Y)

	light := func(x uint8, y uint8) {
		lightRange := []struct {
			x, y int
		}{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
		for _, r := range lightRange {
			ret[uint8(int(y)+r.y)][uint8(int(x)+r.x)] = true
		}
	}

	for rowNum, row := range m.Blocks {
		ret[rowNum] = make([]bool, m.Size.X)
		for colNum, b := range row {
			o := b.GetOwnerId()
			for _, u := range teamUsers {
				if u == o {
					light(uint8(colNum), uint8(rowNum))
				}
			}
		}
	}

	return ret
}

func getProcessedMap(id GameType.GameId, userId uint16, m *MapType.Map) [][][]uint16 {
	var mr [][][]uint16
	vis := getVisibility(id, userId)
	for rowNum, row := range m.Blocks {
		mr[rowNum] = [][]uint16{}
		for colNum, b := range row {
			if vis[rowNum][colNum] {
				mr[rowNum][colNum] = []uint16{uint16(b.GetMeta().BlockId), b.GetOwnerId(), b.GetNumber()}
			} else {
				const noOwner = uint16(0)
				mr[rowNum][colNum] = []uint16{uint16(b.GetMeta().VisitFallBackType), noOwner, 0}
			}
		}
	}
	return mr
}

func generateResponse(_type string, id GameType.GameId, userId uint16) string {
	switch _type {
	case "start":
		{
			m := data.GetCurrentMap(id)
			res := struct {
				Action    string       `json:"action"`
				MapWidth  uint8        `json:"mapWidth"`
				MapHeight uint8        `json:"mapHeight"`
				Map       [][][]uint16 `json:"map"`
			}{"start", m.Size.X, m.Size.Y, getProcessedMap(id, userId, m)}
			ret, _ := json.Marshal(res)
			return string(ret)
		}
	case "wait":
		{
			g := data.GetGameInfo(id)
			g.UserList = data.GetCurrentUserList(id)
			res := struct {
				Action   string            `json:"action"`
				Players  []GameType.User   `json:"players"`
				GameMode GameType.GameMode `json:"gameMode"`
			}{"wait", g.UserList, g.Mode}
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
	case "newTurn":
		{
			g := data.GetGameInfo(id)
			m := data.GetCurrentMap(id)
			res := struct {
				Action     string
				TurnNumber uint16
				Map        [][][]uint16 `json:"map"`
			}{"newTurn", g.RoundNum, getProcessedMap(id, userId, m)}
			ret, _ := json.Marshal(res)
			return string(ret)
		}
	default:
		{
			return "" // TODO
		}
	}
}
