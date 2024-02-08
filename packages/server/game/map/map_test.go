package _map

import (
	"github.com/davecgh/go-spew/spew"
	"reflect"
	"server/game/block"
	_ "server/game/block"
	"server/game/user"
	"testing"
)

func TestConvJsonStrMap(t *testing.T) {
	type args struct {
		jsonStr string
	}
	tests := []struct {
		name string
		args args
		want *Map
	}{
		{
			name: "basic map string",
			args: args{
				jsonStr: `{"mappings":{"block":["blank"]},"game_def":[[0,0,0]]}`,
			},
			want: New(0, [][]block.Block{
				{
					&block.Blank{},
					&block.Blank{},
					&block.Blank{},
				},
			}),
		},
		{
			name: "with not expected owner field 1",
			args: args{
				jsonStr: `{"mappings":{"block":["blank"],"owner":[1,2]},"game_def":[[0,0,0]]}`,
			},
			want: New(0, [][]block.Block{
				{
					&block.Blank{},
					&block.Blank{},
					&block.Blank{},
				},
			}),
		},
		{
			name: "with not expected owner field 2",
			args: args{
				jsonStr: `{"mappings":{"block":["blank"]},"game_def":[[0,0,0]],"owner":[[1,1,1]]}`,
			},
			want: New(0, [][]block.Block{
				{
					&block.Blank{},
					&block.Blank{},
					&block.Blank{},
				},
			}),
		},
		{
			name: "with owner field",
			args: args{
				jsonStr: `{"mappings":{"block":["soldier"],"owner":[1]},"game_def":[[0,0,0]],"owner":[[1,1,1]]}`,
			},
			want: New(0, [][]block.Block{
				{
					block.NewBlock(1, 0, 1),
					block.NewBlock(1, 0, 1),
					block.NewBlock(1, 0, 1),
				},
			}),
		},
		{
			name: "with owner field and number field",
			args: args{
				jsonStr: `{"mappings":{"block":["soldier"],"owner":[1]},"game_def":[[0,0,0]],"owner":[[1,1,1]],"number":[[1,2,255]]}`,
			},
			want: New(0, [][]block.Block{
				{
					block.NewBlock(1, 1, 1),
					block.NewBlock(1, 2, 1),
					block.NewBlock(1, 255, 1),
				},
			}),
		},
	}
	for _, tt := range tests {
		t.Run("JsonStrToMap:"+tt.name, func(t *testing.T) {
			if got := JsonStrToMap(tt.args.jsonStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonStrToMap() = \n%s, want \n%s", spew.Sdump(got), spew.Sdump(tt.want))
			}
		})
	}
	for _, tt := range tests {
		t.Run("MapToJsonStr:"+tt.name, func(t *testing.T) {
			got := MapToJsonStr(tt.want)
			var j1, j2 *Map
			j1 = JsonStrToMap(got)
			j2 = JsonStrToMap(tt.args.jsonStr)
			if !reflect.DeepEqual(j1, j2) {
				t.Errorf("MapToJsonStr() = \n%s, want \n%s\nStr: \n%s", spew.Sdump(j1), spew.Sdump(j2), spew.Sdump(got))
			}
		})
	}
}

