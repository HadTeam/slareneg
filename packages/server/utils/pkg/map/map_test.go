package _map

import (
	"reflect"
	"server/utils/pkg/map/blockManager"
	"server/utils/pkg/map/blockManager/block"
	_ "server/utils/pkg/map/blockManager/block"
	"server/utils/pkg/map/type"
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
				jsonStr: `{"mappings":{"block":["blank"]},"type":[[0,0,0]]}`,
			},
			want: &Map{
				[][]_type.Block{
					{
						&block.Blank{},
						&block.Blank{},
						&block.Blank{},
					},
				},
				mapInfo{
					size: MapSize{3, 1},
					id:   0,
				},
			},
		},
		{
			name: "with not expected owner field 1",
			args: args{
				jsonStr: `{"mappings":{"block":["blank"],"owner":[1,2]},"type":[[0,0,0]]}`,
			},
			want: &Map{
				[][]_type.Block{
					{
						&block.Blank{},
						&block.Blank{},
						&block.Blank{},
					},
				},
				mapInfo{
					size: MapSize{3, 1},
					id:   0,
				},
			},
		},
		{
			name: "with not expected owner field 2",
			args: args{
				jsonStr: `{"mappings":{"block":["blank"]},"type":[[0,0,0]],"owner":[[1,1,1]]}`,
			},
			want: &Map{
				[][]_type.Block{
					{
						&block.Blank{},
						&block.Blank{},
						&block.Blank{},
					},
				},
				mapInfo{
					size: MapSize{3, 1},
					id:   0,
				},
			},
		},
		{
			name: "with owner field",
			args: args{
				jsonStr: `{"mappings":{"block":["soldier"],"owner":[1]},"type":[[0,0,0]],"owner":[[1,1,1]]}`,
			},
			want: &Map{
				[][]_type.Block{
					{
						blockManager.NewBlock(1, 0, 1),
						blockManager.NewBlock(1, 0, 1),
						blockManager.NewBlock(1, 0, 1),
					},
				},
				mapInfo{
					size: MapSize{3, 1},
					id:   0,
				},
			},
		},
		{
			name: "with owner field and number field",
			args: args{
				jsonStr: `{"mappings":{"block":["soldier"],"owner":[1]},"type":[[0,0,0]],"owner":[[1,1,1]],"number":[[1,2,255]]}`,
			},
			want: &Map{
				[][]_type.Block{
					{
						blockManager.NewBlock(1, 1, 1),
						blockManager.NewBlock(1, 2, 1),
						blockManager.NewBlock(1, 255, 1),
					},
				},
				mapInfo{
					size: MapSize{3, 1},
					id:   0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run("JsonStrToMap: "+tt.name, func(t *testing.T) {
			if got := JsonStrToMap(tt.args.jsonStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonStrToMap() = %v, want %v", got, tt.want)
			}
		})
	}
	for _, tt := range tests {
		t.Run("MapToJsonStr: "+tt.name, func(t *testing.T) {
			got := MapToJsonStr(tt.want)
			var j1, j2 *Map
			j1 = JsonStrToMap(got)
			j2 = JsonStrToMap(tt.args.jsonStr)
			if !reflect.DeepEqual(j1, j2) {
				t.Errorf("MapToJsonStr() = %v, want %v", j1, j2)
			}
		})
	}
}
