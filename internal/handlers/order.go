package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/MWT-proger/go-loyalty-system/internal/auth"
	"github.com/MWT-proger/go-loyalty-system/internal/luhn"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
	"github.com/MWT-proger/go-loyalty-system/internal/store"
)

type OrderForm struct {
	Number string
}

func (d *OrderForm) IsValid() bool {
	if d.Number == "" {
		return false
	}
	return luhn.Validate(d.Number)
}

// SetUserOrder(w http.ResponseWriter, r *http.Request)
// Хендлер добавляет авторизованному пользователю новый заказ
func (h *APIHandler) SetUserOrder(w http.ResponseWriter, r *http.Request) {

	var data OrderForm
	userID, ok := auth.UserIDFrom(r.Context())

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

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
	args := map[string]interface{}{"number": data.Number}
	obj, err := h.OrderStore.GetFirstByParameters(context.TODO(), args)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if obj != nil {

		if obj.UserID != userID {
			http.Error(w, "", http.StatusConflict)
			return
		}

		http.Error(w, "", http.StatusOK)
		return

	}

	newOrder := models.NewOrder()
	newOrder.Number = data.Number
	newOrder.UserID = userID

	err = h.OrderStore.Insert(context.TODO(), newOrder)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *APIHandler) GetListOrdersUser(w http.ResponseWriter, r *http.Request) {

	userID, ok := auth.UserIDFrom(r.Context())

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	filterParams := []store.FilterParams{
		{Field: "user_id", Value: userID},
	}
	objs, err := h.OrderStore.GetAllByParameters(context.TODO(), &store.OptionsQuery{
		Filter: filterParams, Sorting: []store.SortingParams{{Key: "updated_at"}}})

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if len(objs) == 0 {
		http.Error(w, "", http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(objs)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}
