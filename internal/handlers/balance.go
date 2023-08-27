package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	lErrors "github.com/MWT-proger/go-loyalty-system/internal/errors"
	"github.com/MWT-proger/go-loyalty-system/internal/luhn"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
	"github.com/MWT-proger/go-loyalty-system/internal/request"
)

type WithdrawForm struct {
	Order string `json:"order"`
	Sum   int64  `json:"-"`
}

func (d *WithdrawForm) UnmarshalJSON(data []byte) error {

	type alias WithdrawForm

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

func (d *WithdrawForm) IsValid() bool {
	return luhn.Validate(d.Order)

}

func (h *APIHandler) GetUserBalance(w http.ResponseWriter, r *http.Request) {

	userID, ok := request.UserIDFrom(r.Context())

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	args := map[string]interface{}{"user_id": userID}

	obj, err := h.AccountStore.GetFirstByParameters(context.TODO(), args)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if obj == nil {

		obj = models.NewAccount()
		obj.UserID = userID

		if err := h.AccountStore.Insert(context.TODO(), obj); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

	}

	resp, err := json.Marshal(obj)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

func (h *APIHandler) WithdrawWithUserBalance(w http.ResponseWriter, r *http.Request) {

	var data WithdrawForm
	userID, ok := request.UserIDFrom(r.Context())

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	if err := h.unmarshalBody(r.Body, &data); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if ok := data.IsValid(); !ok {
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	}
	args := map[string]interface{}{"number": data.Order}
	obj, err := h.WithdrawalStore.GetFirstByParameters(context.TODO(), args)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if obj != nil {

		if obj.UserID != userID || obj.Bonuses.Int64 != data.Sum {
			http.Error(w, "", http.StatusConflict)
			return
		}

		http.Error(w, "", http.StatusOK)
		return

	}

	newWithdrawal := models.NewWithdrawal()
	newWithdrawal.Number = data.Order
	newWithdrawal.UserID = userID
	newWithdrawal.Bonuses = sql.NullInt64{Int64: data.Sum, Valid: true}

	err = h.WithdrawalStore.Insert(context.TODO(), newWithdrawal)

	if err != nil {

		if errors.Is(err, &lErrors.ErrorNotBonuses{}) {
			http.Error(w, "", http.StatusPaymentRequired)
			return
		}

		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	fmt.Println(newWithdrawal)
	w.WriteHeader(http.StatusOK)
}

func (h *APIHandler) GetListWithdrawUserBalance(w http.ResponseWriter, r *http.Request) {

	userID, ok := request.UserIDFrom(r.Context())

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	args := map[string]interface{}{"user_id": userID}
	objs, err := h.WithdrawalStore.GetAllByParameters(context.TODO(), args)

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
