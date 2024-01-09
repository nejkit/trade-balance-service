package provider

import (
	"context"
	"github.com/jackc/pgx/v4"
	"trade-balance-service/dto"
	"trade-balance-service/staticserr"
)

var (
	getBalanceByCurrencyQuery  = "select id, asset_id, currency_id, amount, locked_amount from balances where asset_id = $1 and currency_id = $2"
	getCurrencyInfoByShortName = "select id, currency_name, currency_code from currencies where currency_code = $1"
	getBalanceByIdQuery        = "select id, asset_id, currency_id, amount, locked_amount from balances where if = $1"
	emitBalanceByIdQuery       = "update balances set amount = amount + $2 where id = $1"
	chargeBalanceByIdQuery     = "update balances set locked_amount = locked_amount - $2 where id = $1"
	refundBalanceByIdQuery     = "update balances set amount = amount + $2, locked_amount = locked_amount - $2 where id = $1"
	lockBalanceByIdQuery       = "update balances set amount = amount - $2, locked_amount = locked_amount + $2 where id = $1"
	insertBalanceQuery         = "insert into balances ($1, $2, $3, $4, $5)"
	deleteBalanceQuery         = "delete balances where id = $1"
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

func (b *BalancesProvider) InsertBalanceInfo(ctx context.Context, assetId string, currencyId string, amount float64) {

}

func parseResponse(row pgx.Row, err error) (*dto.BalanceModel, error) {
	if err == pgx.ErrNoRows {
		return nil, staticserr.ErrorNotEnoughBalance
	}

	if err != nil {
		return nil, err
	}

	result := dto.BalanceModel{}

	err = row.Scan(&result.Id, &result.AssetId, &result.CurrencyId, &result.Amount, &result.LockedAmount)

	if err != nil {
		return nil, err
	}

	return &result, nil
}
