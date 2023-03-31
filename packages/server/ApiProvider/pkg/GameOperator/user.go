package GameOperator

import (
	GameType2 "server/Untils/pkg/GameType"
)

func getUserFromList(userList *[]GameType2.User, userId uint8) *GameType2.User {
	for i, u := range *userList {
		if u.UserId == userId {
			return &(*userList)[i]
		}
	}
	return nil
}

// UserJoin TODO: Add unit test
func UserJoin(id GameType2.GameId, user GameType2.User) bool {
	game := data.GetCurrentGame(id)
	if game.Status == GameType2.GameStatusWaiting {
		game.UserList = append(game.UserList, user)

		if l := uint8(len(game.UserList)); l >= game.Mode.MinUserNum {
			if l == game.Mode.MaxUserNum { // Start game
				StartGame(id)
			} else { // Try to force start
				TryForceStart(id)
			}
		}

		return true
	}
	if game.Status == GameType2.GameStatusRunning {
		if u := getUserFromList(&game.UserList, user.UserId); u != nil && u.Status == GameType2.UserStatusDisconnected {
			u.Status = GameType2.UserStatusConnected
			return true
		}
	}
	return false
}

// UserQuit TODO: Add unit test
func UserQuit(id GameType2.GameId, user GameType2.User) bool {
	game := data.GetCurrentGame(id)
	if game.Status == GameType2.GameStatusWaiting {
		for i, u := range game.UserList {
			if u.UserId == user.UserId {
				game.UserList = append(game.UserList[i:], game.UserList[:i+1]...)
				return true
			}
		}
	}
	if game.Status == GameType2.GameStatusRunning {
		if u := getUserFromList(&game.UserList, user.UserId); u != nil && u.Status == GameType2.UserStatusConnected {
			u.Status = GameType2.UserStatusDisconnected
			return true
		}
	}
	return false
}
