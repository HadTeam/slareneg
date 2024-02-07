package judge

import (
	"server/game"
	"server/game/block"
	_ "server/game/block"
	"server/game/map"
	"testing"
)

func Test_allocateKing(t *testing.T) {
	type args struct {
		ctx *gameContext
	}

	tests := []struct {
		name      string
		args      args
		wantKings int // Expected number of Kings after allocation
	}{
		{
			name: "All users already have a King",
			args: args{
				ctx: &gameContext{
					g: &game.Game{
						UserList: []game.User{{UserId: 1}, {UserId: 2}},
						Map: _map.CreateMapWithBlocks(0, [][]block.Block{
							{
								block.NewBlock(block.KingMeta.BlockId, 0, 1),
								block.NewBlock(block.KingMeta.BlockId, 0, 2),
							}}),
					},
					kingPos: []block.Position{{X: 1, Y: 1}, {X: 2, Y: 1}},
				},
			},
			wantKings: 2,
		},
		{
			name: "Some users do not have a King",
			args: args{
				ctx: &gameContext{
					g: &game.Game{
						UserList: []game.User{{UserId: 1}, {UserId: 2}},
						Map: _map.CreateMapWithBlocks(0, [][]block.Block{
							{
								block.NewBlock(block.KingMeta.BlockId, 0, 0),
								block.NewBlock(block.KingMeta.BlockId, 0, 2),
							},
						}),
					},
					kingPos: []block.Position{{X: 1, Y: 1}, {X: 2, Y: 1}},
				},
			},
			wantKings: 2,
		},
		{
			name: "No users have a King",
			args: args{
				ctx: &gameContext{
					g: &game.Game{
						UserList: []game.User{{UserId: 1}, {UserId: 2}, {UserId: 3}},
						Map: _map.CreateMapWithBlocks(0, [][]block.Block{
							{
								block.NewBlock(block.KingMeta.BlockId, 0, 0),
								block.NewBlock(block.KingMeta.BlockId, 0, 0),
								block.NewBlock(block.KingMeta.BlockId, 0, 0),
							},
						}),
					},
					kingPos: []block.Position{{X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1}},
				},
			},
			wantKings: 3,
		},
		{
			name: "Transform King to Castle",
			args: args{
				ctx: &gameContext{
					g: &game.Game{
						UserList: []game.User{{UserId: 1}, {UserId: 2}},
						Map: _map.CreateMapWithBlocks(0, [][]block.Block{
							{
								block.NewBlock(block.KingMeta.BlockId, 0, 0),
								block.NewBlock(block.KingMeta.BlockId, 0, 0),
								block.NewBlock(block.KingMeta.BlockId, 0, 0),
							},
						}),
					},
					kingPos: []block.Position{{X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1}},
				},
			},
			wantKings: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allocateKing(tt.args.ctx)
			gotKings := len(tt.args.ctx.kingPos)
			if gotKings != tt.wantKings {
				t.Errorf("allocateKing() = %v, want %v", gotKings, tt.wantKings)
			}
			realGetKings := len(getKingPos(tt.args.ctx.g))
			if realGetKings != tt.wantKings {
				t.Errorf("allocateKing() real = %v, want %v", realGetKings, tt.wantKings)
			}
		})
	}
}
