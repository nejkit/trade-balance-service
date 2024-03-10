package handler

import (
	"context"
	"trade-balance-service/constants"
	"trade-balance-service/external/bps"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type IFlow interface {
	CreateAsset(ctx context.Context, request *bps.BpsCreateAssetRequest) *bps.BpsCreateAssetResponse
	EmmitAsset(ctx context.Context, request *bps.BpsEmmitAssetRequest)
	GetAssetsById(ctx context.Context, request *bps.BbsGetAssetInfoRequest) *bps.BpsGetAssetInfoResponse
	AddNewCurrency(ctx context.Context, request *bps.BpsAddCurrencyRequest) *bps.BpsAddCurrencyResponse
	DeactivateAsset(ctx context.Context, request *bps.BpsDeactivateAssetRequest) *bps.BpsDeactivateAssetResponse
	LockBalanceAsset(ctx context.Context, request *bps.BpsLockBalanceRequest) *bps.BpsLockBalanceResponse
	RefundBalanceAsset(ctx context.Context, request *bps.BpsRefundBalanceRequest) *bps.BpsRefundBalanceResponse
	CreateTransfer(ctx context.Context, request *bps.BpsCreateTransferRequest)
}

type ISender interface {
	SendMessage(ctx context.Context, message protoreflect.ProtoMessage, exchange, rk string) error
}

type HandlerCollection struct {
	flow   IFlow
	sender ISender
}

func NewHandlerCollection(flow IFlow, sender ISender) HandlerCollection {
	return HandlerCollection{flow: flow, sender: sender}
}

func (h *HandlerCollection) HandleCreateAsset(ctx context.Context, request *bps.BpsCreateAssetRequest) {
	resp := h.flow.CreateAsset(ctx, request)
	h.sender.SendMessage(ctx, resp, constants.BpsExchange, constants.RkCreateAssetResponse)
}

func (h *HandlerCollection) HandleEmmitAsset(ctx context.Context, request *bps.BpsEmmitAssetRequest) {
	h.flow.EmmitAsset(ctx, request)
}

func (h *HandlerCollection) HandleGetAssetsById(ctx context.Context, request *bps.BbsGetAssetInfoRequest) {
	resp := h.flow.GetAssetsById(ctx, request)
	logrus.Infoln("Asset Info for http: ", resp.String())
	h.sender.SendMessage(ctx, resp, constants.BpsExchange, constants.RkGetAssetsResponse)
}

func (h *HandlerCollection) HandleAddCurrency(ctx context.Context, request *bps.BpsAddCurrencyRequest) {
	resp := h.flow.AddNewCurrency(ctx, request)
	h.sender.SendMessage(ctx, resp, constants.BpsExchange, constants.RkAddCurrencyResponse)
}

func (h *HandlerCollection) HandleDeactivateAsset(ctx context.Context, request *bps.BpsDeactivateAssetRequest) {
	resp := h.flow.DeactivateAsset(ctx, request)
	h.sender.SendMessage(ctx, resp, constants.BpsExchange, constants.RkDeactivateAssetResponse)
}

func (h *HandlerCollection) HandleLockBalanceAsset(ctx context.Context, request *bps.BpsLockBalanceRequest) {
	resp := h.flow.LockBalanceAsset(ctx, request)
	h.sender.SendMessage(ctx, resp, constants.BpsExchange, constants.RkLockBalanceAssetResponse)
}

func (h *HandlerCollection) HandleRefundBalanceAsset(ctx context.Context, request *bps.BpsRefundBalanceRequest) {
	resp := h.flow.RefundBalanceAsset(ctx, request)
	h.sender.SendMessage(ctx, resp, constants.BpsExchange, constants.RkUnlockBalanceAssetResponse)
}

func (h *HandlerCollection) HandleCreateTransfer(ctx context.Context, request *bps.BpsCreateTransferRequest) {
	go h.flow.CreateTransfer(ctx, request)
}
