package flow

import (
	"context"
	"time"
	"trade-balance-service/constants"
	"trade-balance-service/dto"
	"trade-balance-service/external/balances"
	"trade-balance-service/utils"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
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

type ISender interface {
	SendMessage(ctx context.Context, message protoreflect.ProtoMessage, exchange, rk string) error
}

type Flow struct {
	assetService   IAssetService
	balanceService IBalanceService
	sender         ISender
}

func NewFlow(assetService IAssetService, balanceService IBalanceService, sender ISender) *Flow {
	return &Flow{assetService: assetService, balanceService: balanceService, sender: sender}
}

func (f *Flow) CreateAsset(ctx context.Context, request *balances.BpsCreateAssetRequest) *balances.BpsCreateAssetResponse {
	assetId, err := f.assetService.CreateNewAsset(ctx)

	if err != nil {
		return &balances.BpsCreateAssetResponse{Id: request.Id, Error: &balances.BpsError{Message: err.Error(), ErrorCode: utils.MapError(err)}}
	}

	response := balances.BpsCreateAssetResponse{Id: request.Id, AssetId: assetId}

	go func() {
		for _, emmitInfo := range request.EmmitInfo {
			resp := balances.BpsEmmitAssetResponse{
				Id:           request.Id,
				AssetId:      assetId,
				CurrencyCode: emmitInfo.CurrencyName,
				Amount:       emmitInfo.Amount,
			}
			err := f.balanceService.EmmitBalance(ctx, assetId, emmitInfo.CurrencyName, emmitInfo.Amount)
			if err != nil {
				resp.Error = &balances.BpsError{Message: err.Error(), ErrorCode: utils.MapError(err)}
			}
			f.sender.SendMessage(ctx, &resp, constants.BpsExchange, constants.RkEmmitAssetResponse)
		}
	}()

	return &response

}

func (f *Flow) EmmitAsset(ctx context.Context, request *balances.BpsEmmitAssetRequest) {
	_, err := f.assetService.GetAssetInfoById(ctx, request.GetAssetId())

	response := balances.BpsEmmitAssetResponse{Id: request.Id, AssetId: request.AssetId}

	if err != nil {
		response.Error = &balances.BpsError{ErrorCode: utils.MapError(err), Message: err.Error()}
		f.sender.SendMessage(ctx, &response, constants.BpsExchange, constants.RkEmmitAssetResponse)
		return
	}

	for _, emmitData := range request.EmitBalancesInfo {
		resp := balances.BpsEmmitAssetResponse{
			Id:           request.Id,
			AssetId:      request.AssetId,
			CurrencyCode: emmitData.CurrencyName,
			Amount:       emmitData.Amount,
		}
		err := f.balanceService.EmmitBalance(ctx, request.AssetId, emmitData.CurrencyName, emmitData.Amount)
		if err != nil {
			resp.Error = &balances.BpsError{Message: err.Error(), ErrorCode: utils.MapError(err)}
		}
		f.sender.SendMessage(ctx, &resp, constants.BpsExchange, constants.RkEmmitAssetResponse)
	}
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
