package errors

import (
	"net/http"
)

type ErrorNotBonuses struct{}

func (m *ErrorNotBonuses) Error() string {
	return "недостаточно средств для списания"
}

type ErrorAccrualStatusCode500 struct{}

func (m *ErrorAccrualStatusCode500) Error() string {
	return "accrual - внутренняя ошибка сервера"
}

type ErrorAccrualStatusCode429 struct{}

func (m *ErrorAccrualStatusCode429) Error() string {
	return "accrual - превышено количество запросов к сервису"
}

func NewServicesError(text string, httpCode int) *ServicesError {
	return &ServicesError{text, httpCode}
}

type ServicesError struct {
	s        string
	HTTPCode int
}

func (e *ServicesError) Error() string {
	return e.s
}

var GetUserServicesError = NewServicesError(
	"не получилось обработать запрос получения пользователя",
	http.StatusInternalServerError,
)

var UserNotFoundServicesError = NewServicesError(
	"пользователь не авторизован",
	http.StatusUnauthorized,
)

var UserExistsServicesError = NewServicesError(
	"пользователь уже существует",
	http.StatusConflict,
)

var InternalServicesError = NewServicesError(
	"внутренняя ошибка сервера",
	http.StatusInternalServerError,
)

var OrderExistsOtherUserServicesError = NewServicesError(
	"заказ уже загружен другим пользователем",
	http.StatusConflict,
)

var OrderExistsServicesError = NewServicesError(
	"заказ уже был загружен вами",
	http.StatusOK,
)

var ListOrdersEmptyServicesError = NewServicesError(
	"список заказов пуст",
	http.StatusNoContent,
)

var ListWithdrawalsEmptyServicesError = NewServicesError(
	"список списаний пуст",
	http.StatusNoContent,
)

var WithdrawalExistsOtherUserServicesError = NewServicesError(
	"списание уже было произведено другим пользователем",
	http.StatusConflict,
)
var WithdrawalExistsServicesError = NewServicesError(
	"данное вы уже оформили",
	http.StatusOK,
)

var NotBonusesWithdrawalServicesError = NewServicesError(
	"недостаточно средств для списания",
	http.StatusPaymentRequired,
)
