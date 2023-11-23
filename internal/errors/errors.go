package errors

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
