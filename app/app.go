package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trade-balance-service/constants"
	"trade-balance-service/external/balances"
	"trade-balance-service/flow"
	"trade-balance-service/handler"
	"trade-balance-service/provider"
	"trade-balance-service/rabbit"
	"trade-balance-service/services"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rabbitmq/amqp091-go"
)

func StartProgram(ctx context.Context, postgreeUrl, rabbitUrl string) {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	flow, _ := getFlow(ctxWithCancel, postgreeUrl)

	initRabbit(ctxWithCancel, rabbitUrl, *flow)
	handleGracefulShutdown(cancel)
}

func getFlow(ctx context.Context, postgreeUrl string) (*flow.Flow, error) {

	pgxPool, err := pgxpool.Connect(ctx, postgreeUrl)

	if err != nil {
		return nil, err
	}

	pgxProvider := provider.NewPgxProvider(pgxPool)

	assetsProvider := provider.NewAssetsProvider(&pgxProvider)
	balancesProvider := provider.NewBalancesProvider(&pgxProvider)
	currenciesProvider := provider.NewCurrencyProvider(&pgxProvider)

	assetService := services.NewAssetService(&assetsProvider)
	balancesService := services.NewBalanceService(&balancesProvider, &currenciesProvider)

	flow := flow.NewFlow(&assetService, &balancesService)

	return flow, nil
}

func initRabbit(ctx context.Context, rabbitUrl string, flow flow.Flow) error {
	rabbitConnection, err := rabbit.GetRabbitConnection(rabbitUrl)

	if err != nil {
		return err
	}

	commonChannel, err := rabbitConnection.Channel()

	if err != nil {
		return err
	}

	initRabbitInfrastructure(commonChannel)

	creationLisChannel, err := rabbitConnection.Channel()

	if err != nil {
		return err
	}

	emmitLisChannel, err := rabbitConnection.Channel()

	if err != nil {
		return err
	}

	getAssetsLisChannel, err := rabbitConnection.Channel()

	if err != nil {
		return err
	}

	senderChannel, err := rabbitConnection.Channel()

	if err != nil {
		return err
	}

	sender := rabbit.NewSender(ctx, senderChannel)

	handlerCollection := handler.NewHandlerCollection(&flow, &sender)

	creationAssetProcessor := rabbit.NewProcessor[balances.BpsCreateAssetRequest](rabbit.GetParserForCreationAssetRequest(), handlerCollection.HandleCreateAsset)
	emmitAssetProcessor := rabbit.NewProcessor[balances.EmmitBalanceRequest](rabbit.GetParserForEmmitAssetRequest(), handlerCollection.HandleEmmitAsset)
	getAssetsProcessor := rabbit.NewProcessor[balances.BbsGetAssetInfoRequest](rabbit.GetParserForGetAssetsById(), handlerCollection.HandleGetAssetsById)

	creationAssetListener, err := rabbit.NewListener[balances.BpsCreateAssetRequest](ctx, creationLisChannel, constants.CreateAssetQueueName, creationAssetProcessor)
	emmitAssetListener, err := rabbit.NewListener[balances.EmmitBalanceRequest](ctx, emmitLisChannel, constants.EmmitAssetQueueName, emmitAssetProcessor)
	getAssetsListener, err := rabbit.NewListener[balances.BbsGetAssetInfoRequest](ctx, getAssetsLisChannel, constants.GetAssetsByIdQueueName, getAssetsProcessor)

	go creationAssetListener.Run(ctx)
	go emmitAssetListener.Run(ctx)
	go getAssetsListener.Run(ctx)
	return nil
}

func handleGracefulShutdown(cancelFunc context.CancelFunc) {
	exit := make(chan os.Signal, 1)
	for {
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		select {
		case <-exit:
			{
				cancelFunc()
				return
			}
		default:
			time.Sleep(2 * time.Second)
		}

	}
}

func initRabbitInfrastructure(channel *amqp091.Channel) error {
	defer channel.Close()

	if err := channel.ExchangeDeclare(constants.BpsExchange, "topic", true, false, false, false, nil); err != nil {
		return err
	}
	if _, err := channel.QueueDeclare(constants.CreateAssetQueueName, true, false, false, false, nil); err != nil {
		return err
	}
	if _, err := channel.QueueDeclare(constants.EmmitAssetQueueName, true, false, false, false, nil); err != nil {
		return err
	}
	if _, err := channel.QueueDeclare(constants.GetAssetsByIdQueueName, true, false, false, false, nil); err != nil {
		return err
	}
	if err := channel.QueueBind(constants.CreateAssetQueueName, constants.RkCreateAssetRequest, constants.BpsExchange, false, nil); err != nil {
		return nil
	}
	if err := channel.QueueBind(constants.EmmitAssetQueueName, constants.RkEmmitAssetRequest, constants.BpsExchange, false, nil); err != nil {
		return nil
	}
	if err := channel.QueueBind(constants.GetAssetsByIdQueueName, constants.RkGetAssetsRequest, constants.BpsExchange, false, nil); err != nil {
		return nil
	}
	return nil
}
