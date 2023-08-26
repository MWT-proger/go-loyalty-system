package handlers

import (
	"context"
	"encoding/json"
	"fmt"
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
	newUser := models.NewUser()

	if ok := h.getBodyData(w, r, &data); !ok {
		return
	}

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
	fmt.Println(newUser)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// TODO: генерить токен
	// TODO: авторизовывать пользователя

	resp, err := json.Marshal(newUser)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

func (h *APIHandler) UserLogin(w http.ResponseWriter, r *http.Request) {

}
