package dto

type AssetState int

const (
	ACTIVE  AssetState = 0
	DELETED AssetState = 1
)

type TradeAsset struct {
	Id          string
	CreatedDate int64
	State       AssetState
}
