package dto

type CurrencyModel struct {
	Id           string
	CurrencyName string
	CurrencyCode string
}

type BalanceModel struct {
	Id           string
	AssetId      string
	CurrencyId   string
	Amount       float64
	LockedAmount float64
}

type PublicBalanceModel struct {
	CurrencyName string
	Amount       float64
	LockedAmount float64
}
