package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/nhx-finance/kesy/internal/app"
	"github.com/nhx-finance/kesy/internal/routes"
)

func main() {
	kesy, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
			panic(err)
		}
	}(kesy.DB)

	var port int
	flag.IntVar(&port, "port", 8080, "application port")
	flag.Parse()

	kesy.Logger.Println("Application running")

	r := routes.SetUpRoutes(kesy)

	server := &http.Server{
		Addr:               fmt.Sprintf(":%d", port),
		Handler:        	r,
		IdleTimeout:        time.Minute,
		ReadHeaderTimeout:  time.Second * 10,
		WriteTimeout:       time.Second * 30,
	}

	kesy.Logger.Printf("Starting server on port %d", port)

	err = server.ListenAndServe()
	if err != nil {
		kesy.Logger.Fatalf("Error starting server: %v", err)
	}
}