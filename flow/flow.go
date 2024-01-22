package flow

import (
	"context"
	"time"
	"trade-balance-service/constants"
	"trade-balance-service/dto"
	"trade-balance-service/external/bps"
	"trade-balance-service/staticserr"
	"trade-balance-service/utils"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IAssetService interface {
	CreateNewAsset(ctx context.Context, accountId string) (string, error)
	GetAssetInfoById(ctx context.Context, id string) (*dto.TradeAsset, error)
	DeactivateAsset(ctx context.Context, id string) error
}

type IBalanceService interface {
	EmmitBalance(ctx context.Context, assetId string, currencyCode string, amount float64) error
	AddCurrency(ctx context.Context, currencyName string, currencyCode string) error
	GetInfoAboutAssets(ctx context.Context, assetId string) ([]dto.PublicBalanceModel, error)
	LockBalance(ctx context.Context, currencyCode string, assetId string, amount float64) (string, error)
	RefundBalance(ctx context.Context, id string, amount float64) error
	CreateTransfer(ctx context.Context, request *bps.BpsCreateTransferRequest) error
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

func (f *Flow) CreateAsset(ctx context.Context, request *bps.BpsCreateAssetRequest) *bps.BpsCreateAssetResponse {
	if request.AccountId == "" {
		return &bps.BpsCreateAssetResponse{Id: request.Id, Error: &bps.BpsError{Message: staticserr.ErrorNotRelatedAccount.Error(), ErrorCode: utils.MapError(staticserr.ErrorNotRelatedAccount)}}
	}

	assetId, err := f.assetService.CreateNewAsset(ctx, request.AccountId)

	if err != nil {
		return &bps.BpsCreateAssetResponse{Id: request.Id, Error: &bps.BpsError{Message: err.Error(), ErrorCode: utils.MapError(err)}}
	}

	response := bps.BpsCreateAssetResponse{Id: request.Id, AssetId: assetId}

	go func() {
		for _, emmitInfo := range request.EmmitInfo {
			resp := bps.BpsEmmitAssetResponse{
				Id:           request.Id,
				AccountId:    request.AccountId,
				AssetId:      assetId,
				CurrencyCode: emmitInfo.CurrencyName,
				Amount:       emmitInfo.Amount,
			}
			err := f.balanceService.EmmitBalance(ctx, assetId, emmitInfo.CurrencyName, emmitInfo.Amount)
			if err != nil {
				resp.Error = &bps.BpsError{Message: err.Error(), ErrorCode: utils.MapError(err)}
			}
			f.sender.SendMessage(ctx, &resp, constants.BpsExchange, constants.RkEmmitAssetResponse)
		}
	}()

	return &response

}

func (f *Flow) EmmitAsset(ctx context.Context, request *bps.BpsEmmitAssetRequest) {
	logrus.Infoln("Received request for emmit: ", request.String())
	assetInfo, err := f.assetService.GetAssetInfoById(ctx, request.GetAssetId())

	response := bps.BpsEmmitAssetResponse{Id: request.Id, AssetId: request.AssetId}

	if err != nil {
		response.Error = &bps.BpsError{ErrorCode: utils.MapError(err), Message: err.Error()}
		f.sender.SendMessage(ctx, &response, constants.BpsExchange, constants.RkEmmitAssetResponse)
		return
	}

	for _, emmitData := range request.EmitBalancesInfo {
		resp := bps.BpsEmmitAssetResponse{
			Id:           request.Id,
			AccountId:    assetInfo.AccountId,
			AssetId:      request.AssetId,
			CurrencyCode: emmitData.CurrencyName,
			Amount:       emmitData.Amount,
		}
		err := f.balanceService.EmmitBalance(ctx, request.AssetId, emmitData.CurrencyName, emmitData.Amount)
		if err != nil {
			resp.Error = &bps.BpsError{Message: err.Error(), ErrorCode: utils.MapError(err)}
		}
		f.sender.SendMessage(ctx, &resp, constants.BpsExchange, constants.RkEmmitAssetResponse)
	}
}

func (f *Flow) GetAssetsById(ctx context.Context, request *bps.BbsGetAssetInfoRequest) *bps.BpsGetAssetInfoResponse {
	assetInfo, err := f.assetService.GetAssetInfoById(ctx, request.AssetId)

	if err != nil {
		return &bps.BpsGetAssetInfoResponse{
			Id:    request.GetId(),
			Error: &bps.BpsError{Message: err.Error(), ErrorCode: utils.MapError(err)}}
	}

	if assetInfo.AccountId != request.AccountId {
		return &bps.BpsGetAssetInfoResponse{
			Id:    request.Id,
			Error: &bps.BpsError{Message: staticserr.ErrorNotRelatedAccount.Error(), ErrorCode: utils.MapError(staticserr.ErrorNotRelatedAccount)},
		}
	}

	assetBalancesInfo, err := f.balanceService.GetInfoAboutAssets(ctx, assetInfo.Id)

	resp := mapAssetInfoToProto(*assetInfo, assetBalancesInfo)
	resp.Id = request.Id
	return resp
}

func (f *Flow) DeactivateAsset(ctx context.Context, request *bps.BpsDeactivateAssetRequest) *bps.BpsDeactivateAssetResponse {
	assetInfo, err := f.assetService.GetAssetInfoById(ctx, request.AssetId)

	if err != nil {
		return &bps.BpsDeactivateAssetResponse{
			Id:    request.Id,
			Error: &bps.BpsError{Message: err.Error(), ErrorCode: utils.MapError(err)},
		}
	}

	if assetInfo.AccountId != request.AccountId {
		return &bps.BpsDeactivateAssetResponse{
			Id:    request.Id,
			Error: &bps.BpsError{Message: staticserr.ErrorNotRelatedAccount.Error(), ErrorCode: bps.BpsErrorCode_BPS_ERROR_CODE_ASSET_NOT_RELATED_TO_ACCOUNT},
		}
	}

	err = f.assetService.DeactivateAsset(ctx, request.AssetId)

	if err != nil {
		return &bps.BpsDeactivateAssetResponse{
			Id:    request.Id,
			Error: &bps.BpsError{Message: err.Error(), ErrorCode: utils.MapError(err)},
		}
	}

	return &bps.BpsDeactivateAssetResponse{Id: request.Id}
}

func (f *Flow) AddNewCurrency(ctx context.Context, request *bps.BpsAddCurrencyRequest) *bps.BpsAddCurrencyResponse {
	response := &bps.BpsAddCurrencyResponse{Id: request.Id}

	err := f.balanceService.AddCurrency(ctx, request.CurrencyName, request.CurrencyCode)

	if err != nil {
		response.Error = &bps.BpsError{Message: err.Error(), ErrorCode: utils.MapError(err)}
	}

	return response
}

func (f *Flow) LockBalanceAsset(ctx context.Context, request *bps.BpsLockBalanceRequest) *bps.BpsLockBalanceResponse {
	response := &bps.BpsLockBalanceResponse{Id: request.Id}

	assetInfo, err := f.assetService.GetAssetInfoById(ctx, request.AssetId)

	if err != nil {
		response.Error = &bps.BpsError{Message: err.Error(), ErrorCode: utils.MapError(err)}
		return response
	}

	if assetInfo.AccountId != request.AccountId {
		response.Error = &bps.BpsError{Message: staticserr.ErrorNotRelatedAccount.Error(), ErrorCode: bps.BpsErrorCode_BPS_ERROR_CODE_ASSET_NOT_RELATED_TO_ACCOUNT}
		return response
	}

	id, err := f.balanceService.LockBalance(ctx, request.CurrencyCode, request.AssetId, request.Amount)

	if err != nil {
		response.Error = &bps.BpsError{Message: err.Error(), ErrorCode: utils.MapError(err)}
		return response
	}

	response.BalanceId = id
	return response
}

func (f *Flow) RefundBalanceAsset(ctx context.Context, request *bps.BpsRefundBalanceRequest) *bps.BpsRefundBalanceResponse {
	response := &bps.BpsRefundBalanceResponse{Id: request.Id}

	err := f.balanceService.RefundBalance(ctx, request.BalanceId, request.Amount)

	if err != nil {
		response.Error = &bps.BpsError{Message: err.Error(), ErrorCode: utils.MapError(err)}
	}

	return response
}

func (f *Flow) CreateTransfer(ctx context.Context, request *bps.BpsCreateTransferRequest) *bps.BpsTransfer {

	response := &bps.BpsTransfer{
		Id:           request.Id,
		TransferData: request.TransferData,
	}

	err := f.balanceService.CreateTransfer(ctx, request)

	if err == staticserr.ErrorNotEnoughBalance {
		response.TransferState = bps.BpsTransferState_BPS_TRANSFER_STATE_REJECTED
		return response
	}

	if err != nil {
		response.TransferState = bps.BpsTransferState_BPS_TRANSFER_STATE_ERROR
		response.Error = &bps.BpsError{Message: err.Error(), ErrorCode: utils.MapError(err)}
		return response
	}

	response.TransferState = bps.BpsTransferState_BPS_TRANSFER_STATE_DONE

	return response

}

func mapAssetInfoToProto(asset dto.TradeAsset, balancesInfo []dto.PublicBalanceModel) *bps.BpsGetAssetInfoResponse {
	protoModel := bps.BpsGetAssetInfoResponse{
		AssetId:     asset.Id,
		CreatedDate: timestamppb.New(time.UnixMilli(asset.CreatedDate)),
	}

	for _, balance := range balancesInfo {
		protoModel.BalancesInfo = append(protoModel.BalancesInfo, &bps.BalanceInfo{
			CurrencyName: balance.CurrencyName,
			Amount:       balance.Amount,
			LockedAmount: balance.LockedAmount,
		})
	}

	return &protoModel
}
