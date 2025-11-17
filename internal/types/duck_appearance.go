package types

type SkinType string      // @name SkinType
type AccessoryType string // @name AccessoryType

const (
	SkinGiraffe  SkinType = "giraffe"
	SkinLGBT     SkinType = "lgbt"
	SkinSuperman SkinType = "superman"
)

const (
	AccessoryFlowerCrown  AccessoryType = "flower_crown"
	AccessoryKingCrown    AccessoryType = "king_crown"
	AccessorySupermanCape AccessoryType = "superman_cape"
	AccessoryVespaHelmet  AccessoryType = "vespa_helmet"
)

type DuckAppearance struct {
	Skin        SkinType        `json:"skin"`
	Accessories []AccessoryType `json:"accessories"`
} // @name DuckAppearance
