package utils

import (
	"trade-balance-service/external/balances"
	"trade-balance-service/staticserr"
)

func MapError(err error) balances.BpsErrorCode {
	switch err {
	case staticserr.ErrorNotFoundAsset:
		return balances.BpsErrorCode_BPS_ERROR_CODE_NOT_EXISTS_ASSET
	case staticserr.ErrorNotEnoughBalance:
		return balances.BpsErrorCode_BPS_ERROR_CODE_NOT_ENOUGH_BALANCE
	case staticserr.ErrorNotExistsCurrency:
		return balances.BpsErrorCode_BPS_ERROR_CODE_NOT_EXISTS_CURRENCY
	case staticserr.ErrorNotRelatedAccount:
		return balances.BpsErrorCode_BPS_ERROR_CODE_ASSET_NOT_RELATED_TO_ACCOUNT
	default:
		return balances.BpsErrorCode_BPS_ERROR_CODE_INTERNAL
	}
}
