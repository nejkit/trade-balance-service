package provider

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"trade-balance-service/dto"
	"trade-balance-service/staticserr"
)

var (
	getCurrencyInfoByShortName = "select id, currency_name, currency_code from currencies where currency_code = $1"
	addNewCurrencyInfoQuery    = "insert into currencies ($1, $2, $3)"
)

type CurrencyProvider struct {
	commonProvider *PgxProvider
}

func NewCurrencyProvider(commonProvider *PgxProvider) CurrencyProvider {
	return CurrencyProvider{commonProvider: commonProvider}
}

func (c CurrencyProvider) GetCurrencyInfoByCode(ctx context.Context, code string) (*dto.CurrencyModel, error) {
	row, err := c.commonProvider.ExecuteQueryWithRow(ctx, getCurrencyInfoByShortName, code)

	if err == pgx.ErrNoRows {
		return nil, staticserr.ErrorNotExistsCurrency
	}
	if err != nil {
		return nil, err
	}

	data := dto.CurrencyModel{}

	if err = row.Scan(&data.Id, &data.CurrencyName, &data.CurrencyCode); err != nil {
		return nil, err
	}

	return &data, nil
}

func (c CurrencyProvider) InsertCurrencyInfo(ctx context.Context, code string, fullName string) error {
	id := uuid.NewString()
	if err := c.commonProvider.ExecuteQuery(ctx, addNewCurrencyInfoQuery, id, fullName, code); err != nil {
		return err
	}
	return nil
}
