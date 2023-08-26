package handlers

import (
	"context"
	"net/http"

	"github.com/MWT-proger/go-loyalty-system/internal/auth"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
)

type UserRegister struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

func (d *UserRegister) IsValid() bool {
	return auth.ValidatePassword(d.Password)
}

func (h *APIHandler) UserRegister(w http.ResponseWriter, r *http.Request) {
	var data UserRegister

	if ok := h.getBodyData(w, r, &data); !ok {
		return
	}

	newUser := models.NewUser()
	newUser.Login = data.Login

	args := map[string]interface{}{"login": newUser.Login}
	obj, err := h.UserStore.GetFirstByParameters(context.TODO(), args)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if obj != nil {
		http.Error(w, "", http.StatusConflict)
		return
	}

	newUser.Password, err = auth.HashPassword(data.Password)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	err = h.UserStore.Insert(context.TODO(), newUser)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	tokenString, err := auth.BuildJWTString(newUser.ID)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	auth.SetAuthTokenToCookie(w, tokenString)

	w.WriteHeader(http.StatusOK)

}

func (h *APIHandler) UserLogin(w http.ResponseWriter, r *http.Request) {

}
