package handler

import (
	"context"
	"trade-balance-service/constants"
	"trade-balance-service/external/balances"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type IFlow interface {
	CreateAsset(ctx context.Context, request *balances.BpsCreateAssetRequest) *balances.BpsCreateAssetResponse
	EmmitAsset(ctx context.Context, request *balances.EmmitBalanceRequest) *balances.EmmitBalanceResponse
	GetAssetsById(ctx context.Context, request *balances.BbsGetAssetInfoRequest) *balances.BpsGetAssetInfoResponse
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

func (h *HandlerCollection) HandleCreateAsset(ctx context.Context, request *balances.BpsCreateAssetRequest) {
	resp := h.flow.CreateAsset(ctx, request)
	h.sender.SendMessage(ctx, resp, constants.BpsExchange, constants.RkCreateAssetResponse)
}

func (h *HandlerCollection) HandleEmmitAsset(ctx context.Context, request *balances.EmmitBalanceRequest) {
	resp := h.flow.EmmitAsset(ctx, request)
	h.sender.SendMessage(ctx, resp, constants.BpsExchange, constants.RkEmmitAssetResponse)
}

func (h *HandlerCollection) HandleGetAssetsById(ctx context.Context, request *balances.BbsGetAssetInfoRequest) {
	resp := h.flow.GetAssetsById(ctx, request)
	h.sender.SendMessage(ctx, resp, constants.BpsExchange, constants.RkGetAssetsResponse)
}
