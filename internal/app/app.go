package app

import (
	"context"
	"os"

	"github.com/nel349/bz-findata/config"
	"github.com/nel349/bz-findata/internal/app/delivery/websocket"
	"github.com/nel349/bz-findata/internal/app/repository"
	"github.com/nel349/bz-findata/internal/app/usecase"
	"github.com/nel349/bz-findata/pkg/database/mysql"
	"github.com/nel349/bz-findata/pkg/exchange/coinbase"
	"github.com/nel349/bz-findata/pkg/logger/zap"
)

// Run started application
func Run(ctx context.Context, cfg *config.Config) {
	// logger
	loggerProvider := zap.NewZapLogger(cfg.Logger.Level, cfg.Logger.DisableCaller, cfg.Logger.DisableStacktrace)
	loggerProvider.InitLogger()

	// database
	dbClient, err := mysql.NewMysqlClient(cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.Base)
	if err != nil {
		loggerProvider.Fatal(err)
	}
	defer dbClient.CloseConnect()

	// exchange
	exchangeClient, err := coinbase.NewCoinbaseClient(cfg.Exchange.Url, cfg.Exchange.Protocol, cfg.Exchange.Origin)
	if err != nil {
		loggerProvider.Fatal(err)
	}
	defer exchangeClient.CloseConnection()

	// repositories & business logic
	repo := repository.NewRepositories(dbClient.DB)
	uc := usecase.NewUseCase(repo, &usecase.Packages{
		Logger: loggerProvider,
	})

	// init client
	client, err := websocket.NewSocketClient(exchangeClient, uc, loggerProvider, cfg.Exchange)
	if err != nil {
		loggerProvider.Fatal(err)
	}

	// run
	go func() {
		loggerProvider.Info("socket starting...")
		if err = client.Run(ctx); err != nil {
			loggerProvider.Fatal(err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	loggerProvider.Info("socket stopping...")
}
