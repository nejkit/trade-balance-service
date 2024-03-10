package provider

import (
	"context"
	"time"
	"trade-balance-service/dto"
	"trade-balance-service/external/bps"
	"trade-balance-service/staticserr"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

var (
	getBalanceByCurrencyQuery      = "select id, asset_id, currency_id, amount, locked_amount from balances where asset_id = $1 and currency_id = $2 for update"
	getBalanceByCurrencyQueryNotTx = "select id, asset_id, currency_id, amount, locked_amount from balances where asset_id = $1 and currency_id = $2"
	getBalanceByIdQuery            = "select id, asset_id, currency_id, amount, locked_amount from balances where id = $1 for update"
	getBalancesByAssetIdQuery      = "select id, asset_id, currency_id, amount, locked_amount from balances where asset_id = $1"
	emmitBalanceByIdQuery          = "update balances set amount = amount + $2 where id = $1"
	emmitBalanceByCurrencyQuery    = "update balances set amount = amount + $3 where asset_id = $1 and currency_id = $2"
	chargeBalanceByIdQuery         = "update balances set locked_amount = locked_amount - $2 where id = $1"
	refundBalanceByIdQuery         = "update balances set amount = amount + $2, locked_amount = locked_amount - $2 where id = $1"
	lockBalanceByIdQuery           = "update balances set amount = amount - $2, locked_amount = locked_amount + $2 where id = $1"
	insertBalanceQuery             = "insert into balances values ($1, $2, $3, $4, $5)"
	deleteBalanceQuery             = "delete balances where id = $1"
)

type BalancesProvider struct {
	commonProvider *PgxProvider
}

func NewBalancesProvider(commonProvider *PgxProvider) BalancesProvider {
	return BalancesProvider{commonProvider: commonProvider}
}

func (b *BalancesProvider) GetInfoAboutBalanceByCurrency(ctx context.Context, assetId string, currencyId string) (*dto.BalanceModel, error) {
	row, err := b.commonProvider.ExecuteQueryWithRow(ctx, getBalanceByCurrencyQueryNotTx, assetId, currencyId)

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
	tx, err := b.commonProvider.PerformTx(ctx)

	if err != nil {
		return err
	}

	row := tx.ExecuteQueryWithRow(ctx, getBalanceByCurrencyQuery, assetId, currencyId)

	if row.Scan() != nil {
		tx.tx.Rollback(ctx)
		tx.tx.Conn().Close(ctx)
		return err
	}

	err = tx.ExecuteQuery(ctx, emmitBalanceByCurrencyQuery, assetId, currencyId, amount)

	if err != nil {
		return err
	}

	err = tx.CommitTx(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (b *BalancesProvider) EmmitBalanceById(ctx context.Context, id string, amount float64) error {

	tx, err := b.commonProvider.PerformTx(ctx)

	defer tx.tx.Conn().Close(ctx)

	if err != nil {
		return err
	}

	row := tx.ExecuteQueryWithRow(ctx, getBalanceByIdQuery, id)

	var balanceModel dto.BalanceModel

	if row.Scan(&balanceModel.Id, &balanceModel.AssetId, &balanceModel.CurrencyId, &balanceModel.Amount, &balanceModel.LockedAmount) != nil {
		logrus.Errorln("Error while scan: ", row.Scan().Error())
		tx.tx.Rollback(ctx)
		tx.tx.Conn().Close(ctx)
		return err
	}

	err = tx.ExecuteQuery(ctx, emmitBalanceByIdQuery, id, amount)

	if err != nil {
		return err
	}

	err = tx.CommitTx(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (b *BalancesProvider) LockBalanceByCurrency(ctx context.Context, assetId, currencyId string, amount float64) (string, error) {

	tx, err := b.commonProvider.PerformTx(ctx)

	defer tx.tx.Conn().Close(ctx)

	if err != nil {
		return "", err
	}

	row := tx.ExecuteQueryWithRow(ctx, getBalanceByCurrencyQuery, assetId, currencyId)

	var balanceModel dto.BalanceModel

	err = row.Scan(&balanceModel.Id, &balanceModel.AssetId, &balanceModel.CurrencyId, &balanceModel.Amount, &balanceModel.LockedAmount)

	if err == pgx.ErrNoRows {
		tx.tx.Rollback(ctx)
		return "", staticserr.ErrorNotEnoughBalance
	}

	if err != nil {
		tx.tx.Rollback(ctx)
		return "", err
	}

	if balanceModel.Amount < amount {
		tx.tx.Rollback(ctx)
		return "", staticserr.ErrorNotEnoughBalance
	}

	err = tx.ExecuteQuery(ctx, lockBalanceByIdQuery, balanceModel.Id, amount)

	if err != nil {
		return "", err
	}

	err = tx.CommitTx(ctx)

	if err != nil {
		return "", err
	}

	return balanceModel.Id, nil
}

func (b *BalancesProvider) RefundBalanceById(ctx context.Context, id string, amount float64) error {

	tx, err := b.commonProvider.PerformTx(ctx)

	defer tx.tx.Conn().Close(ctx)

	if err != nil {
		return err
	}

	row := tx.ExecuteQueryWithRow(ctx, getBalanceByIdQuery, id)

	var balanceModel dto.BalanceModel

	err = row.Scan(&balanceModel.Id, &balanceModel.AssetId, &balanceModel.CurrencyId, &balanceModel.Amount, &balanceModel.LockedAmount)
	if err == pgx.ErrNoRows {
		tx.tx.Rollback(ctx)
		return staticserr.ErrorNotEnoughBalance
	}

	if err != nil {
		tx.tx.Rollback(ctx)
		return err
	}

	err = tx.ExecuteQuery(ctx, refundBalanceByIdQuery, id, amount)

	if err != nil {
		return err
	}

	err = tx.CommitTx(ctx)

	if err != nil {
		return err
	}

	return nil
}

func chargeBalance(ctx context.Context, model dto.BalanceModel, amount float64, tx *txContainer) error {

	return tx.ExecuteQuery(ctx, chargeBalanceByIdQuery, model.Id, amount)
}

func emmitBalance(ctx context.Context, model dto.BalanceModel, amount float64, tx *txContainer) error {

	return tx.ExecuteQuery(ctx, emmitBalanceByIdQuery, model.Id, amount)
}

func (b *BalancesProvider) ChargeBalancesByIds(
	ctx context.Context,
	matchingInfos []*bps.BpsTransferData,
	extRespChan chan dto.TransferState) {

	childCtx, cancel := context.WithTimeout(ctx, time.Second*10)

	defer time.AfterFunc(time.Second, func() {
		cancel()
		close(extRespChan)
	})

	extRespChan <- dto.TransferState{State: bps.BpsTransferState_BPS_TRANSFER_STATE_NEW}
	tx, err := b.commonProvider.PerformTx(childCtx)
	if err != nil {
		extRespChan <- dto.TransferState{State: bps.BpsTransferState_BPS_TRANSFER_STATE_ERROR, Err: err}
		return
	}

	balances := make([]dto.BalanceModel, 2)
	responses := make([]chan pgx.Row, 2)
	chargeData := make(map[string]float64)
	transferData := make(map[string]float64)

	extRespChan <- dto.TransferState{State: bps.BpsTransferState_BPS_TRANSFER_STATE_IN_PROCESS}

	for _, tData := range matchingInfos {
		respChan := make(chan pgx.Row)
		tData := tData
		go func() {
			data := tx.ExecuteQueryWithRow(childCtx, getBalanceByIdQuery, tData.From)
			respChan <- data
		}()
		chargeData[tData.From] = tData.Amount
		transferData[tData.To] = tData.Amount
		responses = append(responses, respChan)
	}

	for _, resp := range responses {

		var balInfo dto.BalanceModel

		data := <-resp

		if err = data.Scan(&balInfo.Id, &balInfo.AssetId, &balInfo.CurrencyId, &balInfo.Amount, &balInfo.LockedAmount); err != nil {
			tx.tx.Rollback(childCtx)
			extRespChan <- dto.TransferState{State: bps.BpsTransferState_BPS_TRANSFER_STATE_ERROR, Err: err}
			return
		}

		balances = append(balances, balInfo)
	}

	for _, info := range balances {

		if info.LockedAmount < chargeData[info.Id] {
			tx.tx.Rollback(childCtx)
			extRespChan <- dto.TransferState{State: bps.BpsTransferState_BPS_TRANSFER_STATE_REJECTED}
			return
		}

		if err = chargeBalance(childCtx, info, chargeData[info.Id], tx); err != nil {
			tx.tx.Rollback(childCtx)
			extRespChan <- dto.TransferState{State: bps.BpsTransferState_BPS_TRANSFER_STATE_ERROR, Err: err}
			return
		}
	}

	extRespChan <- dto.TransferState{State: bps.BpsTransferState_BPS_TRANSFER_STATE_CHARGED}

	for _, info := range balances {
		if err = emmitBalance(childCtx, info, transferData[info.Id], tx); err != nil {
			tx.tx.Rollback(childCtx)
			extRespChan <- dto.TransferState{State: bps.BpsTransferState_BPS_TRANSFER_STATE_ERROR, Err: err}
			return
		}
	}

	if err = tx.CommitTx(childCtx); err != nil {
		tx.tx.Rollback(childCtx)
		extRespChan <- dto.TransferState{State: bps.BpsTransferState_BPS_TRANSFER_STATE_ERROR, Err: err}
		return
	}

	extRespChan <- dto.TransferState{State: bps.BpsTransferState_BPS_TRANSFER_STATE_DONE}
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
