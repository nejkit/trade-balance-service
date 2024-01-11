package staticserr

import "errors"

var (
	ErrorNotFoundAsset         = errors.New("AssetNotFound")
	ErrorNotEnoughBalance      = errors.New("NotEnoughBalance")
	ErrorNotExistsCurrency     = errors.New("CurrencyNotExists")
	ErrorCurrencyAlreadyExists = errors.New("CurrencyAlreadyExists")
	ErrorRabbitConnectionFail  = errors.New("ErrorRabbitConnectionFail")
)
