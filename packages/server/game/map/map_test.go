package _map

import (
	"github.com/davecgh/go-spew/spew"
	"reflect"
	"server/game/block"
	_ "server/game/block"
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
			want: CreateMapWithBlocks(0, [][]block.Block{
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
			want: CreateMapWithBlocks(0, [][]block.Block{
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
			want: CreateMapWithBlocks(0, [][]block.Block{
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
			want: CreateMapWithBlocks(0, [][]block.Block{
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
			want: CreateMapWithBlocks(0, [][]block.Block{
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
				t.Errorf("MapToJsonStr() = \n%s, want \n%s", spew.Sdump(j1), spew.Sdump(j2))
			}
		})
	}
}
