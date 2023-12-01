package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MWT-proger/go-loyalty-system/internal/models"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIHandlerUserRegister(t *testing.T) {
	testCases := []struct {
		name               string
		requestData        string
		requestMethod      string
		requestURL         string
		responseStatusCode int
		responseStoreData  MockResponseStoreUserData
	}{
		{
			name:               "Тест №1 - Всегда успешный",
			requestData:        `{"login": "ivan", "password": "Ivan@asdsl123"}`,
			requestMethod:      http.MethodPost,
			requestURL:         "/api/user/register",
			responseStatusCode: http.StatusOK,
			responseStoreData:  MockResponseStoreUserData{data: nil, err: nil},
		},
		{
			name:               "Тест №2 - Пользователь уже существует",
			requestData:        `{"login": "ivan", "password": "Ivan@asdsl123"}`,
			requestMethod:      http.MethodPost,
			requestURL:         "/api/user/register",
			responseStatusCode: http.StatusConflict,
			responseStoreData:  MockResponseStoreUserData{data: models.NewUser(), err: nil},
		},
		{
			name:               "Тест №3 - UerStore возвращает ошибку",
			requestData:        `{"login": "ivan", "password": "Ivan@asdsl123"}`,
			requestMethod:      http.MethodPost,
			requestURL:         "/api/user/register",
			responseStatusCode: http.StatusInternalServerError,
			responseStoreData:  MockResponseStoreUserData{data: nil, err: errors.New("")},
		},
		{
			name:               "Тест №4 - Не валидное тело запроса",
			requestData:        `{"login": "ivan"}`,
			requestMethod:      http.MethodPost,
			requestURL:         "/api/user/register",
			responseStatusCode: http.StatusBadRequest,
			responseStoreData:  MockResponseStoreUserData{data: nil, err: errors.New("")},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			mockStore := &MockStore{}
			mockUserStore := NewMockUserStore(mockStore, tt.responseStoreData)
			mockOrderstore := NewMockOrderStore(mockStore)
			mockWithdrawalstore := NewMockWithdrawalStore(mockStore)
			mockAccountstore := NewMockAccountStore(mockStore)

			h, err := NewAPIHandler(mockUserStore, mockOrderstore, mockWithdrawalstore, mockAccountstore)

			require.NoError(t, err, "Ошибка при инициализации APIHandler")

			router := chi.NewRouter()
			router.Post("/api/user/register", h.UserRegister)

			ts := httptest.NewServer(router)

			bodyRequest := strings.NewReader(tt.requestData)

			result, _ := testRequest(t, ts, tt.requestMethod, tt.requestURL, bodyRequest)

			defer result.Body.Close()

			assert.Equal(t, tt.responseStatusCode, result.StatusCode, "Код ответа не совпадает с ожидаемым")

			if result.StatusCode == http.StatusOK {
				var isToken bool
				for _, cookie := range result.Cookies() {
					if cookie.Name == "token" {
						assert.NotNil(t, cookie.Value, "Токен авторизации пуст")
						isToken = true
					}

				}

				assert.True(t, isToken, "Токен авторизации не найден в cookies")
			}

		})

	}
}
