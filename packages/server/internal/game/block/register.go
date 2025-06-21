package block

import (
	"log/slog"
)

type tranFunc func(Block) Block

var transBlockTypeFunc map[Name]tranFunc
var GetMetaByName map[Name]Meta
var blockNames []Name

// Register registers a new block type and returns the Name for reference
func Register(name string, description string, transFunc tranFunc) Name {
	if transBlockTypeFunc == nil {
		transBlockTypeFunc = make(map[Name]tranFunc)
	}
	if GetMetaByName == nil {
		GetMetaByName = make(map[Name]Meta)
	}

	blockName := Name(name)

	meta := Meta{
		Name:        blockName,
		Description: description,
	}

	GetMetaByName[blockName] = meta
	transBlockTypeFunc[blockName] = transFunc
	blockNames = append(blockNames, blockName)

	slog.Debug("Registered a block", "name", meta.Name, "description", meta.Description)

	return blockName
}

// GetAllBlockNames returns all registered block names
func GetAllBlockNames() []Name {
	result := make([]Name, len(blockNames))
	copy(result, blockNames)
	return result
}

// GetBlockMeta returns the meta for a given block name
func GetBlockMeta(name Name) (Meta, bool) {
	meta, exists := GetMetaByName[name]
	return meta, exists
}

// BlockExists checks if a block type with the given name exists
func BlockExists(name Name) bool {
	_, exists := GetMetaByName[name]
	return exists
}

func ToBlockByName(name Name, block Block) Block {
	transFunc, exists := transBlockTypeFunc[name]
	if !exists {
		slog.Warn("Unknown block type", "name", name, "available", GetAllBlockNames())
		
		transFunc = transBlockTypeFunc["blank"] // Fallback to blank
	}
	return transFunc(block)
}

func NewBlock(name Name, num Num, owner Owner) Block {
	return ToBlockByName(name, &BaseBlock{num: num, owner: owner})
}
