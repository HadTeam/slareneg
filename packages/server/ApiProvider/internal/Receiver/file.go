package Receiver

import (
	"context"
	"fmt"
	"os"
	"server/ApiProvider/internal/CommandPauser"
	"server/Utils/pkg/DataSource"
	"server/Utils/pkg/GameType"
	"strconv"
	"strings"
	"time"
)

// NOTE: DEBUG ONLY
// Use to receive instructions from local file, in order to test the game functions

var data DataSource.TempDataSource
var fileDir = "./test/replay"

func ApplyDataSource(source any) {
	data = source.(DataSource.TempDataSource)
}

func NewFileReceiver() {
	f := LoadFile()

	selectWaitingGame := func() []GameType.Game {
		var gl []GameType.Game
		for {
			if gl = data.GetGameList(GameType.GameMode1v1); gl != nil {
				break
			}
		}
		var ret []GameType.Game
		for _, g := range gl {
			if g.Status == GameType.GameStatusWaiting {
				ret = append(ret, g)
			}
		}
		return ret
	}
	for _, r := range f {
		time.Sleep(time.Millisecond * 200)
		fmt.Printf("Start game by reply file, index %d\n", r.Id)
		g := selectWaitingGame()[0]

		for _, c := range r.UserPack {
			ctx := Context{
				Context: context.TODO(),
				Game:    &g,
				User:    c.User,
				Command: make(chan string),
				Message: make(chan string),
			}
			go receiver(&ctx)
			go fakePlayer(&ctx, c.Ins)
		}
	}

}

type command struct {
	User GameType.User
	Ins  []string
}

type reply struct {
	Id       GameType.GameId
	UserPack []command
}

func LoadFile() []reply {
	//time.Sleep(1 * time.Second)
	var ret []reply
	dir, err := os.ReadDir(fileDir)
	if err != nil {
		panic("file dir cannot be visited")
	}
	for _, c := range dir {
		if c.IsDir() {
			continue
		}
		s := strings.Split(c.Name(), ".")
		if len(s) != 3 {
			continue
		}
		if s[1] != "gioreplay" || s[2] != "processed" {
			continue
		}

		id, _ := strconv.Atoi(s[0])
		r := reply{
			Id:       GameType.GameId(id),
			UserPack: []command{},
		}

		fileBuf, err := os.ReadFile(fileDir + "/" + c.Name())
		if err != nil {
			panic(fmt.Errorf("cannot read file %s", fileDir+"/"+c.Name()))
			return nil
		}
		userPart := strings.Split(string(fileBuf), "|")
		for userId, p := range userPart { // set the index of user part as user id
			if userId == 0 {
				continue
			}

			t := strings.Split(p, ":")
			name := t[0]
			cmdStr := strings.Split(t[1], "\n")

			cmd := command{
				User: GameType.User{
					Name:             name,
					UserId:           uint16(userId) - 1,
					Status:           GameType.UserStatusConnected,
					TeamId:           uint8(userId) - 1,
					ForceStartStatus: false,
				},
				Ins: append([]string{""}, cmdStr...),
			}

			r.UserPack = append(r.UserPack, cmd)
		}

		ret = append(ret, r)
	}
	return ret
}

func fakePlayer(ctx *Context, c []string) {
	ticker := time.NewTicker(50 * time.Millisecond)
	currentRound := uint16(0)
	for {
		select {
		case <-ticker.C:
			{
				if ctx.Game.RoundNum > currentRound {
					currentRound = ctx.Game.RoundNum
					ctx.Command <- c[currentRound]
				}
			}
		case msg := <-ctx.Message:
			{
				fmt.Printf("Game %d User %s Msg: %s\n", ctx.Game.Id, ctx.User.Name, msg)
			}
		}
	}
}

func receiver(ctx *Context) {
	ctx.User.Name = strconv.Itoa(int(ctx.User.UserId)) // TODO DEBUG ONLY
	ctx.User.Status = GameType.UserStatusConnected
	data.SetUserStatus(ctx.Game.Id, ctx.User)

	fmt.Printf("Game %d User %s join\n", ctx.Game.Id, ctx.User.Name)

	defer func() {
		ctx.User.Status = GameType.UserStatusDisconnected
		data.SetUserStatus(ctx.Game.Id, ctx.User)
	}()

	data.SetUserStatus(ctx.Game.Id, ctx.User)

	ticker := time.NewTicker(50 * time.Millisecond)
	for i := 1; true; i++ {
		select {
		case <-ctx.Done():
			{
				fmt.Println("ctx done")
				return
			}
		case cmd := <-ctx.Command:
			{
				if cmd == "" {
					continue
				}
				ins, err := CommandPauser.PauseCommandStr(ctx.User.UserId, cmd)
				if err != nil {
					panic(fmt.Errorf("cannot parse command: %s", cmd))
				}

				fmt.Println(ctx.Game.Id, ctx.User, ins)

				data.UpdateInstruction(ctx.Game.Id, ctx.User, ins)
			}
		case <-ticker.C:
			{
				done := func(g *GameType.Game, d string) {
					res := CommandPauser.GenerateMessage(d, ctx.Game.Id, ctx.User.UserId)
					ctx.Message <- res
					ctx.Game = g
				}

				// Check game status
				g := data.GetGameInfo(ctx.Game.Id)
				if i%100 == 0 {
					fmt.Println(ctx.User.UserId, i, ctx.Game.RoundNum, g.RoundNum)
				}
				if g.Status != ctx.Game.Status {
					if ctx.Game.Status == GameType.GameStatusWaiting && g.Status == GameType.GameStatusRunning {
						done(g, "start")
						continue
					} else if ctx.Game.Status == GameType.GameStatusRunning && g.Status == GameType.GameStatusEnd {
						done(g, "end")
						ticker.Stop()
					}
				} else if ctx.Game.Status == GameType.GameStatusWaiting && i%20 == 0 {
					done(g, "wait")
				}
				if ctx.Game.RoundNum != g.RoundNum {
					done(g, "newTurn")
				}
			}
		}
	}
}
