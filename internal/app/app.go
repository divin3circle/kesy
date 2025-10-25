package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/nhx-finance/kesy/internal/middleware"
	store "github.com/nhx-finance/kesy/internal/stores"
)

/** Resources needed by the application
* 1. Postgres Database
* 2. Stores: token, user, wallet, mints, kyc, deposits,
* 3. Handlers for each store
* 4. Middleware
* 5. Logger
* 6. Port
 */

 type Application struct {
	Logger *log.Logger
	Port int
	DB *sql.DB
	Middleware *middleware.UserAuthMiddleware
 }

 func loadEnvironmentVariables() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("No .env file found (using environment variables from system)")
	} else {
		fmt.Println("Environment variables loaded from .env file")
	}
} 

 func NewApplication() (*Application, error) {
	loadEnvironmentVariables()

	strport := os.Getenv("PORT")
	port, err := strconv.Atoi(strport)
	if err != nil {
		fmt.Printf("Failed to parse PORT env variable: %v\n", err)
		fmt.Println("Will use default PORT if --port isn't used")
		port = 8080
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	db, err := store.Open()
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}

	return &Application{
		Logger: logger,
		Port:   port,
		DB:     db,
	}, nil
 }

 func (a *Application) HandleStatus(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Healthy!")
 }