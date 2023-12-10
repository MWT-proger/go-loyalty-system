package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/MWT-proger/go-loyalty-system/internal/auth"
	lErrors "github.com/MWT-proger/go-loyalty-system/internal/errors"
)

type UserServicer interface {
	UserLogin(ctx context.Context, login string, password string) (string, error)
	UserRegister(ctx context.Context, login string, password string) (string, error)
}

type UserFormRegister struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

func (d *UserFormRegister) IsValid() bool {
	return auth.ValidatePassword(d.Password)
}

type UserFormLogin struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

func (d *UserFormLogin) IsValid() bool {
	return auth.ValidatePassword(d.Password)
}

func (h *APIHandler) UserRegister(w http.ResponseWriter, r *http.Request) {
	var (
		data UserFormRegister
		ctx  = r.Context()
	)

	if ok := h.getBodyData(w, r, &data); !ok {
		return
	}

	tokenString, err := h.UserService.UserRegister(ctx, data.Login, data.Password)

	if err != nil {
		var serviceError *lErrors.ServicesError
		if errors.As(err, &serviceError) {

			http.Error(w, serviceError.Error(), serviceError.HttpCode)
			return
		}

		http.Error(w, "Ошибка сервера, попробуйте позже.", http.StatusInternalServerError)
		return
	}

	auth.SetAuthTokenToCookie(w, tokenString)

	w.WriteHeader(http.StatusOK)

}

func (h *APIHandler) UserLogin(w http.ResponseWriter, r *http.Request) {

	var (
		data UserFormLogin
		ctx  = r.Context()
	)

	if ok := h.getBodyData(w, r, &data); !ok {
		return
	}
	tokenString, err := h.UserService.UserLogin(ctx, data.Login, data.Password)

	if err != nil {

		var serviceError *lErrors.ServicesError
		if errors.As(err, &serviceError) {

			http.Error(w, serviceError.Error(), serviceError.HttpCode)
			return
		}

		http.Error(w, "Ошибка сервера, попробуйте позже.", http.StatusInternalServerError)
		return
	}

	auth.SetAuthTokenToCookie(w, tokenString)

	w.WriteHeader(http.StatusOK)

}
