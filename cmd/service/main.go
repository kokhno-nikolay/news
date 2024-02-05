package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kokhno-nikolay/news/config"
	"github.com/kokhno-nikolay/news/internal/handler"
	"github.com/kokhno-nikolay/news/internal/repository"
	"github.com/kokhno-nikolay/news/internal/repository/postgresql"
	"github.com/kokhno-nikolay/news/internal/server"
)

const (
	// тільки для зручності (щоб не парився з конфігом)
	// у реальному проекті, звичайно ж так не роблю :)
	dns = "user=postgres password=postgres dbname=boosters sslmode=disable host=postgres"
)

//  @title          News service
//  @version        1.0
//  @description    Test task for Promova

// @host            localhost:8000
// @BasePath        /
func main() {
	ctx := context.Background()
	cfg := config.GetConfig()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	db, err := postgresql.NewClient(dns)
	if err != nil {
		panic(err.Error())
	}

	repos := repository.NewRepository(db)
	handlers := handler.NewHandler(repos)

	srv := server.NewServer(cfg, handlers.InitRoutes())
	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Printf("error occurred while running http server: %s\n", err.Error())
		}
	}()
	log.Printf("Server started")

	<-quit

	if err := srv.Stop(ctx); err != nil {
		log.Printf("failed to stop server: %v", err.Error())
	}

}
