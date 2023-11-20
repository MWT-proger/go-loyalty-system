package main

import (
	"context"

	"github.com/MWT-proger/go-loyalty-system/configs"
	"github.com/MWT-proger/go-loyalty-system/internal/handlers"
	"github.com/MWT-proger/go-loyalty-system/internal/logger"
	"github.com/MWT-proger/go-loyalty-system/internal/router"
	"github.com/MWT-proger/go-loyalty-system/internal/server"
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

// initProject() иницилизирует все необходимые переменный проекта
func initProject(ctx context.Context) error {
	var configInit = configs.InitConfig()
	storage = store.Store{}

	parseFlags(configInit)

	conf := configs.SetConfigFromEnv()

	if err := logger.Initialize(conf.LogLevel); err != nil {
		return err
	}

	if err := storage.Init(ctx); err != nil {
		return err
	}

	return nil
}

// run() выполняет все предворительные действия и вызывает функцию запуска сервера
func run(ctx context.Context) error {
	err := initProject(ctx)

	if err != nil {
		return err
	}
	userStore := userstore.New(&storage)
	orderstore := orderstore.New(&storage)
	withdrawalstore := withdrawalstore.New(&storage)
	accountstore := accountstore.New(&storage)
	h, err := handlers.NewAPIHandler(userStore, orderstore, withdrawalstore, accountstore)

	if err != nil {
		return err
	}

	w, err := worker.NewWorkerAccural(orderstore, withdrawalstore, accountstore)

	// Временный блок ()запускать в горутине бесконечный цикл
	w.CheckInfoAndUpdateOrder()
	// Временный блок

	if err != nil {
		return err
	}

	r := router.Router(h)

	err = server.Run(r)

	if err != nil {
		return err
	}

	return nil
}
