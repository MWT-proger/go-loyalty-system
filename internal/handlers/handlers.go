package handlers

import (
	"github.com/MWT-proger/go-loyalty-system/internal/store/accountstore"
	"github.com/MWT-proger/go-loyalty-system/internal/store/orderstore"
	"github.com/MWT-proger/go-loyalty-system/internal/store/userstore"
	"github.com/MWT-proger/go-loyalty-system/internal/store/withdrawalstore"
)

type OrderServicer interface {
}
type WithdrawalServicer interface {
}
type AccountServicer interface {
}
type APIHandler struct {
	UserStore       userstore.UserStorer
	OrderStore      orderstore.OrderStorer
	WithdrawalStore withdrawalstore.WithdrawalStorer
	AccountStore    accountstore.AccountStorer
	UserService     UserServicer
	// OrderService      OrderServicer
	// WithdrawalService WithdrawalServicer
	// AccountService    AccountServicer
}

func NewAPIHandler(
	orderstore orderstore.OrderStorer,
	withdrawalstore withdrawalstore.WithdrawalStorer,
	accountstore accountstore.AccountStorer,
	userService UserServicer,
	// orderService OrderServicer,
	// withdrawalService WithdrawalServicer,
	// accountService AccountServicer,
) (h *APIHandler, err error) {
	hh := &APIHandler{
		OrderStore:      orderstore,
		WithdrawalStore: withdrawalstore,
		AccountStore:    accountstore,
		UserService:     userService,
		// OrderService:      orderService,
		// WithdrawalService: withdrawalService,
		// AccountService:    accountService,
	}

	return hh, err
}

type BaseBodyDater interface {
	IsValid() bool
}
