package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kokhno-nikolay/news/config"
	"github.com/kokhno-nikolay/news/internal/repository"
	"github.com/kokhno-nikolay/news/internal/repository/postgresql"
	"github.com/kokhno-nikolay/news/internal/server"
	"github.com/kokhno-nikolay/news/internal/service"
)

const (
	// тільки для зручності (щоб не парився з конфігом)
	// у реальному проекті, звичайно ж так не роблю :)
	dns = "user=postgres password=postgres dbname=boosters sslmode=disable host=postgres"
)

//  @title          News service
//  @version        1.0
//  @description    Test task

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
	services := service.NewService(repos)
	server := server.NewServer(services.PostService)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	// starting grpc server
	go func() {
		defer wg.Done()

		if err := server.StartGrpcServer(cfg); err != nil {
			log.Fatal(err)
		}
	}()

	// starting http server
	go func() {
		defer wg.Done()

		if err := server.StartHttpServer(ctx, cfg); err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()
}
