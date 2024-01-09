package provider

import "github.com/jackc/pgx"

type AssetsProvider struct {
	connectionPool *pgx.ConnPool
}

var (
	getAssetQuery = "select id, creation_date, state from assets where id = $1"
	getBalanceByCurrencyQuery = "select id, asset_id, currency_id, amount, locked_amount from balances where asset_id = $1 and currency_id = $2"
	getCurrencyInfoByShortName = "select id, currency_name, currency_code from currencies where currency_code = $1"
	getBalanceByIdQuery = "select id, asset_id, currency_id, amount, locked_amount from balances where if = $1"
	emitBalanceByIdQuery = "update balances set amount = amount + $2 where id = $1"
	chargeBalanceByIdQuery = "update balances set locked_amount = locked_amount - $2 where id = $1"
	refundBalanceByIdQuery = "update balances set amount = amount + $2, locked_amount = locked_amount - $2 where id = $1"
	lockBalanceByIdQuery = "update balances set amount = amount - $2, locked_amount = locked_amount + $2 where id = $1"
	insertBalanceQuery = "insert into balances ($1, $2, $3, $4, $5)"
	deleteBalanceQuery = "delete balances where id = $1"
)

func NewAssetsProvider(connectionPool *pgx.ConnPool) AssetsProvider {
	return AssetsProvider{connectionPool: connectionPool}
}

func (p *AssetsProvider) ()  {

}