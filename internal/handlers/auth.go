package handlers

import (
	"encoding/json"
	"net/http"

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

	if ok := h.getBodyData(w, r, &data); !ok {
		return
	}
	// Тут процесс
	resp, err := json.Marshal(data)
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
