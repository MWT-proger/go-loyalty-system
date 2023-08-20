package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/MWT-proger/go-loyalty-system/configs"
	"github.com/MWT-proger/go-loyalty-system/internal/logger"
)

// Run() запускает сервер и слушает его по указанному хосту
func Run(r *chi.Mux) error {
	conf := configs.GetConfig()

	logger.Log.Info("Running server on", zap.String("host", conf.HostServer))
	return http.ListenAndServe(conf.HostServer, r)
}
