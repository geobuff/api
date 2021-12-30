package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/geobuff/api/auth"
	"github.com/geobuff/api/avatars"
	"github.com/geobuff/api/badges"
	"github.com/geobuff/api/checkout"
	"github.com/geobuff/api/dailytrivia"
	"github.com/geobuff/api/discounts"
	"github.com/geobuff/api/leaderboard"
	"github.com/geobuff/api/mappings"
	"github.com/geobuff/api/merch"
	"github.com/geobuff/api/orders"
	"github.com/geobuff/api/plays"
	"github.com/geobuff/api/quizzes"
	"github.com/geobuff/api/repo"
	"github.com/geobuff/api/support"
	"github.com/geobuff/api/tempscores"
	"github.com/geobuff/api/users"
	"github.com/geobuff/api/validation"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {
	err := loadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println("successfully loaded .env config")

	err = repo.OpenConnection()
	if err != nil {
		panic(err)
	}
	fmt.Println("successfully connected to database")

	err = runMigrations()
	if err != nil {
		panic(err)
	}
	fmt.Println("successfully ran database migrations")

	err = validation.Init()
	if err != nil {
		panic(err)
	}
	fmt.Println("successfully initialized validator")

	err = serve()
	if err != nil {
		panic(err)
	}
}

var loadConfig = func() error {
	return godotenv.Load()
}

var runMigrations = func() error {
	driver, err := postgres.WithInstance(repo.Connection, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

var serve = func() error {
	return http.ListenAndServe(":8080", handler(router()))
}

func handler(router http.Handler) http.Handler {
	corsOptions := cors.New(cors.Options{
		AllowedOrigins: strings.Split(os.Getenv("CORS_ORIGINS"), ","),
		AllowedMethods: strings.Split(os.Getenv("CORS_METHODS"), ","),
		AllowedHeaders: strings.Split(os.Getenv("CORS_HEADERS"), ","),
	})

	return corsOptions.Handler(router)
}

func router() http.Handler {
	router := mux.NewRouter()

	// Quiz endpoints.
	router.HandleFunc("/api/quizzes", quizzes.GetQuizzes).Methods("GET")
	router.HandleFunc("/api/quizzes/{id}", quizzes.GetQuiz).Methods("GET")

	// Daily Trivia endpoints.
	router.HandleFunc("/api/daily-trivia", dailytrivia.GenerateDailyTrivia).Methods("POST")
	router.HandleFunc("/api/daily-trivia/{date}", dailytrivia.GetDailyTrivia).Methods("GET")

	// Mapping endpoints.
	router.HandleFunc("/api/mappings/{key}", mappings.GetMapping).Methods("GET")

	// Auth endpoints.
	router.HandleFunc("/api/auth/login", auth.Login).Methods("POST")
	router.HandleFunc("/api/auth/register", auth.Register).Methods("POST")
	router.HandleFunc("/api/auth/send-reset-token", auth.SendResetToken).Methods("POST")
	router.HandleFunc("/api/auth/reset-token-valid/{userId}/{token}", auth.ResetTokenValid).Methods("GET")
	router.HandleFunc("/api/auth", auth.UpdatePasswordUsingToken).Methods("PUT")

	// User endpoints.
	router.HandleFunc("/api/users", users.GetUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", users.GetUser).Methods("GET")
	router.HandleFunc("/api/users/{id}", users.UpdateUser).Methods("PUT")
	router.HandleFunc("/api/users/{id}", users.DeleteUser).Methods("DELETE")

	// Badge endpoints.
	router.HandleFunc("/api/badges", badges.GetBadges).Methods("GET")

	// Temp Score endpoints.
	router.HandleFunc("/api/tempscores/{id}", tempscores.GetTempScore).Methods("GET")
	router.HandleFunc("/api/tempscores", tempscores.CreateTempScore).Methods("POST")

	// Leaderboard endpoints.
	router.HandleFunc("/api/leaderboard/all/{quizId}", leaderboard.GetEntries).Methods("POST")
	router.HandleFunc("/api/leaderboard/{userId}", leaderboard.GetUserEntries).Methods("GET")
	router.HandleFunc("/api/leaderboard/{quizId}/{userId}", leaderboard.GetEntry).Methods("GET")
	router.HandleFunc("/api/leaderboard", leaderboard.CreateEntry).Methods("POST")
	router.HandleFunc("/api/leaderboard/{id}", leaderboard.UpdateEntry).Methods("PUT")
	router.HandleFunc("/api/leaderboard/{id}", leaderboard.DeleteEntry).Methods("DELETE")

	// Play endpoints.
	router.HandleFunc("/api/plays", plays.GetAllPlays).Methods("GET")
	router.HandleFunc("/api/plays/{quizId}", plays.GetPlays).Methods("GET")
	router.HandleFunc("/api/plays/{quizId}", plays.IncrementPlays).Methods("PUT")

	// Checkout endpoints.
	router.HandleFunc("/api/checkout/create-checkout-session", checkout.HandleCreateCheckoutSession).Methods("POST")
	router.HandleFunc("/api/checkout/webhook", checkout.HandleWebhook).Methods("POST")

	// Order endpoints.
	router.HandleFunc("/api/orders/{email}", orders.GetOrders).Methods("GET")
	router.HandleFunc("/api/orders/{email}", orders.CancelOrder).Methods("DELETE")

	// Avatar endpoints.
	router.HandleFunc("/api/avatars", avatars.GetAvatars).Methods("GET")

	// Merch endpoints.
	router.HandleFunc("/api/merch", merch.GetMerch).Methods("GET")

	// Discount endpoints.
	router.HandleFunc("/api/discounts/{code}", discounts.GetDiscount).Methods("GET")

	// Support endpoints.
	router.HandleFunc("/api/support", support.SendSupportRequest).Methods("POST")

	return router
}
