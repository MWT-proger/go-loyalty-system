package handlers

import (
	"github.com/MWT-proger/go-loyalty-system/internal/store/accountstore"
	"github.com/MWT-proger/go-loyalty-system/internal/store/orderstore"
	"github.com/MWT-proger/go-loyalty-system/internal/store/userstore"
	"github.com/MWT-proger/go-loyalty-system/internal/store/withdrawalstore"
)

type APIHandler struct {
	UserStore       userstore.UserStorer
	OrderStore      orderstore.OrderStorer
	WithdrawalStore withdrawalstore.WithdrawalStorer
	AccountStore    accountstore.AccountStorer
}

func NewAPIHandler(
	userStore userstore.UserStorer,
	orderstore orderstore.OrderStorer,
	withdrawalstore withdrawalstore.WithdrawalStorer,
	accountstore accountstore.AccountStorer,
) (h *APIHandler, err error) {
	hh := &APIHandler{
		UserStore:       userStore,
		OrderStore:      orderstore,
		WithdrawalStore: withdrawalstore,
		AccountStore:    accountstore,
	}

	return hh, err
}

type BaseBodyDater interface {
	IsValid() bool
}
