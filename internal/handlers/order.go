package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"

	"github.com/MWT-proger/go-loyalty-system/internal/auth"
	"github.com/MWT-proger/go-loyalty-system/internal/luhn"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
)

type OrderServicer interface {
	Set(ctx context.Context, userID uuid.UUID, numberOrder string) error
	GetList(ctx context.Context, userID uuid.UUID) ([]*models.Order, error)
}

type orderForm struct {
	Number string
}

// SetUserOrder(w http.ResponseWriter, r *http.Request)
// Хендлер добавляет авторизованному пользователю новый заказ
func (h *APIHandler) SetUserOrder(w http.ResponseWriter, r *http.Request) {

	var (
		data       orderForm
		ctx        = r.Context()
		userID, ok = auth.UserIDFrom(ctx)
	)

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

	if ok := data.isValid(); !ok {
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	}

	err = h.OrderService.Set(ctx, userID, data.Number)

	if err != nil {
		h.setHttpError(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *APIHandler) GetListOrdersUser(w http.ResponseWriter, r *http.Request) {

	var (
		ctx        = r.Context()
		userID, ok = auth.UserIDFrom(ctx)
	)

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	objs, err := h.OrderService.GetList(ctx, userID)

	if err != nil {
		h.setHttpError(w, err)
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

func (d *orderForm) isValid() bool {
	if d.Number == "" {
		return false
	}
	return luhn.Validate(d.Number)
}
