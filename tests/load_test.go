package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"wallet-service/internal/handler"
	"wallet-service/internal/repository"
	"wallet-service/internal/service"
	"wallet-service/pkg/testdb"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func TestLoad_Deposit(t *testing.T) {
	testdb.WithDB(t, []string{"../migrations", "../migrations/test"}, func(pool *pgxpool.Pool) {
		repo, err := repository.NewPostgresRepository(pool)
		assert.NoError(t, err)

		services := service.NewService(repo)
		handlers := handler.NewHandler(services)
		router := handlers.GetRouter()

		server := httptest.NewServer(router)
		defer server.Close()

		walletID := testdb.WalletEmptyWalletID
		rateFreq := 1000
		durationS := 10
		amount := int64(1)

		targeter := vegeta.NewStaticTargeter(vegeta.Target{
			Method: "POST",
			URL:    server.URL + "/api/v1/wallet",
			Body:   mustJSON(handler.UpdateWalletRequest{WalletID: walletID, OperationType: "DEPOSIT", Amount: amount}),
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
		})

		attacker := vegeta.NewAttacker()
		rate := vegeta.Rate{Freq: rateFreq, Per: time.Second}
		duration := time.Duration(durationS) * time.Second

		var metrics vegeta.Metrics
		for res := range attacker.Attack(targeter, rate, duration, "Deposit Concurrency Test") {
			metrics.Add(res)
		}
		metrics.Close()

		// Проверка итогового баланса
		resp := mustGetWallet(server, walletID)
		expectedBalance := amount * int64(rateFreq) * int64(durationS)
		assert.Equal(t, expectedBalance, resp.Balance, "incorrect final balance under concurrency")

		// Проверка отсутствия 50x ошибок
		assert.Equal(t, float64(1), metrics.Success, "some requests failed under concurrency")
	})
}

func TestLoad_Withdraw(t *testing.T) {
	testdb.WithDB(t, []string{"../migrations", "../migrations/test"}, func(pool *pgxpool.Pool) {
		repo, err := repository.NewPostgresRepository(pool)
		assert.NoError(t, err)

		services := service.NewService(repo)
		handlers := handler.NewHandler(services)
		router := handlers.GetRouter()

		server := httptest.NewServer(router)
		defer server.Close()

		walletID := testdb.Wallet10000AmountID
		rateFreq := 1000
		durationS := 10
		amount := int64(1)

		targeter := vegeta.NewStaticTargeter(vegeta.Target{
			Method: "POST",
			URL:    server.URL + "/api/v1/wallet",
			Body:   mustJSON(handler.UpdateWalletRequest{WalletID: walletID, OperationType: "WITHDRAW", Amount: amount}),
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
		})

		attacker := vegeta.NewAttacker()
		rate := vegeta.Rate{Freq: rateFreq, Per: time.Second}
		duration := time.Duration(durationS) * time.Second

		var metrics vegeta.Metrics
		for res := range attacker.Attack(targeter, rate, duration, "Withdraw Concurrency Test") {
			metrics.Add(res)
		}
		metrics.Close()

		// Проверка итогового баланса
		resp := mustGetWallet(server, walletID)
		assert.Equal(t, int64(0), resp.Balance, "incorrect final balance under concurrency")

		// Проверка отсутствия 50x ошибок
		assert.Equal(t, float64(1), metrics.Success, "some requests failed under concurrency")
	})
}

func TestLoad_Get(t *testing.T) {
	testdb.WithDB(t, []string{"../migrations", "../migrations/test"}, func(pool *pgxpool.Pool) {
		repo, err := repository.NewPostgresRepository(pool)
		assert.NoError(t, err)

		services := service.NewService(repo)
		handlers := handler.NewHandler(services)
		router := handlers.GetRouter()

		server := httptest.NewServer(router)
		defer server.Close()

		walletID := testdb.Wallet10000AmountID
		rateFreq := 1000
		durationS := 10

		targeter := vegeta.NewStaticTargeter(vegeta.Target{
			Method: "GET",
			URL:    server.URL + "/api/v1/wallets/" + walletID,
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
		})

		attacker := vegeta.NewAttacker()
		rate := vegeta.Rate{Freq: rateFreq, Per: time.Second}
		duration := time.Duration(durationS) * time.Second

		var metrics vegeta.Metrics
		for res := range attacker.Attack(targeter, rate, duration, "Withdraw Concurrency Test") {
			metrics.Add(res)
		}
		metrics.Close()

		// Проверка отсутствия 50x ошибок
		assert.Equal(t, float64(1), metrics.Success, "some requests failed under concurrency")
	})
}

func mustJSON(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func mustGetWallet(server *httptest.Server, walletID string) handler.GetWalletResponse {
	url := server.URL + "/api/v1/wallets/" + walletID
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result handler.GetWalletResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}
	return result
}
