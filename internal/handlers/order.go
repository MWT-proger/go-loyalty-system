package handlers

import (
	"net/http"

	"github.com/MWT-proger/go-loyalty-system/internal/luhn"
)

type OrderForm struct {
	Number string
}

func (d *OrderForm) IsValid() bool {
	return luhn.Validate(d.Number)
}

// SetUserOrder(w http.ResponseWriter, r *http.Request)
// Хендлер добавляет авторизованному пользователю новый заказ
func (h *APIHandler) SetUserOrder(w http.ResponseWriter, r *http.Request) {

	var data OrderForm

	defer r.Body.Close()

	text, err := h.getTextBody(r.Body)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	data.Number = text

	if ok := data.IsValid(); !ok {
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	}
}

func (h *APIHandler) GetListOrdersUser(w http.ResponseWriter, r *http.Request) {

}
