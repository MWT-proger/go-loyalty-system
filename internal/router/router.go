package router

import (
	"github.com/go-chi/chi"

	"github.com/MWT-proger/go-loyalty-system/internal/handlers"
)

// Router() Перенаправляет запросы на необходимые хендлеры
func Router(h *handlers.APIHandler) *chi.Mux {

	r := chi.NewRouter()

	r.Post("/api/user/register", h.UserRegister)
	r.Post("/api/user/login", h.UserLogin)
	r.Post("/api/user/orders", h.SetUserOrder)
	r.Get("/api/user/orders", h.GetListOrdersUser)
	r.Get("/api/user/balance", h.GetUserBalance)
	r.Post("/api/user/balance/withdraw", h.WithdrawWithUserBalance)
	r.Get("/api/user/withdrawals", h.GetListWithdrawUserBalance)

	return r
}
