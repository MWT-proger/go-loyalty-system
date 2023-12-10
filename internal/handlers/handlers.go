package handlers

import (
	"github.com/MWT-proger/go-loyalty-system/internal/store/accountstore"
	"github.com/MWT-proger/go-loyalty-system/internal/store/withdrawalstore"
)

type WithdrawalServicer interface {
}
type AccountServicer interface {
}
type APIHandler struct {
	WithdrawalStore withdrawalstore.WithdrawalStorer
	AccountStore    accountstore.AccountStorer
	UserService     UserServicer
	OrderService    OrderServicer
	// WithdrawalService WithdrawalServicer
	// AccountService    AccountServicer
}

func NewAPIHandler(
	withdrawalstore withdrawalstore.WithdrawalStorer,
	accountstore accountstore.AccountStorer,
	userService UserServicer,
	orderService OrderServicer,
	// withdrawalService WithdrawalServicer,
	// accountService AccountServicer,
) (h *APIHandler, err error) {
	hh := &APIHandler{
		WithdrawalStore: withdrawalstore,
		AccountStore:    accountstore,
		UserService:     userService,
		OrderService:    orderService,
		// WithdrawalService: withdrawalService,
		// AccountService:    accountService,
	}

	return hh, err
}

type BaseBodyDater interface {
	IsValid() bool
}
