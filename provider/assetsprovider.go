package provider

import (
	"context"
	"time"
	"trade-balance-service/dto"
	"trade-balance-service/staticserr"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

var (
	getAssetQuery    = "select id, account_id, creation_date, state from assets where id = $1"
	insertAssetQuery = "insert into assets values ($1, $2, $3, $4)"
	deleteAssetQuery = "update assets set state = $2 where id = $1"
)

type AssetsProvider struct {
	commonProvider *PgxProvider
}

func NewAssetsProvider(commonProvider *PgxProvider) AssetsProvider {
	return AssetsProvider{commonProvider: commonProvider}
}

func (a *AssetsProvider) GetAssetInfoById(ctx context.Context, id string) (*dto.TradeAsset, error) {
	row, err := a.commonProvider.ExecuteQueryWithRow(ctx, getAssetQuery, id)

	if err == pgx.ErrNoRows {
		return nil, staticserr.ErrorNotFoundAsset
	}

	if err != nil {
		return nil, err
	}

	result := dto.TradeAsset{}

	err = row.Scan(&result.Id, &result.AccountId, &result.CreatedDate, &result.State)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (a *AssetsProvider) InsertNewAssetInfo(ctx context.Context, accountId string) (string, error) {
	id := uuid.NewString()
	createdDate := time.Now().UTC().UnixMilli()

	err := a.commonProvider.ExecuteQuery(ctx, insertAssetQuery, id, createdDate, dto.ACTIVE, accountId)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (a *AssetsProvider) DeleteAssetById(ctx context.Context, id string) error {
	_, err := a.GetAssetInfoById(ctx, id)

	if err != nil {
		return err
	}

	err = a.commonProvider.ExecuteQuery(ctx, deleteAssetQuery, id, dto.DELETED)

	if err != nil {
		return err
	}

	return nil
}
