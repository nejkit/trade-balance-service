package constants

const (
	BpsExchange = "bps.forward"
)

const (
	CreateAssetQueueName   = "q.bps.request.CreateAsset"
	EmmitAssetQueueName    = "q.bps.request.EmmitAsset"
	GetAssetsByIdQueueName = "q.bps.request.GetAssets"
)

const (
	RkCreateAssetRequest  = "r.#.CreateAssetRequest.#"
	RkCreateAssetResponse = "r.bps.CreateAssetResponse.#"

	RkEmmitAssetRequest  = "r.#.EmmitAssetRequest.#"
	RkEmmitAssetResponse = "r.bps.EmmitAssetResponse.#"

	RkGetAssetsRequest  = "r.#.GetAssetsRequest.#"
	RkGetAssetsResponse = "r.bps.GetAssetsResponse.#"
)
