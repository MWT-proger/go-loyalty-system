package handlers

import (
	"context"
	"net/http"

	"github.com/MWT-proger/go-loyalty-system/internal/auth"
)

type UserServicer interface {
	UserLogin(ctx context.Context, login string, password string) (string, error)
	UserRegister(ctx context.Context, login string, password string) (string, error)
}

type userFormRegister struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type userFormLogin struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

func (h *APIHandler) UserRegister(w http.ResponseWriter, r *http.Request) {
	var (
		data userFormRegister
		ctx  = r.Context()
	)

	if ok := h.getBodyData(w, r, &data); !ok {
		return
	}

	tokenString, err := h.UserService.UserRegister(ctx, data.Login, data.Password)

	if err != nil {
		h.setHttpError(w, err)
		return
	}

	auth.SetAuthTokenToCookie(w, tokenString)

	w.WriteHeader(http.StatusOK)

}

func (h *APIHandler) UserLogin(w http.ResponseWriter, r *http.Request) {

	var (
		data userFormLogin
		ctx  = r.Context()
	)

	if ok := h.getBodyData(w, r, &data); !ok {
		return
	}
	tokenString, err := h.UserService.UserLogin(ctx, data.Login, data.Password)

	if err != nil {
		h.setHttpError(w, err)
		return
	}

	auth.SetAuthTokenToCookie(w, tokenString)

	w.WriteHeader(http.StatusOK)

}

func (d *userFormRegister) isValid() bool {
	return auth.ValidatePassword(d.Password)
}

func (d *userFormLogin) isValid() bool {
	return auth.ValidatePassword(d.Password)
}
