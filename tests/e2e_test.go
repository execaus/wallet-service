package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"wallet-service/internal/handler"
	"wallet-service/internal/repository"
	"wallet-service/internal/service"
	"wallet-service/pkg/testdb"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

var migrationsPath = []string{"../migrations", "../migrations/test"}

func Test_Deposit_Success(t *testing.T) {
	t.Parallel()

	run(t, func(router *gin.Engine) {
		var amount int64 = 100

		resp, _ := request[handler.UpdateWalletResponse](t, router, "POST", "/api/v1/wallet", &handler.UpdateWalletRequest{
			WalletID:      testdb.WalletEmptyWalletID,
			OperationType: "DEPOSIT",
			Amount:        amount,
		}, http.StatusOK)

		assert.Equal(t, amount, resp.NewBalance)
	})
}

func Test_Get_Success(t *testing.T) {
	t.Parallel()

	run(t, func(router *gin.Engine) {
		var amount int64 = 100

		_, _ = request[handler.UpdateWalletResponse](t, router, "POST", "/api/v1/wallet", &handler.UpdateWalletRequest{
			WalletID:      testdb.WalletEmptyWalletID,
			OperationType: "DEPOSIT",
			Amount:        amount,
		}, http.StatusOK)

		resp, _ := request[handler.GetWalletResponse](t, router, "GET", "/api/v1/wallets/"+testdb.WalletEmptyWalletID, nil, http.StatusOK)
		assert.Equal(t, amount, resp.Balance)
	})
}

func Test_Withdraw_Conflict(t *testing.T) {
	t.Parallel()

	run(t, func(router *gin.Engine) {
		var amount int64 = 101

		request[handler.UpdateWalletResponse](t, router, "POST", "/api/v1/wallet", &handler.UpdateWalletRequest{
			WalletID:      testdb.WalletCorrectID,
			OperationType: "WITHDRAW",
			Amount:        amount,
		}, http.StatusConflict)
	})
}

func Test_ConcurrentDeposit_CorrectBalance(t *testing.T) {
	t.Parallel()

	run(t, func(router *gin.Engine) {
		var (
			amount      int64 = 100
			numRoutines       = 10
		)
		var wg sync.WaitGroup

		for i := 0; i < numRoutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, _ = request[handler.UpdateWalletResponse](t, router, "POST", "/api/v1/wallet", &handler.UpdateWalletRequest{
					WalletID:      testdb.WalletEmptyWalletID,
					OperationType: "DEPOSIT",
					Amount:        amount,
				}, http.StatusOK)
			}()
		}
		wg.Wait()

		resp, _ := request[handler.GetWalletResponse](t, router, "GET", "/api/v1/wallets/"+testdb.WalletEmptyWalletID, nil, http.StatusOK)
		assert.Equal(t, amount*int64(numRoutines), resp.Balance)
	})
}

func Test_ConcurrentWithdraw_OneOperationSuccess(t *testing.T) {
	t.Parallel()

	run(t, func(router *gin.Engine) {
		const (
			amount      int64 = 50
			numRoutines       = 10
		)

		var wg sync.WaitGroup
		var statuses [numRoutines]int

		for i := 0; i < numRoutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, status := request[handler.UpdateWalletResponse](t, router, "POST", "/api/v1/wallet", &handler.UpdateWalletRequest{
					WalletID:      testdb.WalletCorrectID,
					OperationType: "WITHDRAW",
					Amount:        amount,
				}, -1)
				statuses[i] = status
			}()
		}
		wg.Wait()

		resp, _ := request[handler.GetWalletResponse](t, router, "GET", "/api/v1/wallets/"+testdb.WalletCorrectID, nil, http.StatusOK)

		okCount := 0
		conflictCount := 0
		for _, s := range statuses {
			if s == http.StatusOK {
				okCount++
			} else if s == http.StatusConflict {
				conflictCount++
			}
		}
		assert.Equal(t, 2, okCount)
		assert.Equal(t, numRoutines-2, conflictCount)

		assert.Equal(t, int64(0), resp.Balance)
	})
}

func run(t *testing.T, fn func(router *gin.Engine)) {
	testdb.WithDB(t, migrationsPath, func(pool *pgxpool.Pool) {
		repo, err := repository.NewPostgresRepository(pool)
		assert.NoError(t, err)

		services := service.NewService(repo)
		handlers := handler.NewHandler(services)

		router := handlers.GetRouter()

		fn(router)
	})
}

func request[ResponseT any](t *testing.T, router *gin.Engine, method string, url string, in interface{}, expectedCode int) (resp ResponseT, code int) {
	testdb.WithDB(t, migrationsPath, func(pool *pgxpool.Pool) {
		body, err := json.Marshal(in)
		assert.NoError(t, err)
		req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if expectedCode != -1 {
			assert.Equal(t, expectedCode, w.Code)
		}

		code = w.Code

		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
	})

	gin.SetMode(gin.TestMode)

	return resp, code
}
