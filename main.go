package main

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wallet-service/config"
	"wallet-service/internal/handler"
	"wallet-service/internal/repository"
	"wallet-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations
var embedMigrations embed.FS

func main() {
	gin.SetMode(gin.ReleaseMode)

	cfg := config.LoadConfig()

	ctx := context.Background()
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	runMigrations(dsn, cfg.Database.Test)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal("error to open connect to database")
	}

	repositories, err := repository.NewPostgresRepository(pool)
	if err != nil {
		log.Fatal("error to open connect to database")
	}

	services := service.NewService(repositories)
	handlers := handler.NewHandler(services)

	router := handlers.GetRouter()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	go func() {
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Could not listen on %s: %v\n", cfg.Server.Port, err)
		}
	}()

	log.Printf("Server is running on port %s\n", cfg.Server.Port)

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	if pool != nil {
		pool.Close()
	}

	log.Println("Server gracefully stopped")
}

func runMigrations(dsn string, withTestData bool) {
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	goose.SetBaseFS(embedMigrations)
	if err = goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}
	if err = goose.Up(sqlDB, "migrations"); err != nil {
		log.Fatal(err)
	}

	if withTestData {
		if err = goose.Up(sqlDB, "migrations/test"); err != nil {
			log.Fatal(err)
		}
	}

	if err = sqlDB.Close(); err != nil {
		log.Println("failed to close sqlDB:", err)
	}
}
