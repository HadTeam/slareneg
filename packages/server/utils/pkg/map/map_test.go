package _map

import (
	_ "server/utils/pkg/map/blockManager/block"
	"testing"
)

func TestMapToJsonStr(t *testing.T) {
	t.Run("t1", func(t *testing.T) {
		str := `
	{
		"mappings": {
			"block":[ 
				"blank",
				"castle",
				"king",
				"mountain",
				"soldier"
			],
			"owner": [
				1,
				2
			] 
		},
		"type": [
			[1, 0, 0, 0],
			[0, 2, 3, 0],
			[0, 3, 2, 0],
			[0, 0, 0, 1]
		],
		"owner": [
			[0, 0, 0, 0],
			[0, 1, 0, 0],
			[0, 0, 2, 0],
			[0, 0, 0, 0]
		],
		"number": [
			[43, 0, 0, 0],
			[0, 1, 0, 0],
			[0, 0, 1, 0],
			[0, 0, 0, 43]
		]
	}
`
		m := JsonStrToMap(str)
		got := MapToJsonStr(m)
		t.Log(got)
	})
}
