package server

import (
	"net/http"

	"github.com/MWT-proger/go-loyalty-system/configs"
	"github.com/MWT-proger/go-loyalty-system/internal/handlers"
	"github.com/MWT-proger/go-loyalty-system/internal/logger"
)

// Run() запускает сервер и слушает его по указанному хосту
func Run(h *handlers.APIHandler, conf *configs.Config) error {
	r := router(h, conf)

	logger.Log.Info("Running server on", logger.StringField("host", conf.HostServer))

	return http.ListenAndServe(conf.HostServer, r)
}
