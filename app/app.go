package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trade-balance-service/external/balances"
	"trade-balance-service/flow"
	"trade-balance-service/handler"
	"trade-balance-service/provider"
	"trade-balance-service/rabbit"
	"trade-balance-service/services"

	"github.com/jackc/pgx/v4/pgxpool"
)

func StartProgram(ctx context.Context, postgreeUrl, rabbitUrl string) {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	handlerCollection, _ := getHandler(ctxWithCancel, postgreeUrl)

	initRabbit(ctxWithCancel, rabbitUrl, *handlerCollection)
	go handleGracefulShutdown(cancel)
}

func getHandler(ctx context.Context, postgreeUrl string) (*handler.HandlerCollection, error) {

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

	handlerCollection := handler.NewHandlerCollection(flow)

	return &handlerCollection, nil
}

func initRabbit(ctx context.Context, rabbitUrl string, handlerCollection handler.HandlerCollection) error {
	rabbitConnection, err := rabbit.GetRabbitConnection(rabbitUrl)

	if err != nil {
		return err
	}

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

	creationAssetProcessor := rabbit.NewProcessor[balances.BpsCreateAssetRequest](rabbit.GetParserForCreationAssetRequest(), handlerCollection.HandleCreateAsset)
	emmitAssetProcessor := rabbit.NewProcessor[balances.EmmitBalanceRequest](rabbit.GetParserForEmmitAssetRequest(), handlerCollection.HandleEmmitAsset)
	getAssetsProcessor := rabbit.NewProcessor[balances.BbsGetAssetInfoRequest](rabbit.GetParserForGetAssetsById(), handlerCollection.HandleGetAssetsById)

	creationAssetListener, err := rabbit.NewListener[balances.BpsCreateAssetRequest](ctx, creationLisChannel, "", creationAssetProcessor)
	emmitAssetListener, err := rabbit.NewListener[balances.EmmitBalanceRequest](ctx, emmitLisChannel, "", emmitAssetProcessor)
	getAssetsListener, err := rabbit.NewListener[balances.BbsGetAssetInfoRequest](ctx, getAssetsLisChannel, "", getAssetsProcessor)

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
				break
			}
		default:
			time.Sleep(2 * time.Second)
		}

	}
}
