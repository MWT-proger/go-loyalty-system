package errors

type ErrorNotBonuses struct{}

func (m *ErrorNotBonuses) Error() string {
	return "недостаточно средств для списания"
}
