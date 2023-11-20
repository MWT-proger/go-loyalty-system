package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	"github.com/MWT-proger/go-loyalty-system/configs"
	"github.com/MWT-proger/go-loyalty-system/internal/auth"
	"github.com/MWT-proger/go-loyalty-system/internal/gzip"
	"github.com/MWT-proger/go-loyalty-system/internal/handlers"
	"github.com/MWT-proger/go-loyalty-system/internal/logger"
)

// Router() Перенаправляет запросы на необходимые хендлеры
func Router(h *handlers.APIHandler) *chi.Mux {
	conf := configs.GetConfig()

	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   conf.Cors.AllowedOrigins,
		AllowedMethods:   conf.Cors.AllowedMethods,
		AllowedHeaders:   conf.Cors.AllowedHeaders,
		AllowCredentials: conf.Cors.AllowCredentials,
		Debug:            conf.Cors.Debug,
	}))
	r.Use(gzip.GzipMiddleware)

	r.Post("/api/user/register", h.UserRegister)
	r.Post("/api/user/login", h.UserLogin)

	r.Group(func(r chi.Router) {

		r.Use(auth.AuthCookieMiddleware)
		r.Post("/api/user/orders", h.SetUserOrder)
		r.Get("/api/user/orders", h.GetListOrdersUser)
		r.Get("/api/user/balance", h.GetUserBalance)
		r.Post("/api/user/balance/withdraw", h.WithdrawWithUserBalance)
		r.Get("/api/user/withdrawals", h.GetListWithdrawUserBalance)

	})

	return r
}
