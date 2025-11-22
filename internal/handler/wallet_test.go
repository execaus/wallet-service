package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"wallet-service/internal/domain"
	"wallet-service/internal/service"
	mock_service "wallet-service/internal/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return h.GetRouter()
}

func getBodyReader(t *testing.T, in map[string]interface{}) *bytes.Reader {
	b, err := json.Marshal(in)
	assert.NoError(t, err)
	return bytes.NewReader(b)
}

func TestUpdateWallet_CorrectDeposit_200(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var amount int64 = 1000
	id := uuid.New()
	wallet, err := domain.NewWallet(id, amount)
	assert.NoError(t, err)

	mockWallet := mock_service.NewMockWallet(ctrl)
	mockWallet.
		EXPECT().
		Deposit(gomock.Any(), id, amount).
		Return(wallet, nil)

	srv := service.Service{
		Wallet: mockWallet,
	}
	h := NewHandler(&srv)
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", getBodyReader(t, map[string]interface{}{
		"walletId":      id.String(),
		"operationType": "DEPOSIT",
		"amount":        amount,
	}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateWallet_CorrectWithdraw_200(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var amount int64 = 1000
	id := uuid.New()
	wallet, err := domain.NewWallet(id, amount)
	assert.NoError(t, err)

	mockWallet := mock_service.NewMockWallet(ctrl)
	mockWallet.
		EXPECT().
		Withdraw(gomock.Any(), id, amount).
		Return(wallet, nil)

	srv := service.Service{
		Wallet: mockWallet,
	}
	h := NewHandler(&srv)
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", getBodyReader(t, map[string]interface{}{
		"walletId":      id.String(),
		"operationType": "WITHDRAW",
		"amount":        amount,
	}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateWallet_InvalidWalletId_400(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWallet := mock_service.NewMockWallet(ctrl)

	srv := service.Service{
		Wallet: mockWallet,
	}
	h := NewHandler(&srv)
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", getBodyReader(t, map[string]interface{}{
		"walletId":      "8759432",
		"operationType": "DEPOSIT",
		"amount":        10,
	}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateWallet_InvalidOperation_400(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWallet := mock_service.NewMockWallet(ctrl)

	srv := service.Service{
		Wallet: mockWallet,
	}
	h := NewHandler(&srv)
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", getBodyReader(t, map[string]interface{}{
		"walletId":      uuid.New().String(),
		"operationType": "ADD",
		"amount":        10,
	}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateWallet_ZeroAmount_400(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWallet := mock_service.NewMockWallet(ctrl)

	srv := service.Service{
		Wallet: mockWallet,
	}
	h := NewHandler(&srv)
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", getBodyReader(t, map[string]interface{}{
		"walletId":      uuid.New().String(),
		"operationType": "DEPOSIT",
		"amount":        0,
	}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateWallet_NegativeAmount_400(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWallet := mock_service.NewMockWallet(ctrl)

	srv := service.Service{
		Wallet: mockWallet,
	}
	h := NewHandler(&srv)
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", getBodyReader(t, map[string]interface{}{
		"walletId":      uuid.New().String(),
		"operationType": "DEPOSIT",
		"amount":        -1,
	}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateWallet_ServiceError_500(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var amount int64 = 1000
	id := uuid.New()

	mockWallet := mock_service.NewMockWallet(ctrl)
	mockWallet.
		EXPECT().
		Deposit(gomock.Any(), id, amount).
		Return(nil, errors.New(""))

	srv := service.Service{
		Wallet: mockWallet,
	}
	h := NewHandler(&srv)
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", getBodyReader(t, map[string]interface{}{
		"walletId":      id.String(),
		"operationType": "DEPOSIT",
		"amount":        amount,
	}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateWallet_NonExistentWallet_404(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var amount int64 = 1000
	id := uuid.New()

	mockWallet := mock_service.NewMockWallet(ctrl)
	mockWallet.
		EXPECT().
		Deposit(gomock.Any(), id, amount).
		Return(nil, domain.ErrWalletNotFound)

	srv := service.Service{
		Wallet: mockWallet,
	}
	h := NewHandler(&srv)
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", getBodyReader(t, map[string]interface{}{
		"walletId":      id.String(),
		"operationType": "DEPOSIT",
		"amount":        amount,
	}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateWallet_AmountExceedsBalance_409(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var amount int64 = 1000
	id := uuid.New()

	mockWallet := mock_service.NewMockWallet(ctrl)
	mockWallet.
		EXPECT().
		Deposit(gomock.Any(), id, amount).
		Return(nil, domain.ErrInsufficientBalance)

	srv := service.Service{
		Wallet: mockWallet,
	}
	h := NewHandler(&srv)
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", getBodyReader(t, map[string]interface{}{
		"walletId":      id.String(),
		"operationType": "DEPOSIT",
		"amount":        amount,
	}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestGetWallet_CorrectID_200(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	wallet, err := domain.NewWallet(id, 100)
	assert.NoError(t, err)

	mockWallet := mock_service.NewMockWallet(ctrl)
	mockWallet.
		EXPECT().
		Get(gomock.Any(), id).
		Return(wallet, nil)

	srv := service.Service{
		Wallet: mockWallet,
	}
	h := NewHandler(&srv)
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetWallet_InvalidUUID_400(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWallet := mock_service.NewMockWallet(ctrl)
	srv := service.Service{
		Wallet: mockWallet,
	}
	h := NewHandler(&srv)
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/12345", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetWallet_NonExistentWallet_404(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()

	mockWallet := mock_service.NewMockWallet(ctrl)
	mockWallet.
		EXPECT().
		Get(gomock.Any(), id).
		Return(nil, domain.ErrWalletNotFound)

	srv := service.Service{
		Wallet: mockWallet,
	}
	h := NewHandler(&srv)
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetWallet_ServiceError_500(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()

	mockWallet := mock_service.NewMockWallet(ctrl)
	mockWallet.
		EXPECT().
		Get(gomock.Any(), id).
		Return(nil, errors.New("some error"))

	srv := service.Service{
		Wallet: mockWallet,
	}
	h := NewHandler(&srv)
	router := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/"+id.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
