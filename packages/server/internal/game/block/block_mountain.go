package block

type Mountain struct {
	BaseBlock
}

var MountainMeta Meta
var MountainName Name

func init() {
	MountainName = Register("mountain", "", toBlockMountain)
	MountainMeta = GetMetaByName[MountainName]
}

func toBlockMountain(Block) Block {
	return &Mountain{}
}

func (*Mountain) Meta() Meta {
	return MountainMeta
}

func (b *Mountain) AllowMove() AllowMove {
	return AllowMove{
		From:   false,
		To:     false,
		Reason: "Mountain block cannot be moved",
	}
}

func (block *Mountain) Fog(isOwner bool, isSight bool) Block {
	// Mountains are always visible as mountains (terrain feature)
	return block
}
