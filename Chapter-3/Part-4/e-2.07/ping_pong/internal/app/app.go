package app

import (
	"common/db"
	"context"
	"fmt"
	"sync"

	handler "ping_pong/internal/api"
	"ping_pong/internal/migrations"
	"ping_pong/internal/store"
)

type Application struct {
	PingpongHandler *handler.PingPongHandler
	wg              sync.WaitGroup
}

func NewApplication() (*Application, error) {
	postgresDB, err := db.Open()
	if err != nil {
		panic(err)
	}

	err = db.MigrateFS(postgresDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	pingpongRepo := store.NewPingPongStore(postgresDB)
	pingpongHandler := handler.NewPingPongHandler(pingpongRepo)

	app := &Application{
		PingpongHandler: pingpongHandler,
	}

	return app, nil
}

func (a *Application) Start(ctx context.Context) error {
	fmt.Println("Starting application services....")

	a.wg.Add(1)

	return nil
}

func (a *Application) Stop() error {
	fmt.Println("Stopping application services...")

	a.wg.Wait()

	fmt.Println("Application stopped succesfully")
	return nil
}
