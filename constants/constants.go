package constants

const (
	BpsExchange = "bps.forward"
)

const (
	CreateAssetQueueName        = "q.bps.request.CreateAsset"
	EmmitAssetQueueName         = "q.bps.request.EmmitAsset"
	GetAssetsByIdQueueName      = "q.bps.request.GetAssets"
	AddNewCurrencyQueueName     = "q.bps.request.AddCurrency"
	DeactivateAssetQueueName    = "q.bps.request.DeactivateAsset"
	LockBalanceAssetQueueName   = "q.bps.request.LockBalanceAsset"
	UnlockBalanceAssetQueueName = "q.bps.request.UnlockBalanceAsset"
	CreateTransferQueueName     = "q.bps.request.CreateTransfer"
)

const (
	RkCreateAssetRequest  = "r.#.CreateAssetRequest.#"
	RkCreateAssetResponse = "r.bps.CreateAssetResponse.#"

	RkEmmitAssetRequest  = "r.#.EmmitAssetRequest.#"
	RkEmmitAssetResponse = "r.bps.EmmitAssetResponse.#"

	RkGetAssetsRequest  = "r.#.GetAssetsRequest.#"
	RkGetAssetsResponse = "r.bps.GetAssetsResponse.#"

	RkAddCurrencyRequest  = "r.#.AddCurrencyRequest.#"
	RkAddCurrencyResponse = "r.bps.AddCurrencyResponse.#"

	RkDeactivateAssetRequest  = "r.#.DeactivateAssetRequest.#"
	RkDeactivateAssetResponse = "r.bps.DeactivateAssetResponse.#"

	RkLockBalanceAssetRequest  = "r.#.LockBalanceAssetRequest.#"
	RkLockBalanceAssetResponse = "r.bps.LockBalanceAssetResponse.#"

	RkUnlockBalanceAssetRequest  = "r.#.UnlockBalanceAssetRequest.#"
	RkUnlockBalanceAssetResponse = "r.bps.UnlockBalanceAssetResponse.#"

	RkCreateTransferRequest = "r.#.CreateTransferRequest.#"
	RkTransferResponse      = "r.bps.TransferResponse.#"
)
