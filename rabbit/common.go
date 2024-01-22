package rabbit

import (
	"time"
	"trade-balance-service/external/bps"
	"trade-balance-service/staticserr"

	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

func GetRabbitConnection(connectionString string) (*amqp091.Connection, error) {
	timeout := time.After(time.Minute * 5)
	for {
		select {
		case <-timeout:
			return nil, staticserr.ErrorRabbitConnectionFail
		default:
			connect, err := amqp091.Dial(connectionString)

			if err != nil {
				time.Sleep(time.Millisecond * 100)
				continue
			}

			return connect, nil
		}
	}
}

func GetParserForCreationAssetRequest() ParserFunc[bps.BpsCreateAssetRequest] {
	return func(b []byte) (*bps.BpsCreateAssetRequest, error) {
		var request bps.BpsCreateAssetRequest
		err := proto.Unmarshal(b, &request)

		if err != nil {
			return nil, err
		}

		return &request, nil
	}
}

func GetParserForEmmitAssetRequest() ParserFunc[bps.BpsEmmitAssetRequest] {
	return func(b []byte) (*bps.BpsEmmitAssetRequest, error) {
		var request bps.BpsEmmitAssetRequest
		err := proto.Unmarshal(b, &request)

		if err != nil {
			return nil, err
		}

		return &request, nil
	}
}

func GetParserForGetAssetsById() ParserFunc[bps.BbsGetAssetInfoRequest] {
	return func(b []byte) (*bps.BbsGetAssetInfoRequest, error) {
		var request bps.BbsGetAssetInfoRequest
		err := proto.Unmarshal(b, &request)

		if err != nil {
			return nil, err
		}

		return &request, nil
	}
}

func GetParserForDeactivateAsset() ParserFunc[bps.BpsDeactivateAssetRequest] {
	return func(b []byte) (*bps.BpsDeactivateAssetRequest, error) {
		var request bps.BpsDeactivateAssetRequest
		err := proto.Unmarshal(b, &request)

		if err != nil {
			return nil, err
		}

		return &request, nil
	}
}

func GetParserForAddCurrency() ParserFunc[bps.BpsAddCurrencyRequest] {
	return func(b []byte) (*bps.BpsAddCurrencyRequest, error) {
		var request bps.BpsAddCurrencyRequest
		err := proto.Unmarshal(b, &request)

		if err != nil {
			return nil, err
		}

		return &request, nil
	}
}

func GetParserForLockBalanceAsset() ParserFunc[bps.BpsLockBalanceRequest] {
	return func(b []byte) (*bps.BpsLockBalanceRequest, error) {
		var request bps.BpsLockBalanceRequest
		err := proto.Unmarshal(b, &request)

		if err != nil {
			return nil, err
		}

		return &request, nil
	}
}

func GetParserForRefundBalanceAsset() ParserFunc[bps.BpsRefundBalanceRequest] {
	return func(b []byte) (*bps.BpsRefundBalanceRequest, error) {
		var request bps.BpsRefundBalanceRequest
		err := proto.Unmarshal(b, &request)

		if err != nil {
			return nil, err
		}

		return &request, nil
	}
}

func GetParserForCreateTransfer() ParserFunc[bps.BpsCreateTransferRequest] {
	return func(b []byte) (*bps.BpsCreateTransferRequest, error) {
		var request bps.BpsCreateTransferRequest
		err := proto.Unmarshal(b, &request)

		if err != nil {
			return nil, err
		}

		return &request, nil
	}
}
