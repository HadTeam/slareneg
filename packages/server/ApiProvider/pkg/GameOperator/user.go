package GameOperator

import "server/GameJudge/pkg/GameType"

func getUserFromList(userList *[]GameType.User, userId uint8) *GameType.User {
	for i, u := range *userList {
		if u.UserId == userId {
			return &(*userList)[i]
		}
	}
	return nil
}

func UserJoin(id GameType.GameId, user GameType.User) bool {
	game := data.GetCurrentGame(id)
	if game.Status == GameType.GameStatusWaiting {
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
	if game.Status == GameType.GameStatusRunning {
		if u := getUserFromList(&game.UserList, user.UserId); u != nil && u.Status == GameType.UserStatusDisconnected {
			u.Status = GameType.UserStatusConnected
			return true
		}
	}
	return false
}

func UserQuit(id GameType.GameId, user GameType.User) bool {
	game := data.GetCurrentGame(id)
	if game.Status == GameType.GameStatusWaiting {
		for i, u := range game.UserList {
			if u.UserId == user.UserId {
				game.UserList = append(game.UserList[i:], game.UserList[:i+1]...)
				return true
			}
		}
	}
	if game.Status == GameType.GameStatusRunning {
		if u := getUserFromList(&game.UserList, user.UserId); u != nil && u.Status == GameType.UserStatusConnected {
			u.Status = GameType.UserStatusDisconnected
			return true
		}
	}
	return false
}
