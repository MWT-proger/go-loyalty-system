package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/MWT-proger/go-loyalty-system/internal/logger"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
)

// getListOrdersForCheck служит генаратором потока данных
// возвращает канал с Orders и блокируется до момента
// пока все полученные заказы не обработают и не обновят в БД
func (w *WorkerAccural) getListOrdersForCheck(ctx context.Context) chan *models.Order {

	ordersFromDBCh := make(chan *models.Order)

	go func() {
		defer close(ordersFromDBCh)

		for {
			select {
			case <-ctx.Done():
				logger.Log.Info("ЗАКРЫТА - Задача получения заказов из БД для проверки начисления.")
				return

			default:
				w.getDataDBSemaphore.Acquire()
				// TODO: Необходимо логировать
				objs, _ := w.GetOrderLimit()

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

	infoOrdersCh := make(chan *InfoOrder)

	go func() {

		defer close(infoOrdersCh)
		for {
			// TODO: Необходимо логировать
			select {
			case <-ctx.Done():
				logger.Log.Info("ЗАКРЫТА - Задача получения заказов из Accrual для проверки начисления.")
				return
			case obj := <-ordersFromDBCh:
				infoObj, err := w.GetInfoOrder(obj.Number)

				if err != nil {
					// TODO: Необходимо логировать ошибку
					fmt.Println(err)
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
		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-ctx.Done():
				logger.Log.Info("ЗАКРЫТА - Задача обновления(статусов и начисления) заказов в БД.")
				return
			case obj := <-infoOrdersCh:
				fmt.Println(obj)
				// TODO: Необходимо логировать
				// Тут наверное будет обновдение пачкой в БД
			case <-ticker.C:

				if len(infoOrdersCh) == 0 && len(ordersFromDBCh) == 0 {
					// TODO: Необходимо логировать
					w.getDataDBSemaphore.Release()
				}
			}
		}
	}()
}
