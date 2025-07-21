package gamemap

import (
	"encoding/json"
	"server/internal/game/block"
)

// BlockDTO is a thin wrapper for JSON serialization of block.Block interface
type BlockDTO struct {
	Num   block.Num   `json:"num"`
	Owner block.Owner `json:"owner"`
	Meta  block.Meta  `json:"meta"`
}

type ExportedMap struct {
	Size   Size          `json:"size"`
	Info   Info          `json:"info"`
	Blocks [][]BlockDTO  `json:"blocks"`
}

func ExportToJSON(m Map) ([]byte, error) {
	// Convert block.Block interfaces to BlockDTOs
	blocks := m.Blocks()
	dtoBlocks := make([][]BlockDTO, len(blocks))
	
	for i, row := range blocks {
		dtoBlocks[i] = make([]BlockDTO, len(row))
		for j, b := range row {
			if b != nil {
				dtoBlocks[i][j] = BlockDTO{
					Num:   b.Num(),
					Owner: b.Owner(),
					Meta:  b.Meta(),
				}
			}
		}
	}
	
	e := ExportedMap{
		Size:   m.Size(),
		Info:   m.Info(),
		Blocks: dtoBlocks,
	}
	return json.MarshalIndent(e, "", "  ")
}

