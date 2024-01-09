package staticserr

import "errors"

var (
	ErrorNotFoundAsset    = errors.New("AssetNotFound")
	ErrorNotEnoughBalance = errors.New("NotEnoughBalance")
)
