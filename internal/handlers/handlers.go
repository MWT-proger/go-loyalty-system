package handlers

import "github.com/MWT-proger/go-loyalty-system/internal/store/userstore"

type APIHandler struct {
	UserStore userstore.UserStorer
}

func NewAPIHandler(userStore userstore.UserStorer) (h *APIHandler, err error) {
	hh := &APIHandler{
		UserStore: userStore,
	}

	return hh, err
}

type BaseBodyDater interface {
	IsValid() bool
}
