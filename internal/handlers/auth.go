package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/MWT-proger/go-loyalty-system/internal/models"
	"github.com/MWT-proger/go-loyalty-system/internal/utils"
)

type UserRegister struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

func (d *UserRegister) IsValid() bool {
	return utils.ValidatePassword(d.Password)
}

func (h *APIHandler) UserRegister(w http.ResponseWriter, r *http.Request) {

	var data UserRegister
	newUser := models.User{}

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

	// TODO: Генерить хэш
	// TODO: генерить токен
	// TODO: авторизовывать пользователя
	newUser.Password = data.Password

	err = h.UserStore.Insert(context.TODO(), &newUser)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

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
