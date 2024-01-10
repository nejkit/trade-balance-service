package flow

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
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
	GetInfoAboutAssets(ctx context.Context, assetId string) ([]dto.PublicBalanceModel, error)
}

type Flow struct {
	assetService   IAssetService
	balanceService IBalanceService
}

func NewFlow(assetService IAssetService, balanceService IBalanceService) *Flow {
	return &Flow{assetService: assetService, balanceService: balanceService}
}

func (f *Flow) CreateAsset(ctx context.Context, request *balances.BpsCreateAssetRequest) *balances.BpsCreateAssetResponse {
	assetId, err := f.assetService.CreateNewAsset(ctx)

	if err != nil {
		return &balances.BpsCreateAssetResponse{Id: request.Id, Errors: map[string]*balances.BpsError{
			"asset": {Message: err.Error(), ErrorCode: utils.MapError(err)},
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

func (f *Flow) EmmitAsset(ctx context.Context, request *balances.EmmitBalanceRequest) *balances.EmmitBalanceResponse {
	_, err := f.assetService.GetAssetInfoById(ctx, request.GetAssetId())

	response := balances.EmmitBalanceResponse{Id: request.Id, AssetId: request.AssetId}

	if err != nil {
		return &balances.EmmitBalanceResponse{Id: request.Id, AssetId: request.AssetId, Errors: map[string]*balances.BpsError{
			"asset": {Message: err.Error(), ErrorCode: utils.MapError(err)},
		}}
	}

	for _, emmitData := range request.EmitBalancesInfo {
		err := f.balanceService.EmmitBalance(ctx, request.AssetId, emmitData.CurrencyName, emmitData.Amount)
		if err != nil {
			if response.Errors == nil {
				response.Errors = make(map[string]*balances.BpsError)
			}
			response.Errors[emmitData.CurrencyName] = &balances.BpsError{Message: err.Error(), ErrorCode: utils.MapError(err)}
		}
	}

	return &response
}

func (f *Flow) GetAssetsById(ctx context.Context, request *balances.BbsGetAssetInfoRequest) *balances.BpsGetAssetInfoResponse {
	assetInfo, err := f.assetService.GetAssetInfoById(ctx, request.AssetId)

	if err != nil {
		return &balances.BpsGetAssetInfoResponse{
			Id:    request.GetId(),
			Error: &balances.BpsError{Message: err.Error(), ErrorCode: utils.MapError(err)}}
	}

	assetBalancesInfo, err := f.balanceService.GetInfoAboutAssets(ctx, assetInfo.Id)

	resp := mapAssetInfoToProto(*assetInfo, assetBalancesInfo)
	resp.Id = request.Id
	return resp
}

func mapAssetInfoToProto(asset dto.TradeAsset, balancesInfo []dto.PublicBalanceModel) *balances.BpsGetAssetInfoResponse {
	protoModel := balances.BpsGetAssetInfoResponse{
		AssetId:     asset.Id,
		CreatedDate: timestamppb.New(time.UnixMilli(asset.CreatedDate)),
	}

	for _, balance := range balancesInfo {
		protoModel.BalancesInfo = append(protoModel.BalancesInfo, &balances.BalanceInfo{
			CurrencyName: balance.CurrencyName,
			Amount:       balance.Amount,
			LockedAmount: balance.LockedAmount,
		})
	}

	return &protoModel
}
