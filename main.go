package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	handler "github.com/Quiqui-dev/Earthbound/internal/api/handler/user"
	"github.com/Quiqui-dev/Earthbound/internal/db"
	migration "github.com/Quiqui-dev/Earthbound/internal/db/migrations"
	repo "github.com/Quiqui-dev/Earthbound/internal/repo/user"
	"github.com/Quiqui-dev/Earthbound/internal/router"
	service "github.com/Quiqui-dev/Earthbound/internal/service/user"
	"github.com/Quiqui-dev/Earthbound/util"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env not found, falling back to env variables")
	}

	dbConn, err := db.New()
	if err != nil {
		log.Fatalf("Could not initialise db: %s", err)
	}

	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %s", err)
	}

	log.Println("Connected to database successfully")

	if err := migration.RunMigrations(dbConn); err != nil {
		log.Fatalf("Failed to run migrations", err)
	}

	userRepo := repo.NewUserRepository(dbConn)

	userService := service.NewUserService(userRepo)

	userHandler := handler.NewUserHandler(userService)

	router := router.NewRouter(userHandler)

	done := make(chan bool, 1)
	go gracefulShutdown(router, done)

	log.Printf("server running on port %s", util.GetEnv("PORT", "8080"))
	if err := router.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server", err)
	}

	<-done
	log.Println("Graceful shutdown complete.")
}

func gracefulShutdown(srv *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server forced to shutdown")
	}

	log.Println("Server exiting")

	// notify main goroutine that shutdown is complete
	done <- true
}
