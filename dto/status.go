package dto

import "trade-balance-service/external/bps"

type TransferState struct {
	State bps.BpsTransferState
	Err   error
}
