package services

import (
	"context"
	"trade-balance-service/dto"
	"trade-balance-service/staticserr"
)

type IBalanceProvider interface {
	GetInfoAboutBalanceByCurrency(ctx context.Context, assetId string, currencyId string) (*dto.BalanceModel, error)
	GetInfoAboutBalanceById(ctx context.Context, id string) (*dto.BalanceModel, error)
	GetBalancesInfoByAssetId(ctx context.Context, assetId string) ([]dto.BalanceModel, error)
	InsertBalanceInfo(ctx context.Context, assetId string, currencyId string, amount float64) error
	EmmitBalanceByCurrency(ctx context.Context, assetId string, currencyId string, amount float64) error
	EmmitBalanceById(ctx context.Context, id string, amount float64) error
	LockBalanceById(ctx context.Context, id string, amount float64) error
	RefundBalanceById(ctx context.Context, id string, amount float64) error
	ChargeBalanceById(ctx context.Context, id string, amount float64) error
	DeleteBalanceById(ctx context.Context, id string) error
}

type ICurrencyProvider interface {
	GetCurrencyInfoByCode(ctx context.Context, code string) (*dto.CurrencyModel, error)
	InsertCurrencyInfo(ctx context.Context, code string, fullName string) error
}

type BalanceService struct {
	balanceProvider  IBalanceProvider
	currencyProvider ICurrencyProvider
}

func NewBalanceService(balanceProvider IBalanceProvider, currencyProvider ICurrencyProvider) BalanceService {
	return BalanceService{balanceProvider: balanceProvider, currencyProvider: currencyProvider}
}

func (b BalanceService) EmmitBalance(ctx context.Context, assetId string, currencyCode string, amount float64) error {
	currencyModel, err := b.currencyProvider.GetCurrencyInfoByCode(ctx, currencyCode)

	if err != nil {
		return err
	}

	balanceModel, err := b.balanceProvider.GetInfoAboutBalanceByCurrency(ctx, assetId, currencyModel.Id)

	if err == staticserr.ErrorNotEnoughBalance {
		if err = b.balanceProvider.InsertBalanceInfo(ctx, assetId, currencyModel.Id, amount); err != nil {
			return err
		}
		return nil
	}
	if err != nil {
		return err
	}

	err = b.balanceProvider.EmmitBalanceById(ctx, balanceModel.Id, amount)
	if err != nil {
		return err
	}

	return nil
}

func (b BalanceService) AddCurrency(ctx context.Context, currencyName string, currencyCode string) error {
	_, err := b.currencyProvider.GetCurrencyInfoByCode(ctx, currencyCode)

	if err != nil && err != staticserr.ErrorNotExistsCurrency {
		return err
	}

	if err == nil {
		return staticserr.ErrorCurrencyAlreadyExists
	}

	err = b.currencyProvider.InsertCurrencyInfo(ctx, currencyCode, currencyName)
	if err != nil {
		return err
	}

	return nil
}

func (b BalanceService) GetInfoAboutAssets(ctx context.Context, assetId string) ([]dto.BalanceModel, error) {
	return b.balanceProvider.GetBalancesInfoByAssetId(ctx, assetId)
}
