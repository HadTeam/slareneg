package receiver

import (
	"context"
	"log"
	"os"
	_command "server/api/internal/command"
	"server/judgePool"
	"server/utils/pkg/datasource"
	"server/utils/pkg/game"
	"server/utils/pkg/map"
	"strconv"
	"strings"
	"time"
)

// NOTE: DEBUG ONLY
// Use to receive instructions from local file, in order to test the game functions

var data datasource.TempDataSource
var fileDir = "./test/replay"

func ApplyDataSource(source any) {
	data = source.(datasource.TempDataSource)
}

func NewFileReceiver(pool *judgePool.Pool) {
	f := LoadFile()

	for index, r := range f {
		time.Sleep(time.Millisecond * 200)
		log.Printf("Start game by reply file, index %d\n", r.Id)
		g := &game.Game{
			Map:        r.Map,
			Mode:       game.GameMode1v1,
			Id:         game.GameId(index + 1e3),
			UserList:   []game.User{},
			CreateTime: time.Now().UnixMicro(),
			Status:     game.GameStatusWaiting,
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
	User game.User
	Ins  []string
}

type reply struct {
	Id       game.GameId
	UserPack []command
	Map      *_map.Map
}

func LoadFile() []reply {
	var ret []reply
	dir, err := os.ReadDir(fileDir)
	if err != nil {
		panic("file dir cannot be visited")
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
			Id:       game.GameId(id),
			UserPack: []command{},
		}

		fileBuf, err := os.ReadFile(fileDir + "/" + c.Name())
		if err != nil {
			log.Panicf("cannot read file %s", fileDir+"/"+c.Name())
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
				User: game.User{
					Name:             name,
					UserId:           uint16(userId),
					Status:           game.UserStatusConnected,
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
	for {
		select {
		case <-ticker.C:
			{
				g := data.GetGameInfo(ctx.Game.Id)
				if currentRound >= uint16(len(c)) {
					log.Printf("[Game %d] Fake Player %s Command(tot: %d) runs out, quit\n", ctx.Game.Id, ctx.User.Name, len(c))
					ticker.Stop()
				}
				if ctx.Game.RoundNum > currentRound {
					log.Printf("[Game %d] Fake Player %s Discover new round %d\n", ctx.Game.Id, ctx.User.Name, g.RoundNum)
					currentRound = g.RoundNum
					if currentRound < uint16(len(c)) {
						log.Printf("[Game %d] Fake Player %s Send command '%s'\n", ctx.Game.Id, ctx.User.Name, c[currentRound])
						ctx.Command <- c[currentRound]
					}
				}
			}
		case msg := <-ctx.Message:
			{
				log.Printf("[Game %d] Fake Player %s Msg: %s\n", ctx.Game.Id, ctx.User.Name, msg[:30]) // output a part of message avoiding excessive output
			}
		}
	}
}

func receiver(ctx *Context) {
	ctx.User.Name = strconv.Itoa(int(ctx.User.UserId)) // DEBUG ONLY, avoiding strange username from `gioreply` file
	ctx.User.Status = game.UserStatusConnected
	data.SetUserStatus(ctx.Game.Id, ctx.User)

	log.Printf("[Game %d] User %s join\n", ctx.Game.Id, ctx.User.Name)

	defer func() {
		ctx.User.Status = game.UserStatusDisconnected
		data.SetUserStatus(ctx.Game.Id, ctx.User)
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
				if cmd == "" {
					continue
				}
				ins, err := _command.PauseCommandStr(ctx.User.UserId, cmd)
				if err != nil {
					log.Panicf("cannot parse command: %s", cmd)
				}
				data.UpdateInstruction(ctx.Game.Id, ctx.User, ins)
			}
		case <-ticker.C:
			{
				done := func(g *game.Game, d string) {
					res := _command.GenerateMessage(d, ctx.Game.Id, ctx.User.UserId)
					ctx.Message <- res
					ctx.Game = g
				}

				// Check game status
				g := data.GetGameInfo(ctx.Game.Id)
				if g.Status != ctx.Game.Status {
					if ctx.Game.Status == game.GameStatusWaiting && g.Status == game.GameStatusRunning {
						done(g, "info")
						done(g, "start")
						continue
					} else if g.Status == game.GameStatusEnd {
						done(g, "end")
						ticker.Stop()
						flag = false
						break
					}
				} else {
					if i%20 == 0 {
						done(g, "info")
						if ctx.Game.Status == game.GameStatusWaiting {
							done(g, "wait")
						}
					}

					if ctx.Game.Status == game.GameStatusRunning && ctx.Game.RoundNum != g.RoundNum {
						done(g, "newTurn")
						continue

					}
				}
			}
		}
	}
}
