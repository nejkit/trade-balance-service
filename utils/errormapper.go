package utils

import (
	"trade-balance-service/external/bps"
	"trade-balance-service/staticserr"
)

func MapError(err error) bps.BpsErrorCode {
	switch err {
	case staticserr.ErrorNotFoundAsset:
		return bps.BpsErrorCode_BPS_ERROR_CODE_NOT_EXISTS_ASSET
	case staticserr.ErrorNotEnoughBalance:
		return bps.BpsErrorCode_BPS_ERROR_CODE_NOT_ENOUGH_BALANCE
	case staticserr.ErrorNotExistsCurrency:
		return bps.BpsErrorCode_BPS_ERROR_CODE_NOT_EXISTS_CURRENCY
	case staticserr.ErrorNotRelatedAccount:
		return bps.BpsErrorCode_BPS_ERROR_CODE_ASSET_NOT_RELATED_TO_ACCOUNT
	case staticserr.ErrorCurrencyAlreadyExists:
		return bps.BpsErrorCode_BPS_ERROR_CODE_NOT_EXISTS_CURRENCY
	default:
		return bps.BpsErrorCode_BPS_ERROR_CODE_INTERNAL
	}
}
