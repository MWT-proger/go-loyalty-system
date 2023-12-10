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

type AccountServicer interface {
	GetOrSet(ctx context.Context, userID uuid.UUID) (*models.Account, error)
}
type WithdrawalServicer interface {
	GetList(ctx context.Context, userID uuid.UUID) ([]*models.Withdrawal, error)
	Set(ctx context.Context, userID uuid.UUID, numberOrder string, sum int64) error
}

type withdrawForm struct {
	Order string `json:"order"`
	Sum   int64  `json:"-"`
}

func (h *APIHandler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	var (
		ctx        = r.Context()
		userID, ok = auth.UserIDFrom(ctx)
	)

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	obj, err := h.AccountService.GetOrSet(ctx, userID)

	if err != nil {
		h.setHttpError(w, err)
		return
	}

	resp, err := json.Marshal(obj)

	if err != nil {
		http.Error(w, "Ошибка сервера, попробуйте позже.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

func (h *APIHandler) WithdrawWithUserBalance(w http.ResponseWriter, r *http.Request) {
	var (
		ctx        = r.Context()
		data       withdrawForm
		userID, ok = auth.UserIDFrom(r.Context())
	)

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	if err := h.unmarshalBody(r.Body, &data); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if ok := data.isValid(); !ok {
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	}
	err := h.WithdrawalService.Set(ctx, userID, data.Order, data.Sum)

	if err != nil {
		h.setHttpError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *APIHandler) GetListWithdrawUserBalance(w http.ResponseWriter, r *http.Request) {

	var (
		ctx        = r.Context()
		userID, ok = auth.UserIDFrom(r.Context())
	)

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	objs, err := h.WithdrawalService.GetList(ctx, userID)

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

// UnmarshalJSON(data []byte) подгатавливает данные для вывода на клиент,
// поэтому этот метод публичный
func (d *withdrawForm) UnmarshalJSON(data []byte) error {

	type alias withdrawForm

	aliasValue := &struct {
		*alias
		Sum float64 `json:"sum"`
	}{
		alias: (*alias)(d),
	}

	if err := json.Unmarshal(data, aliasValue); err != nil {
		return err
	}

	d.Sum = int64(aliasValue.Sum * 100)

	return nil

}

func (d *withdrawForm) isValid() bool {
	return luhn.Validate(d.Order)

}
