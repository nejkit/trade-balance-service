package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trade-balance-service/constants"
	"trade-balance-service/external/bps"
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
	initHandler(ctxWithCancel, rabbitUrl, postgreeUrl)
	handleGracefulShutdown(cancel)
}

func getFlow(ctx context.Context, postgreeUrl string, sender *rabbit.Sender) (*flow.Flow, error) {

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

	flow := flow.NewFlow(&assetService, &balancesService, sender)

	return flow, nil
}

func initHandler(ctx context.Context, rabbitUrl string, postgreeUrl string) error {
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

	addCurrencyLisChannel, err := rabbitConnection.Channel()

	if err != nil {
		return err
	}

	deactivateAssetLisChannel, err := rabbitConnection.Channel()

	if err != nil {
		return err
	}

	lockBalanceAssetLisChannel, err := rabbitConnection.Channel()

	if err != nil {
		return err
	}

	refundBalanceAssetLisChannel, err := rabbitConnection.Channel()

	if err != nil {
		return err
	}

	transferLisChannel, err := rabbitConnection.Channel()

	if err != nil {
		return err
	}

	senderChannel, err := rabbitConnection.Channel()

	if err != nil {
		return err
	}

	sender := rabbit.NewSender(ctx, senderChannel)

	flow, err := getFlow(ctx, postgreeUrl, &sender)

	handlerCollection := handler.NewHandlerCollection(flow, &sender)

	creationAssetProcessor := rabbit.NewProcessor[bps.BpsCreateAssetRequest](rabbit.GetParserForCreationAssetRequest(), handlerCollection.HandleCreateAsset)
	emmitAssetProcessor := rabbit.NewProcessor[bps.BpsEmmitAssetRequest](rabbit.GetParserForEmmitAssetRequest(), handlerCollection.HandleEmmitAsset)
	getAssetsProcessor := rabbit.NewProcessor[bps.BbsGetAssetInfoRequest](rabbit.GetParserForGetAssetsById(), handlerCollection.HandleGetAssetsById)
	deactivateAssetProcessor := rabbit.NewProcessor[bps.BpsDeactivateAssetRequest](rabbit.GetParserForDeactivateAsset(), handlerCollection.HandleDeactivateAsset)
	addCurrencyProcessor := rabbit.NewProcessor[bps.BpsAddCurrencyRequest](rabbit.GetParserForAddCurrency(), handlerCollection.HandleAddCurrency)
	lockBalanceProcessor := rabbit.NewProcessor[bps.BpsLockBalanceRequest](rabbit.GetParserForLockBalanceAsset(), handlerCollection.HandleLockBalanceAsset)
	refundBalanceProcessor := rabbit.NewProcessor[bps.BpsRefundBalanceRequest](rabbit.GetParserForRefundBalanceAsset(), handlerCollection.HandleRefundBalanceAsset)
	transferProcessor := rabbit.NewProcessor[bps.BpsCreateTransferRequest](rabbit.GetParserForCreateTransfer(), handlerCollection.HandleCreateTransfer)

	creationAssetListener, err := rabbit.NewListener[bps.BpsCreateAssetRequest](ctx, creationLisChannel, constants.CreateAssetQueueName, creationAssetProcessor)
	emmitAssetListener, err := rabbit.NewListener[bps.BpsEmmitAssetRequest](ctx, emmitLisChannel, constants.EmmitAssetQueueName, emmitAssetProcessor)
	getAssetsListener, err := rabbit.NewListener[bps.BbsGetAssetInfoRequest](ctx, getAssetsLisChannel, constants.GetAssetsByIdQueueName, getAssetsProcessor)
	deactivateAssetListener, err := rabbit.NewListener[bps.BpsDeactivateAssetRequest](ctx, deactivateAssetLisChannel, constants.DeactivateAssetQueueName, deactivateAssetProcessor)
	addCurrencyListener, err := rabbit.NewListener[bps.BpsAddCurrencyRequest](ctx, addCurrencyLisChannel, constants.AddNewCurrencyQueueName, addCurrencyProcessor)
	lockBalanceListener, err := rabbit.NewListener[bps.BpsLockBalanceRequest](ctx, lockBalanceAssetLisChannel, constants.LockBalanceAssetQueueName, lockBalanceProcessor)
	refundBalanceListener, err := rabbit.NewListener[bps.BpsRefundBalanceRequest](ctx, refundBalanceAssetLisChannel, constants.UnlockBalanceAssetQueueName, refundBalanceProcessor)
	transferListener, err := rabbit.NewListener[bps.BpsCreateTransferRequest](ctx, transferLisChannel, constants.CreateTransferQueueName, transferProcessor)

	if err != nil {
		return err
	}

	go creationAssetListener.Run(ctx)
	go emmitAssetListener.Run(ctx)
	go getAssetsListener.Run(ctx)
	go deactivateAssetListener.Run(ctx)
	go addCurrencyListener.Run(ctx)
	go lockBalanceListener.Run(ctx)
	go refundBalanceListener.Run(ctx)
	go transferListener.Run(ctx)

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
	if _, err := channel.QueueDeclare(constants.DeactivateAssetQueueName, true, false, false, false, nil); err != nil {
		return err
	}
	if _, err := channel.QueueDeclare(constants.AddNewCurrencyQueueName, true, false, false, false, nil); err != nil {
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
	if err := channel.QueueBind(constants.AddNewCurrencyQueueName, constants.RkAddCurrencyRequest, constants.BpsExchange, false, nil); err != nil {
		return nil
	}
	if err := channel.QueueBind(constants.DeactivateAssetQueueName, constants.RkDeactivateAssetRequest, constants.BpsExchange, false, nil); err != nil {
		return nil
	}
	return nil
}
