package handlers

import (
	"context"
	"net/http"

	"github.com/MWT-proger/go-loyalty-system/internal/auth"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
)

type UserFormAuth struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

func (d *UserFormAuth) IsValid() bool {
	return auth.ValidatePassword(d.Password)
}

func (h *APIHandler) UserRegister(w http.ResponseWriter, r *http.Request) {
	var data UserFormAuth

	if ok := h.getBodyData(w, r, &data); !ok {
		return
	}

	newUser := models.NewUser()
	newUser.Login = data.Login

	args := map[string]interface{}{"login": newUser.Login}
	objs, err := h.UserStore.GetFirstByParameters(context.TODO(), args)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if len(objs) > 0 {
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

	var data UserFormAuth

	if ok := h.getBodyData(w, r, &data); !ok {
		return
	}

	args := map[string]interface{}{"login": data.Login}
	objs, err := h.UserStore.GetFirstByParameters(context.TODO(), args)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if len(objs) == 0 {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}
	user := objs[0]

	if ok := auth.CheckPasswordHash(data.Password, user.Password); !ok {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	tokenString, err := auth.BuildJWTString(user.ID)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	auth.SetAuthTokenToCookie(w, tokenString)

	w.WriteHeader(http.StatusOK)

}
