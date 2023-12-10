package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gofrs/uuid"

	"github.com/MWT-proger/go-loyalty-system/internal/auth"
	lErrors "github.com/MWT-proger/go-loyalty-system/internal/errors"
	"github.com/MWT-proger/go-loyalty-system/internal/luhn"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
)

type OrderServicer interface {
	Set(ctx context.Context, userID uuid.UUID, numberOrder string) error
	GetList(ctx context.Context, userID uuid.UUID) ([]*models.Order, error)
}

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
	var ctx = r.Context()
	userID, ok := auth.UserIDFrom(ctx)

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

	err = h.OrderService.Set(ctx, userID, data.Number)

	if err != nil {
		var serviceError *lErrors.ServicesError
		if errors.As(err, &serviceError) {

			http.Error(w, serviceError.Error(), serviceError.HttpCode)
			return
		}

		http.Error(w, "Ошибка сервера, попробуйте позже.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *APIHandler) GetListOrdersUser(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()
	userID, ok := auth.UserIDFrom(ctx)

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	objs, err := h.OrderService.GetList(ctx, userID)

	if err != nil {
		var serviceError *lErrors.ServicesError
		if errors.As(err, &serviceError) {

			http.Error(w, serviceError.Error(), serviceError.HttpCode)
			return
		}

		http.Error(w, "Ошибка сервера, попробуйте позже.", http.StatusInternalServerError)
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
