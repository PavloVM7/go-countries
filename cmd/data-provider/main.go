package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"pm.com/go-countries/configs/db"
	app2 "pm.com/go-countries/internal/app"
)

var app = app2.NewApp()

func main() {
	app.Name = "Country data provider"
	var dbConfig db.PostgresConfig
	if err := dbConfig.Read(); err != nil {
		panic(err)
	}

	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\n%s error: %v\n", app.Name, err)
		os.Exit(1)
	}
}

func run() error {
	app.Action = getCountries
	return app.Run()
}
func getCountries(ctx context.Context) error {
	log.Printf("Starting %s..., %v", app.Name, ctx)
	app.Wait()
	return nil
}
