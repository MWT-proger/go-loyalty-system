package worker

import (
	"context"
	"time"

	"github.com/MWT-proger/go-loyalty-system/internal/logger"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
	"github.com/MWT-proger/go-loyalty-system/internal/store"
)

// getListOrdersForCheck служит генаратором потока данных
// возвращает канал с Orders и блокируется до момента
// пока все полученные заказы не обработают и не обновят в БД
func (w *WorkerAccural) getListOrdersForCheck(ctx context.Context) chan *models.Order {

	ordersFromDBCh := make(chan *models.Order, 100)

	go func() {
		defer close(ordersFromDBCh)

		for {
			select {
			case <-ctx.Done():
				logger.Log.Info("ЗАКРЫТА - Задача получения заказов из БД для проверки начисления.")
				return

			default:
				w.getDataDBSemaphore.Acquire()
				logger.Log.Debug("Получение новой пачки заказов из БД")
				objs, _ := w.GetOrderLimit(ctx)

				for _, obj := range objs {
					ordersFromDBCh <- obj
				}

			}
		}
	}()

	return ordersFromDBCh
}

// getAsyncInfoOrder слушает канал ordersFromDBCh
// получает по каждому заказу информацию
// и возвращает в канал infoOrdersCh
func (w *WorkerAccural) getAsyncInfoOrder(ctx context.Context, ordersFromDBCh chan *models.Order) chan *InfoOrder {

	infoOrdersCh := make(chan *InfoOrder, 10)

	go func() {

		defer close(infoOrdersCh)
		for {
			select {
			case <-ctx.Done():
				logger.Log.Info("ЗАКРЫТА - Задача получения заказов из Accrual для проверки начисления.")
				return
			case obj := <-ordersFromDBCh:
				logger.Log.Debug("Получение информации о заказе", logger.StringField("method", obj.Number))
				infoObj, err := w.GetInfoOrder(obj.Number, obj.UserID)

				if err != nil {
					logger.Log.Error(err.Error())
					continue
				}
				infoOrdersCh <- infoObj

			}
		}

	}()
	return infoOrdersCh
}

// updateAsyncIOrderToDB слушает каналы infoOrdersCh & ordersFromDBCh.
// Обновляет пачкой Order в БД в зависимости от Status.
// Если оба канала пусты, освобождает getDataDBSemaphore
// (это позволяет получать новую пачку строк Orders)
func (w *WorkerAccural) updateAsyncIOrderToDB(
	ctx context.Context,
	infoOrdersCh chan *InfoOrder,
	ordersFromDBCh chan *models.Order) {

	go func() {
		var (
			tickerCheckProgress = time.NewTicker(5 * time.Second)
			tickerUpdateOrders  = time.NewTicker(2 * time.Second)

			listInvalidOrders    = []string{}
			listRegistredOrders  = []string{}
			listProcessingOrders = []string{}
		)
		for {
			select {
			case <-ctx.Done():
				logger.Log.Info("ЗАКРЫТА - Задача обновления(статусов и начисления) заказов в БД.")
				return

			// Распределяем объекты по "коробкам"
			case obj := <-infoOrdersCh:

				switch obj.Status {

				case Registred, NotRegistred:
					listRegistredOrders = append(listRegistredOrders, obj.Order)

				case Invaliud:
					listInvalidOrders = append(listInvalidOrders, obj.Order)

				case Processing:
					listProcessingOrders = append(listProcessingOrders, obj.Order)

				case Processed:
					fieldValue := map[string]interface{}{"updated_at": time.Now(), "status": obj.Status, "bonuses": int64(obj.Accrual * 100)}

					options := store.OptionsUpdateQuery{
						ListFieldValue: fieldValue,
						Filter:         []store.FilterParams{{Field: "number", Value: obj.Order}},
					}
					logger.Log.Debug("Обновление заказов в БД")
					w.OrderStore.UpdateOrderPlusUserAccount(ctx, &options, obj.UserID, int64(obj.Accrual*100))

				}

			// Записываем обновления в БД раз в период = timePeriodUpdatedOrdersInDB
			// обновляем updated_at и по необходимости status
			case <-tickerUpdateOrders.C:

				if listInvalidOrders != nil {
					w.updateOrdersBatch(ctx, models.Invaliud, listInvalidOrders)
					listInvalidOrders = nil
				}
				if listRegistredOrders != nil {
					w.updateOrdersBatch(ctx, "", listRegistredOrders)
					listRegistredOrders = nil
				}
				if listProcessingOrders != nil {
					w.updateOrdersBatch(ctx, models.Processing, listProcessingOrders)
					listProcessingOrders = nil
				}

			// раз в период = ticker проверяем пустоту каналов
			// и даём команду на запрос новой пачки заказов из БД
			case <-tickerCheckProgress.C:
				if len(infoOrdersCh) == 0 &&
					len(ordersFromDBCh) == 0 &&
					listInvalidOrders == nil &&
					listRegistredOrders == nil &&
					listProcessingOrders == nil {

					logger.Log.Debug("Очередной список заказов проверен и обновлен")
					w.getDataDBSemaphore.Release()

				}
			}

		}
	}()
}

// TODO: можно тоже вынести в сервисный слой
func (w *WorkerAccural) updateOrdersBatch(ctx context.Context, orderStatus models.StatusOrder, listNumberOrders []string) {

	fieldValue := map[string]interface{}{"updated_at": time.Now()}

	if orderStatus != "" {
		fieldValue["status"] = orderStatus
	}

	options := store.OptionsUpdateQuery{
		ListFieldValue: fieldValue,
		Filter:         []store.FilterParams{{Field: "number", Operator: store.FilterIN, Value: listNumberOrders}},
	}
	logger.Log.Debug("Обновление заказов в БД")
	w.OrderStore.UpdateBatch(ctx, &options)
}
