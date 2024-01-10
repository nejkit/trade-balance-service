package flow

import (
	"context"
	"trade-balance-service/dto"
	"trade-balance-service/external/balances"
	"trade-balance-service/utils"
)

type IAssetService interface {
	CreateNewAsset(ctx context.Context) (string, error)
	GetAssetInfoById(ctx context.Context, id string) (*dto.TradeAsset, error)
	DeactivateAsset(ctx context.Context, id string) error
}

type IBalanceService interface {
	EmmitBalance(ctx context.Context, assetId string, currencyCode string, amount float64) error
	AddCurrency(ctx context.Context, currencyName string, currencyCode string) error
	GetInfoAboutAssets(ctx context.Context, assetId string) ([]dto.BalanceModel, error)
}

type Flow struct {
	assetService   IAssetService
	balanceService IBalanceService
}

func NewFlow(assetService IAssetService, balanceService IBalanceService) *Flow {
	return &Flow{assetService: assetService, balanceService: balanceService}
}

func (f Flow) CreateAsset(ctx context.Context, request *balances.BpsCreateAssetRequest) *balances.BpsCreateAssetResponse {
	assetId, err := f.assetService.CreateNewAsset(ctx)

	if err != nil {
		return &balances.BpsCreateAssetResponse{Id: request.Id, Errors: map[string]*balances.BpsError{
			"CreationAsset": {Message: err.Error(), ErrorCode: utils.MapError(err)},
		}}
	}

	response := balances.BpsCreateAssetResponse{Id: request.Id, AssetId: assetId}

	for _, emmitInfo := range request.EmmitInfo {
		err := f.balanceService.EmmitBalance(ctx, assetId, emmitInfo.CurrencyName, emmitInfo.Amount)
		if err != nil {
			if response.Errors == nil {
				response.Errors = make(map[string]*balances.BpsError)
			}
			response.Errors[emmitInfo.CurrencyName] = &balances.BpsError{Message: err.Error(), ErrorCode: utils.MapError(err)}
		}
	}

	return &response

}