func TestMap_WarFog(t *testing.T) {
	type args struct {
		gameId   uint16
		teamId   uint8
		userList []user.User
	}
	otherBlock := block.NewBlock(block.SoldierMeta.BlockId, 111, 3)
	soldierBlock := block.NewBlock(block.SoldierMeta.BlockId, 10, 1)
	otherBlockData := []uint16{1, 3, 111}
	soldierBlockData := []uint16{1, 1, 10}
	blockDisallowedData := []uint16{0, 0, 0}
	map9x9center := New(0, [][]block.Block{
		{otherBlock, otherBlock, otherBlock},
		{otherBlock, soldierBlock, otherBlock},
		{otherBlock, otherBlock, otherBlock},
	})
	tests1 := []struct {
		name string
		m    *Map
		args args
		want *visibilityInFog
	}{
		{name: "9x9(center) block owner", m: map9x9center, args: args{gameId: 0, teamId: 1, userList: []user.User{{UserId: 1, TeamId: 1}}},
			want: &visibilityInFog{
				{false, true, false},
				{true, true, true},
				{false, true, false},
			},
		},
		{name: "9x9(center) team member", m: map9x9center, args: args{gameId: 0, teamId: 1, userList: []user.User{{UserId: 1, TeamId: 1}, {UserId: 2, TeamId: 1}}},
			want: &visibilityInFog{
				{false, true, false},
				{true, true, true},
				{false, true, false},
			},
		},
		{name: "9x9(center) others", m: map9x9center, args: args{gameId: 0, teamId: 2, userList: []user.User{{UserId: 1, TeamId: 1}, {UserId: 2, TeamId: 2}}},
			want: &visibilityInFog{
				{false, false, false},
				{false, false, false},
				{false, false, false},
			},
		},
		{name: "9x9(corner) block owner",
			m: New(0, [][]block.Block{
				{soldierBlock, otherBlock, otherBlock},
				{otherBlock, otherBlock, otherBlock},
				{otherBlock, otherBlock, soldierBlock},
			}),
			args: args{gameId: 0, teamId: 1, userList: []user.User{{UserId: 1, TeamId: 1}}},
			want: &visibilityInFog{
				{true, true, false},
				{true, false, true},
				{false, true, true},
			},
		},
	}
	for _, tt := range tests1 {
		t.Run("visibilityInFog "+tt.name, func(t *testing.T) {
			if got := tt.m.getVisibilityInFog(tt.args.gameId, tt.args.teamId, tt.args.userList); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getVisibilityInFog() = %v, want %v", got, tt.want)
			}
		})
	}

	tests2 := []struct {
		name string
		m    *Map
		args args
		want *mapWithFog
	}{
		{name: "9x9(center) block owner", m: map9x9center, args: args{gameId: 0, teamId: 1, userList: []user.User{{UserId: 1, TeamId: 1}}},
			want: &mapWithFog{
				{blockDisallowedData, otherBlockData, blockDisallowedData},
				{otherBlockData, soldierBlockData, otherBlockData},
				{blockDisallowedData, otherBlockData, blockDisallowedData},
			},
		},
		{name: "9x9(center) team member", m: map9x9center, args: args{gameId: 0, teamId: 1, userList: []user.User{{UserId: 1, TeamId: 1}, {UserId: 2, TeamId: 1}}},
			want: &mapWithFog{
				{blockDisallowedData, otherBlockData, blockDisallowedData},
				{otherBlockData, soldierBlockData, otherBlockData},
				{blockDisallowedData, otherBlockData, blockDisallowedData},
			},
		},
		{name: "9x9(center) others", m: map9x9center, args: args{gameId: 0, teamId: 2, userList: []user.User{{UserId: 1, TeamId: 1}, {UserId: 2, TeamId: 2}}},
			want: &mapWithFog{
				{blockDisallowedData, blockDisallowedData, blockDisallowedData},
				{blockDisallowedData, blockDisallowedData, blockDisallowedData},
				{blockDisallowedData, blockDisallowedData, blockDisallowedData},
			},
		},
		{name: "9x9(corner) block owner", args: args{gameId: 0, teamId: 1, userList: []user.User{{UserId: 1, TeamId: 1}}},
			m: New(0, [][]block.Block{
				{soldierBlock, otherBlock, otherBlock},
				{otherBlock, otherBlock, otherBlock},
				{otherBlock, otherBlock, soldierBlock},
			}),
			want: &mapWithFog{
				{soldierBlockData, otherBlockData, blockDisallowedData},
				{otherBlockData, blockDisallowedData, otherBlockData},
				{blockDisallowedData, otherBlockData, soldierBlockData},
			},
		},
	}
	for _, tt := range tests2 {
		t.Run("GetProcessedMap "+tt.name, func(t *testing.T) {
			if got := tt.m.GetProcessedMap(tt.args.gameId, tt.args.teamId, tt.args.userList); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProcessedMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
