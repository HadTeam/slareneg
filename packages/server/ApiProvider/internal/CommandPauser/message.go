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
			if o == 0 {
				continue
			}
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
	vis := getVisibility(id, userId)
	mr := make([][][]uint16, len(m.Blocks))
	for rowNum, row := range m.Blocks {
		mr[rowNum] = make([][]uint16, len(row))
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

type playerInfo struct {
	Name       string `json:"name"`
	Id         uint16 `json:"id"`
	ForceStart bool   `json:"forceStart"`
	TeamId     uint8  `json:"teamId"`
	Status     string `json:"status"`
}

func getUserList(id GameType.GameId) []playerInfo {
	l := data.GetCurrentUserList(id)
	ret := make([]playerInfo, len(l))
	var status string
	for i, u := range l {
		if u.Status == GameType.UserStatusConnected {
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

func GenerateMessage(_type string, id GameType.GameId, userId uint16) string {
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
				Action  string            `json:"action"`
				Players []playerInfo      `json:"players"`
				Mode    GameType.GameMode `json:"mode"`
			}{"info", getUserList(id), g.Mode}
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
