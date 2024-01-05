package handlers

type APIHandler struct {
	UserService       UserServicer
	OrderService      OrderServicer
	WithdrawalService WithdrawalServicer
	AccountService    AccountServicer
}

func NewAPIHandler(
	userService UserServicer,
	orderService OrderServicer,
	withdrawalService WithdrawalServicer,
	accountService AccountServicer,
) (h *APIHandler, err error) {
	hh := &APIHandler{
		UserService:       userService,
		OrderService:      orderService,
		WithdrawalService: withdrawalService,
		AccountService:    accountService,
	}

	return hh, err
}

type BaseBodyDater interface {
	isValid() bool
}
