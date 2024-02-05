package receiver

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	_command "server/api/internal/command"
	"server/game_logic"
	"server/game_logic/game_def"
	"server/game_logic/map"
	judge_pool "server/judge_pool"
	"server/utils/pkg/data_source"
	"strconv"
	"strings"
	"time"
)

// NOTE: DEBUG ONLY
// Use to receive instructions from local file, in order to test the game functions

var data data_source.TempDataSource
var fileDir = "./test/replay"

func ApplyDataSource(source any) {
	data = source.(data_source.TempDataSource)
}

func NewFileReceiver(pool *judge_pool.Pool) {
	f := LoadFile()

	for index, r := range f {
		time.Sleep(time.Millisecond * 200)
		logrus.Infof("start game by reply file #%d", r.Id)
		g := &game_logic.Game{
			Map:        r.Map,
			Mode:       _type.Mode1v1,
			Id:         game_logic.Id(index + 1e3),
			UserList:   []_type.User{},
			CreateTime: time.Now().UnixMicro(),
			Status:     game_logic.StatusWaiting,
			RoundNum:   0,
		}
		pool.DebugNewGame(g)
		for _, c := range r.UserPack {
			ctx := Context{
				Context: context.TODO(),
				Game:    g,
				User:    c.User,
				// it's TOO IMPORTANT to use channel with buffer, or the fake player and the receiver will strike now and then
				Command: make(chan string, 1),
				Message: make(chan string, 3),
			}
			go receiver(&ctx)
			go fakePlayer(&ctx, c.Ins)
		}
	}

}

type command struct {
	User _type.User
	Ins  []string
}

type reply struct {
	Id       game_logic.Id
	UserPack []command
	Map      *_map.Map
}

func LoadFile() []reply {
	var ret []reply
	dir, err := os.ReadDir(fileDir)
	if err != nil {
		logrus.Panic("file dir cannot be visited")
	}
	for index, c := range dir {
		if c.IsDir() {
			continue
		}
		s := strings.Split(c.Name(), ".")
		if len(s) != 3 {
			continue
		}
		if s[1] != "gioreplay" || s[2] != "processed" { // allows *.gioreplay.processed files
			continue
		}

		id, _ := strconv.Atoi(s[0])
		r := reply{
			Id:       game_logic.Id(id),
			UserPack: []command{},
		}

		fileBuf, err := os.ReadFile(fileDir + "/" + c.Name())
		if err != nil {
			logrus.Panicf("cannot read file %s", fileDir+"/"+c.Name())
			return nil
		}
		part := strings.Split(string(fileBuf), "|")

		m := part[0]
		r.Map = _map.FullStr2GameMap(uint32(index), m)

		for userId, p := range part { // set the index of user part as user id
			if userId == 0 {
				continue
			}

			t := strings.Split(p, ":")
			name := t[0]
			cmdStr := strings.Split(t[1], "\n")

			cmd := command{
				User: _type.User{
					Name:             name,
					UserId:           uint16(userId),
					Status:           _type.UserStatusConnected,
					TeamId:           uint8(userId) - 1,
					ForceStartStatus: false,
				},
				Ins: append([]string{""}, cmdStr...), // add an empty string to skip round 0
			}

			r.UserPack = append(r.UserPack, cmd)
		}

		ret = append(ret, r)
	}
	return ret
}

func fakePlayer(ctx *Context, c []string) {
	// DO NOT MODIFY `ctx`(except channel sending) 'cause it is read-only for a real player
	ticker := time.NewTicker(10 * time.Millisecond)
	currentRound := uint16(0)
	playerLogger := logrus.WithFields(logrus.Fields{
		"gameId": ctx.Game.Id,
		"user":   ctx.User.UserId,
	})
	for {
		select {
		case <-ticker.C:
			{
				g := data.GetGameInfo(ctx.Game.Id)
				if g.Status == game_logic.StatusEnd {
					return
				}

				if currentRound >= uint16(len(c)) {
					playerLogger.Infof("Command(tot: %d) runs out, quit", len(c))
					ticker.Stop()
					player := ctx.User
					player.Status = _type.UserStatusDisconnected
					data.SetUserStatus(ctx.Game.Id, player)
				}
				if ctx.Game.RoundNum > currentRound {
					playerLogger.Printf("Discover new round %d", g.RoundNum)
					currentRound = g.RoundNum
					if currentRound < uint16(len(c)) {
						playerLogger.Infof("Send command '%s'", c[currentRound])
						ctx.Command <- c[currentRound]
					}
				}
			}
		case msg := <-ctx.Message:
			{
				// output a part of message avoiding excessive output
				var m string
				if len(msg) <= 100 {
					m = msg
				} else {
					m = msg[:100]
				}
				logrus.Infof("Msg: %s", m)
			}
		}
	}
}

func receiver(ctx *Context) {
	//ctx.User.Name = strconv.Itoa(int(ctx.User.UserId)) // DEBUG ONLY, avoiding strange username from `gioreply` file
	ctx.User.Status = _type.UserStatusConnected
	data.SetUserStatus(ctx.Game.Id, ctx.User)

	receiverLogger := logrus.WithFields(logrus.Fields{
		"user": ctx.User.Name,
	})
	receiverLogger.Infof("user join")

	defer func() {
		ctx.User.Status = _type.UserStatusDisconnected
		data.SetUserStatus(ctx.Game.Id, ctx.User)
		receiverLogger.Infof("user quit")
	}()

	data.SetUserStatus(ctx.Game.Id, ctx.User)

	ticker := time.NewTicker(10 * time.Millisecond)
	flag := true
	for i := 1; flag; i++ {
		select {
		case <-ctx.Done():
			{
				return
			}
		case cmd := <-ctx.Command:
			{

				if strings.TrimSpace(cmd) == "" {
					continue
				}
				ins, err := _command.PauseCommandStr(ctx.User.UserId, cmd)
				if err != nil {
					receiverLogger.Panicf("cannot parse command: |%s|", cmd)
				}
				data.UpdateInstruction(ctx.Game.Id, ctx.User, ins)
			}
		case <-ticker.C:
			{
				done := func(g *game_logic.Game, d string) {
					res := _command.GenerateMessage(d, ctx.Game.Id, ctx.User.UserId)
					ctx.Message <- res
					ctx.Game = g
				}

				// Check game status
				g := data.GetGameInfo(ctx.Game.Id)
				if g.Status != ctx.Game.Status {
					if ctx.Game.Status == game_logic.StatusWaiting && g.Status == game_logic.StatusRunning {
						done(g, "info")
						done(g, "start")
						continue
					} else if g.Status == game_logic.StatusEnd {
						done(g, "end")
						ticker.Stop()
						flag = false
						break
					}
				} else {
					if i%20 == 0 {
						done(g, "info")
						if ctx.Game.Status == game_logic.StatusWaiting {
							done(g, "wait")
						}
					}

					if ctx.Game.Status == game_logic.StatusRunning && ctx.Game.RoundNum != g.RoundNum {
						done(g, "newTurn")
						continue
					}
				}
			}
		}
	}
}
