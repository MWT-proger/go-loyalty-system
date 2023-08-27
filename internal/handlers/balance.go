package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/MWT-proger/go-loyalty-system/internal/logger"
	"github.com/MWT-proger/go-loyalty-system/internal/luhn"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
	"github.com/MWT-proger/go-loyalty-system/internal/request"
)

type WithdrawForm struct {
	Order string `json:"order"`
	Sum   int64  `json:"sum"`
}

type UserBalance struct {
	Current   int64 `json:"current"`
	Withdrawn int64 `json:"withdrawn"`
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

	sumAmmount, err := h.OrderStore.GetSumByParameters(context.TODO(), args)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	sumWithdrawn, err := h.WithdrawalStore.GetSumByParameters(context.TODO(), args)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	userBalance := UserBalance{
		Current:   sumAmmount - sumWithdrawn,
		Withdrawn: sumWithdrawn,
	}

	resp, err := json.Marshal(userBalance)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

	logger.Log.Debug("", zap.Int64("ammount", sumAmmount), zap.Int64("withdrawn", sumWithdrawn))

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
