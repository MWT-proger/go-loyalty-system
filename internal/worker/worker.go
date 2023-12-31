package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gofrs/uuid"

	"github.com/MWT-proger/go-loyalty-system/configs"
	"github.com/MWT-proger/go-loyalty-system/internal/errors"
	"github.com/MWT-proger/go-loyalty-system/internal/logger"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
	"github.com/MWT-proger/go-loyalty-system/internal/store"
)

type StatusOrderAccural string

// InfoOrder Структура данных API Accrual
type InfoOrder struct {
	Order   string             `json:"order"`
	Status  StatusOrderAccural `json:"status"`
	Accrual float64            `json:"accrual,omitempty"`
	UserID  uuid.UUID
}

type OrderStorer interface {
	GetAllByParameters(ctx context.Context, options *store.OptionsQuery) ([]*models.Order, error)
	UpdateBatch(ctx context.Context, options *store.OptionsUpdateQuery) error
	UpdateOrderPlusUserAccount(ctx context.Context, options *store.OptionsUpdateQuery, userID uuid.UUID, bonuses int64) error
}

// WorkerAccural структура отвечает параллельную работу с заказами
// OrderStore, WithdrawalStore, AccountStore репозитории объектов в БД
// getDataDBSemaphore семафор ограничивающий колличество запросов к БД
type WorkerAccural struct {
	OrderStore         OrderStorer
	client             *http.Client
	baseURL            string
	getDataDBSemaphore Semaphore
}

// StatusOrderAccural Статусы
// варианты поля status ответа сервиса Accrual  GET /api/orders/{order}
const (
	NotRegistred StatusOrderAccural = "NOT_REGISTERED"
	Registred    StatusOrderAccural = "REGISTERED"
	Processing   StatusOrderAccural = "PROCESSING"
	Invaliud     StatusOrderAccural = "INVALID"
	Processed    StatusOrderAccural = "PROCESSED"
)

func NewWorkerAccural(
	conf *configs.Config,
	orderstore OrderStorer,
) (w *WorkerAccural, err error) {

	ww := &WorkerAccural{
		OrderStore: orderstore,
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return nil
			}},
		baseURL:            conf.AccuralSystemAddress,
		getDataDBSemaphore: *NewSemaphore(1),
	}

	return ww, err
}

// Init запускает горутину вечно проверяющию и обновляющую заказы
func (w *WorkerAccural) Init(ctx context.Context) error {
	logger.Log.Info("WorkerAccural Init - запуск воркера в отдельном потоке")
	go w.StartEternalCycle(ctx)
	return nil
}

// StartEternalCycle Запускает конвейр параллельной обработки заказов
func (w *WorkerAccural) StartEternalCycle(ctx context.Context) {

	ordersFromDBCh := w.getListOrdersForCheck(ctx)

	ordersFromAccrual := w.getAsyncInfoOrder(ctx, ordersFromDBCh)

	w.updateAsyncIOrderToDB(ctx, ordersFromAccrual, ordersFromDBCh)

}

// GetInfoOrder(numberOrder string) (*InfoOrder, error)
// Получает информацию о заказе в Accrual сервисе
// по номеру заказа и возвращает структуру InfoOrder
// TODO: можно вынести в пакет client
func (w *WorkerAccural) GetInfoOrder(numberOrder string, userID uuid.UUID) (*InfoOrder, error) {

	var data InfoOrder

	response, err := w.client.Get(w.baseURL + "/api/orders/" + numberOrder)
	if err != nil {
		return nil, err
	}

	logger.Log.Debug(
		"Ответ сервиса Accrual ",
		logger.StringField("Заказ", numberOrder),
		logger.IntField("Статус ответа", response.StatusCode),
	)

	switch response.StatusCode {

	case 204:
		data = InfoOrder{Order: numberOrder, Status: NotRegistred}

	case 500:
		err := errors.ErrorAccrualStatusCode500{}
		return nil, &err
	case 429:
		err := errors.ErrorAccrualStatusCode429{}
		return nil, &err
	case 200:
		defer response.Body.Close()

		if err := w.unmarshalBody(response.Body, &data); err != nil {
			return nil, err
		}

	}
	data.UserID = userID
	return &data, nil
}

// GetOrderLimit() ([]*models.Order, error) Достает из БД
// Заказы со статусами (New, Processing) в количестве равном Limit
// TODO: можно тоже вынести в сервисный слой
func (w *WorkerAccural) GetOrderLimit(ctx context.Context) ([]*models.Order, error) {

	objs, err := w.OrderStore.GetAllByParameters(
		ctx,
		&store.OptionsQuery{
			Filter: []store.FilterParams{
				{
					Field:    "status",
					Value:    []models.StatusOrder{models.New, models.Processing},
					Operator: store.FilterIN,
				},
			},
			Sorting: []store.SortingParams{
				{Key: "updated_at", Desc: true},
			},
			Limit: 10,
		})

	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}

	return objs, nil
}

// unmarshalBody(body io.ReadCloser, form interface{}) error -
// парсит body и записывает резултат в form
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
