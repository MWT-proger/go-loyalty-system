package services

import (
	"context"

	"github.com/MWT-proger/go-loyalty-system/configs"
	"github.com/MWT-proger/go-loyalty-system/internal/auth"
	lErrors "github.com/MWT-proger/go-loyalty-system/internal/errors"
	"github.com/MWT-proger/go-loyalty-system/internal/models"
)

type UserStorer interface {
	Insert(ctx context.Context, obj *models.User) error
	GetFirstByParameters(ctx context.Context, args map[string]interface{}) (*models.User, error)
}

type UserService struct {
	conf      *configs.Config
	userStore UserStorer
}

func NewUserService(s UserStorer, config *configs.Config) *UserService {

	return &UserService{
		conf:      config,
		userStore: s,
	}
}

// UserLogin(ctx context.Context, login string, password string) проверяет данные пользователя
// при совпадение, формирует токен авторизации и возвращает пользователю
// TODO: МБ и правда генерацию токена оставить в хендлере(для возможности выбора способа авторизации)
func (s *UserService) UserLogin(ctx context.Context, login string, password string) (string, error) {

	args := map[string]interface{}{"login": login}

	user, err := s.userStore.GetFirstByParameters(ctx, args)

	if err != nil {
		return "", lErrors.GetUserServicesError
	}

	if user == nil {
		return "", lErrors.UserNotFoundServicesError
	}

	if ok := auth.CheckPasswordHash(password, user.Password); !ok {
		return "", lErrors.UserNotFoundServicesError
	}

	tokenString, err := auth.BuildJWTString(user.ID, s.conf)

	return tokenString, err
}

// UserRegister(ctx context.Context, login string, password string) создает пользователя
// при успехе, формирует токен авторизации и возвращает пользователю
// TODO: МБ и правда генерацию токена оставить в хендлере(для возможности выбора способа авторизации)
func (s *UserService) UserRegister(ctx context.Context, login string, password string) (string, error) {

	newUser, err := models.NewUser()

	if err != nil {
		return "", lErrors.InternalServicesError
	}

	newUser.Login = login

	args := map[string]interface{}{"login": newUser.Login}
	obj, err := s.userStore.GetFirstByParameters(ctx, args)

	if err != nil {
		return "", lErrors.GetUserServicesError
	}

	if obj != nil {
		return "", lErrors.UserExistsServicesError
	}

	newUser.Password, err = auth.HashPassword(password)

	if err != nil {
		return "", lErrors.InternalServicesError
	}

	err = s.userStore.Insert(ctx, newUser)

	if err != nil {
		return "", lErrors.InternalServicesError
	}

	tokenString, err := auth.BuildJWTString(newUser.ID, s.conf)

	return tokenString, err
}
