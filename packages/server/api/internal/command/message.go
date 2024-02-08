package command

import (
	"encoding/json"
	"server/game"
	"server/game/mode"
	"server/game/user"
)

var data game.TempDataSource

func ApplyDataSource(source any) {
	data = source.(game.TempDataSource)
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
		if u.Status == user.Connected {
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

func GenerateMessage(_type string, id game.Id, u user.User) string {
	switch _type {
	case "start":
		{
			m := data.GetCurrentMap(id)
			res := struct {
				Action    string       `json:"action"`
				MapWidth  uint8        `json:"mapWidth"`
				MapHeight uint8        `json:"mapHeight"`
				Map       [][][]uint16 `json:"map"`
			}{"start", m.Size().W, m.Size().H, *m.GetProcessedMap(uint16(id), u.TeamId, data.GetCurrentUserList(id))}
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
				Mode    mode.Mode    `json:"mode"`
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
			}{Action: "newTurn", TurnNumber: g.RoundNum, Map: *m.GetProcessedMap(uint16(id), u.TeamId, data.GetCurrentUserList(id))}
			ret, _ := json.Marshal(res)
			return string(ret)
		}
	default:
		{
			return "" // TODO
		}
	}
}
