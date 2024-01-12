package provider

import (
	"context"
	"trade-balance-service/dto"
	"trade-balance-service/staticserr"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

var (
	getBalanceByCurrencyQuery   = "select id, asset_id, currency_id, amount, locked_amount from balances where asset_id = $1 and currency_id = $2"
	getBalanceByIdQuery         = "select id, asset_id, currency_id, amount, locked_amount from balances where id = $1"
	getBalancesByAssetIdQuery   = "select id, asset_id, currency_id, amount, locked_amount from balances where asset_id = $1"
	emmitBalanceByIdQuery       = "update balances set amount = amount + $2 where id = $1"
	emmitBalanceByCurrencyQuery = "update balances set amount = amount + $3 where asset_id = $1 and currency_id = $2"
	chargeBalanceByIdQuery      = "update balances set locked_amount = locked_amount - $2 where id = $1"
	refundBalanceByIdQuery      = "update balances set amount = amount + $2, locked_amount = locked_amount - $2 where id = $1"
	lockBalanceByIdQuery        = "update balances set amount = amount - $2, locked_amount = locked_amount + $2 where id = $1"
	insertBalanceQuery          = "insert into balances values ($1, $2, $3, $4, $5)"
	deleteBalanceQuery          = "delete balances where id = $1"
)

type BalancesProvider struct {
	commonProvider *PgxProvider
}

func NewBalancesProvider(commonProvider *PgxProvider) BalancesProvider {
	return BalancesProvider{commonProvider: commonProvider}
}

func (b *BalancesProvider) GetInfoAboutBalanceByCurrency(ctx context.Context, assetId string, currencyId string) (*dto.BalanceModel, error) {
	row, err := b.commonProvider.ExecuteQueryWithRow(ctx, getBalanceByCurrencyQuery, assetId, currencyId)

	return parseResponse(row, err)
}

func (b *BalancesProvider) GetInfoAboutBalanceById(ctx context.Context, id string) (*dto.BalanceModel, error) {
	row, err := b.commonProvider.ExecuteQueryWithRow(ctx, getBalanceByIdQuery, id)

	return parseResponse(row, err)
}

func (b *BalancesProvider) InsertBalanceInfo(ctx context.Context, assetId string, currencyId string, amount float64) error {
	id := uuid.NewString()
	if err := b.commonProvider.ExecuteQuery(ctx, insertBalanceQuery, id, assetId, currencyId, amount, 0); err != nil {
		return err
	}
	return nil
}

func (b *BalancesProvider) EmmitBalanceByCurrency(ctx context.Context, assetId string, currencyId string, amount float64) error {
	if err := b.commonProvider.ExecuteQuery(ctx, emmitBalanceByCurrencyQuery, assetId, currencyId, amount); err != nil {
		return err
	}
	return nil
}

func (b *BalancesProvider) EmmitBalanceById(ctx context.Context, id string, amount float64) error {
	if err := b.commonProvider.ExecuteQuery(ctx, emmitBalanceByIdQuery, id, amount); err != nil {
		return err
	}
	return nil
}

func (b *BalancesProvider) LockBalanceById(ctx context.Context, id string, amount float64) error {
	if err := b.commonProvider.ExecuteQuery(ctx, lockBalanceByIdQuery, id, amount); err != nil {
		return err
	}
	return nil
}

func (b *BalancesProvider) RefundBalanceById(ctx context.Context, id string, amount float64) error {
	if err := b.commonProvider.ExecuteQuery(ctx, refundBalanceByIdQuery, id, amount); err != nil {
		return err
	}

	return nil
}

func (b *BalancesProvider) ChargeBalanceById(ctx context.Context, id string, amount float64) error {
	if err := b.commonProvider.ExecuteQuery(ctx, chargeBalanceByIdQuery, id, amount); err != nil {
		return err
	}

	return nil
}

func (b *BalancesProvider) DeleteBalanceById(ctx context.Context, id string) error {
	if err := b.commonProvider.ExecuteQuery(ctx, deleteBalanceQuery, id); err != nil {
		return err
	}

	return nil
}

func (b *BalancesProvider) GetBalancesInfoByAssetId(ctx context.Context, assetId string) ([]dto.BalanceModel, error) {
	rows, err := b.commonProvider.ExecuteQueryRows(ctx, getBalancesByAssetIdQuery, assetId)

	if err == pgx.ErrNoRows {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	var result []dto.BalanceModel
	for rows.Next() {
		data := dto.BalanceModel{}
		if err = rows.Scan(&data.Id, &data.AssetId, &data.CurrencyId, &data.Amount, &data.LockedAmount); err != nil {
			continue
		}
		result = append(result, data)
	}
	return result, nil
}

func parseResponse(row pgx.Row, err error) (*dto.BalanceModel, error) {

	if err != nil {
		return nil, err
	}

	result := dto.BalanceModel{}

	err = row.Scan(&result.Id, &result.AssetId, &result.CurrencyId, &result.Amount, &result.LockedAmount)

	if err == pgx.ErrNoRows {
		return nil, staticserr.ErrorNotEnoughBalance
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}
