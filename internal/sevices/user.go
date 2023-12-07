package sevices

import (
	"github.com/MWT-proger/go-loyalty-system/internal/store/userstore"
)

type UserService struct {
	store userstore.UserStorer
}

func NewUserService(s userstore.UserStorer) *UserService {

	return &UserService{
		store: s,
	}
}
