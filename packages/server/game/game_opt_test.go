package game

import (
	"server/game/block"
	_ "server/game/block"
	"server/game/map"
	"server/game/user"
	"testing"
)

func TestGame_allocateKing(t *testing.T) {
	tests := []struct {
		name      string
		g         *Game
		wantKings int
	}{
		{
			name: "All users already have a King",
			g: &Game{
				UserList: []user.User{{UserId: 1}, {UserId: 2}},
				Map: _map.New(0, [][]block.Block{
					{
						block.NewBlock(block.KingMeta.BlockId, 0, 1),
						block.NewBlock(block.KingMeta.BlockId, 0, 2),
					}}),
			},

			wantKings: 2,
		},
		{
			name: "Some users do not have a King",
			g: &Game{
				UserList: []user.User{{UserId: 1}, {UserId: 2}},
				Map: _map.New(0, [][]block.Block{
					{
						block.NewBlock(block.KingMeta.BlockId, 0, 0),
						block.NewBlock(block.KingMeta.BlockId, 0, 2),
					},
				}),
			},

			wantKings: 2,
		},
		{
			name: "No users have a King",
			g: &Game{
				UserList: []user.User{{UserId: 1}, {UserId: 2}, {UserId: 3}},
				Map: _map.New(0, [][]block.Block{
					{
						block.NewBlock(block.KingMeta.BlockId, 0, 0),
						block.NewBlock(block.KingMeta.BlockId, 0, 0),
						block.NewBlock(block.KingMeta.BlockId, 0, 0),
					},
				}),
			},
			wantKings: 3,
		},
		{
			name: "Transform King to Castle",
			g: &Game{
				UserList: []user.User{{UserId: 1}, {UserId: 2}},
				Map: _map.New(0, [][]block.Block{
					{
						block.NewBlock(block.KingMeta.BlockId, 0, 0),
						block.NewBlock(block.KingMeta.BlockId, 0, 0),
						block.NewBlock(block.KingMeta.BlockId, 0, 0),
					},
				}),
			},
			wantKings: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.g.allocateKing()
			gotKings := len(tt.g.getKingPos())
			if gotKings != tt.wantKings {
				t.Errorf("allocateKing() = %v, want %v", gotKings, tt.wantKings)
			}
		})
	}
}
