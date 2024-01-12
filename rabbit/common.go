package rabbit

import (
	"time"
	"trade-balance-service/external/balances"
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

func GetParserForCreationAssetRequest() ParserFunc[balances.BpsCreateAssetRequest] {
	return func(b []byte) (*balances.BpsCreateAssetRequest, error) {
		var request balances.BpsCreateAssetRequest
		err := proto.Unmarshal(b, &request)

		if err != nil {
			return nil, err
		}

		return &request, nil
	}
}

func GetParserForEmmitAssetRequest() ParserFunc[balances.BpsEmmitAssetRequest] {
	return func(b []byte) (*balances.BpsEmmitAssetRequest, error) {
		var request balances.BpsEmmitAssetRequest
		err := proto.Unmarshal(b, &request)

		if err != nil {
			return nil, err
		}

		return &request, nil
	}
}

func GetParserForGetAssetsById() ParserFunc[balances.BbsGetAssetInfoRequest] {
	return func(b []byte) (*balances.BbsGetAssetInfoRequest, error) {
		var request balances.BbsGetAssetInfoRequest
		err := proto.Unmarshal(b, &request)

		if err != nil {
			return nil, err
		}

		return &request, nil
	}
}
