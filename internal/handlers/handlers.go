package handlers

import (
	"github.com/MWT-proger/go-loyalty-system/internal/store/orderstore"
	"github.com/MWT-proger/go-loyalty-system/internal/store/userstore"
)

type APIHandler struct {
	UserStore  userstore.UserStorer
	OrderStore orderstore.OrderStorer
}

func NewAPIHandler(userStore userstore.UserStorer, orderstore orderstore.OrderStorer) (h *APIHandler, err error) {
	hh := &APIHandler{
		UserStore:  userStore,
		OrderStore: orderstore,
	}

	return hh, err
}

type BaseBodyDater interface {
	IsValid() bool
}
