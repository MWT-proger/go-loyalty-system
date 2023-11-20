package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/MWT-proger/go-loyalty-system/configs"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
	"github.com/MWT-proger/go-loyalty-system/internal/store"
	"github.com/MWT-proger/go-loyalty-system/internal/store/accountstore"
	"github.com/MWT-proger/go-loyalty-system/internal/store/orderstore"
	"github.com/MWT-proger/go-loyalty-system/internal/store/withdrawalstore"
)

type StatusOrderAccural string

type InfoOrder struct {
	Order   string             `json:"order"`
	Status  StatusOrderAccural `json:"status"`
	Accrual float64            `json:"accrual,omitempty"`
}

type WorkerAccural struct {
	OrderStore      orderstore.OrderStorer
	WithdrawalStore withdrawalstore.WithdrawalStorer
	AccountStore    accountstore.AccountStorer
	client          *http.Client
	baseURL         string
}

const (
	Registred  StatusOrderAccural = "REGISTERED"
	Processing StatusOrderAccural = "PROCESSING"
	Invaliud   StatusOrderAccural = "INVALID"
	Processed  StatusOrderAccural = "PROCESSED"
)

func NewWorkerAccural(
	orderstore orderstore.OrderStorer,
	withdrawalstore withdrawalstore.WithdrawalStorer,
	accountstore accountstore.AccountStorer,
) (w *WorkerAccural, err error) {

	conf := configs.GetConfig()

	ww := &WorkerAccural{
		OrderStore:      orderstore,
		WithdrawalStore: withdrawalstore,
		AccountStore:    accountstore,
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return nil
			}},
		baseURL: conf.AccuralSystemAddress,
	}

	return ww, err
}

func (w *WorkerAccural) GetInfoOrder(numberOrder string) (*InfoOrder, error) {

	var data InfoOrder

	response, err := w.client.Get(w.baseURL + "/api/orders/" + numberOrder)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if err := w.unmarshalBody(response.Body, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (w *WorkerAccural) GetOrderLimit() ([]*models.Order, error) {

	args := map[string]interface{}{}

	objs, err := w.OrderStore.GetAllByParameters(
		context.TODO(),
		&store.OptionsSelect{
			Args: args, Limit: 10, OrderBy: "updated_at", DescOrderBy: true,
		})

	if err != nil {
		return nil, err
	}

	return objs, nil
}

func (w *WorkerAccural) unmarshalBody(body io.ReadCloser, form interface{}) error {

	defer body.Close()

	var buf bytes.Buffer
	_, err := buf.ReadFrom(body)

	if err != nil {
		return err
	}

	if err = json.Unmarshal(buf.Bytes(), form); err != nil {
		return err
	}

	return nil
}
