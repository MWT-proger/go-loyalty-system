package handlers

// import (
// 	"context"

// 	"github.com/MWT-proger/go-loyalty-system/internal/models"
// 	"github.com/MWT-proger/go-loyalty-system/internal/store"
// 	"github.com/gofrs/uuid"
// )

// type MockStore struct {
// }

// type MockResponseStoreUserData struct {
// 	data *models.User
// 	err  error
// }

// func (s *MockStore) Init(ctx context.Context) error {
// 	return nil
// }

// type MockUserStore struct {
// 	*MockStore
// 	testCase MockResponseStoreUserData
// }

// func NewMockUserStore(baseStorage *MockStore, testCase MockResponseStoreUserData) *MockUserStore {
// 	return &MockUserStore{baseStorage, testCase}
// }

// func (s *MockUserStore) GetFirstByParameters(ctx context.Context, args map[string]interface{}) (*models.User, error) {
// 	return s.testCase.data, s.testCase.err
// }
// func (s *MockUserStore) Insert(ctx context.Context, obj *models.User) error {
// 	return nil
// }

// type MockOrderStore struct {
// 	*MockStore
// }

// func NewMockOrderStore(baseStorage *MockStore) *MockOrderStore {
// 	return &MockOrderStore{baseStorage}
// }

// func (s *MockOrderStore) GetFirstByParameters(ctx context.Context, args map[string]interface{}) (*models.Order, error) {
// 	return nil, nil
// }
// func (s *MockOrderStore) Insert(ctx context.Context, obj *models.Order) error {
// 	return nil
// }

// func (s *MockOrderStore) GetAllByParameters(ctx context.Context, options *store.OptionsQuery) ([]*models.Order, error) {
// 	return nil, nil
// }
// func (s *MockOrderStore) GetSumByParameters(ctx context.Context, args map[string]interface{}) (int64, error) {
// 	return 0, nil
// }
// func (s *MockOrderStore) UpdateOrderPlusUserAccount(ctx context.Context, options *store.OptionsUpdateQuery, userID uuid.UUID, bonuses int64) error {
// 	return nil
// }
// func (s *MockOrderStore) UpdateBatch(ctx context.Context, options *store.OptionsUpdateQuery) error {
// 	return nil
// }

// type MockWithdrawalStore struct {
// 	*MockStore
// }

// func NewMockWithdrawalStore(baseStorage *MockStore) *MockWithdrawalStore {
// 	return &MockWithdrawalStore{baseStorage}
// }

// func (s *MockWithdrawalStore) Insert(ctx context.Context, obj *models.Withdrawal) error {
// 	return nil
// }
// func (s *MockWithdrawalStore) GetFirstByParameters(ctx context.Context, args map[string]interface{}) (*models.Withdrawal, error) {
// 	return nil, nil
// }
// func (s *MockWithdrawalStore) GetAllByParameters(ctx context.Context, options *store.OptionsQuery) ([]*models.Withdrawal, error) {
// 	return nil, nil
// }
// func (s *MockWithdrawalStore) GetSumByParameters(ctx context.Context, args map[string]interface{}) (int64, error) {
// 	return 0, nil
// }

// type MockAccountStore struct {
// 	*MockStore
// }

// func NewMockAccountStore(baseStorage *MockStore) *MockAccountStore {
// 	return &MockAccountStore{baseStorage}
// }

// func (s *MockAccountStore) Insert(ctx context.Context, obj *models.Account) error {
// 	return nil
// }
// func (s *MockAccountStore) GetFirstByParameters(ctx context.Context, args map[string]interface{}) (*models.Account, error) {
// 	return nil, nil
// }
