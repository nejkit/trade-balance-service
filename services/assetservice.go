package services

import (
	"context"
	"trade-balance-service/dto"
	"trade-balance-service/staticserr"
)

type IAssetProvider interface {
	GetAssetInfoById(ctx context.Context, id string) (*dto.TradeAsset, error)
	InsertNewAssetInfo(ctx context.Context, accountId string) (string, error)
	DeleteAssetById(ctx context.Context, id string) error
}

type AssetService struct {
	provider IAssetProvider
}

func NewAssetService(provider IAssetProvider) AssetService {
	return AssetService{provider: provider}
}

func (a *AssetService) CreateNewAsset(ctx context.Context, accountId string) (string, error) {
	return a.provider.InsertNewAssetInfo(ctx, accountId)
}

func (a *AssetService) GetAssetInfoById(ctx context.Context, id string) (*dto.TradeAsset, error) {
	info, err := a.provider.GetAssetInfoById(ctx, id)

	if err != nil {
		return nil, err
	}

	if info.State == dto.ACTIVE {
		return info, nil
	}

	return nil, staticserr.ErrorNotFoundAsset
}

func (a *AssetService) DeactivateAsset(ctx context.Context, id string) error {
	return a.provider.DeleteAssetById(ctx, id)
}
