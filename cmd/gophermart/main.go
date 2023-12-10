package main

import (
	"context"

	"github.com/MWT-proger/go-loyalty-system/configs"
	"github.com/MWT-proger/go-loyalty-system/internal/handlers"
	"github.com/MWT-proger/go-loyalty-system/internal/logger"
	"github.com/MWT-proger/go-loyalty-system/internal/server"
	"github.com/MWT-proger/go-loyalty-system/internal/services"
	"github.com/MWT-proger/go-loyalty-system/internal/store"
	"github.com/MWT-proger/go-loyalty-system/internal/store/accountstore"
	"github.com/MWT-proger/go-loyalty-system/internal/store/orderstore"
	"github.com/MWT-proger/go-loyalty-system/internal/store/userstore"
	"github.com/MWT-proger/go-loyalty-system/internal/store/withdrawalstore"
	"github.com/MWT-proger/go-loyalty-system/internal/worker"
)

var storage store.Store

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx); err != nil {
		cancel()
		panic(err)
	}
}

// run() выполняет все предворительные действия и вызывает функцию запуска сервера
func run(ctx context.Context) error {

	var conf = configs.InitConfig()

	if err := logger.Initialize(conf.LogLevel); err != nil {
		return err
	}

	storage, err := store.NewStore(ctx, conf)

	if err != nil {
		return err
	}

	userStore := userstore.New(storage)
	orderStore := orderstore.New(storage)
	withdrawalStore := withdrawalstore.New(storage)
	accountStore := accountstore.New(storage)

	userService := services.NewUserService(userStore, conf)
	orderService := services.NewOrderService(orderStore)
	accountService := services.NewAccountService(accountStore)
	withdrawalService := services.NewWithdrawalService(withdrawalStore)

	h, err := handlers.NewAPIHandler(userService, orderService, withdrawalService, accountService)

	if err != nil {
		return err
	}

	w, err := worker.NewWorkerAccural(conf, orderStore, withdrawalStore, accountStore)

	if err != nil {
		return err
	}

	err = w.Init(ctx)

	if err != nil {
		return err
	}

	err = server.Run(h, conf)

	if err != nil {
		return err
	}

	return nil
}
