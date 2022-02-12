package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/didip/tollbooth"
	"github.com/geobuff/api/auth"
	"github.com/geobuff/api/avatars"
	"github.com/geobuff/api/badges"
	"github.com/geobuff/api/checkout"
	"github.com/geobuff/api/continents"
	"github.com/geobuff/api/discounts"
	"github.com/geobuff/api/leaderboard"
	"github.com/geobuff/api/manualtriviaquestions"
	"github.com/geobuff/api/mappings"
	"github.com/geobuff/api/merch"
	"github.com/geobuff/api/orders"
	"github.com/geobuff/api/quizplays"
	"github.com/geobuff/api/quiztype"
	"github.com/geobuff/api/quizzes"
	"github.com/geobuff/api/repo"
	"github.com/geobuff/api/tempscores"
	"github.com/geobuff/api/trivia"
	"github.com/geobuff/api/triviaplays"
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
	rand.Seed(time.Now().UnixNano())
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
	return http.ListenAndServe(":8080", tollbooth.LimitHandler(tollbooth.NewLimiter(10, nil), (handler(router()))))
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
	router.HandleFunc("/api/quizzes/all", quizzes.GetQuizzes).Methods("POST")
	router.HandleFunc("/api/quizzes/{id}", quizzes.GetQuiz).Methods("GET")
	router.HandleFunc("/api/quizzes", quizzes.CreateQuiz).Methods("POST")
	router.HandleFunc("/api/quizzes/enabled/{id}", quizzes.ToggleQuizEnabled).Methods("PUT")

	// Quiz Type endpoints.
	router.HandleFunc("/api/quiztype", quiztype.GetTypes).Methods("GET")

	// Quiz Plays endpoints.
	router.HandleFunc("/api/quiz-plays", quizplays.GetAllQuizPlays).Methods("GET")
	router.HandleFunc("/api/quiz-plays/{quizId}", quizplays.GetQuizPlays).Methods("GET")
	router.HandleFunc("/api/quiz-plays-top-five", quizplays.GetTopFiveQuizPlays).Methods("GET")
	router.HandleFunc("/api/quiz-plays/{quizId}", quizplays.IncrementQuizPlays).Methods("PUT")

	// Trivia endpoints.
	router.HandleFunc("/api/trivia", trivia.GetAllTrivia).Methods("GET")
	router.HandleFunc("/api/trivia/{date}", trivia.GetTriviaByDate).Methods("GET")
	router.HandleFunc("/api/trivia", trivia.GenerateTrivia).Methods("POST")

	// Trivia Plays endpoints.
	router.HandleFunc("/api/trivia-plays/week", triviaplays.GetLastWeekTriviaPlays).Methods("GET")
	router.HandleFunc("/api/trivia-plays/{id}", triviaplays.IncrementTriviaPlays).Methods("PUT")

	// Manual Trivia Question endpoints.
	router.HandleFunc("/api/manual-trivia-questions", manualtriviaquestions.GetManualTriviaQuestions).Methods("GET")
	router.HandleFunc("/api/manual-trivia-questions/{id}", manualtriviaquestions.DeleteManualTriviaQuestion).Methods("DELETE")

	// Mapping endpoints.
	router.HandleFunc("/api/mappings/{key}", mappings.GetMapping).Methods("GET")

	// Continent endpoints.
	router.HandleFunc("/api/continents", continents.GetContinents).Methods("GET")

	// Auth endpoints.
	router.HandleFunc("/api/auth/login", auth.Login).Methods("POST")
	router.HandleFunc("/api/auth/register", auth.Register).Methods("POST")
	router.HandleFunc("/api/auth/send-reset-token", auth.SendResetToken).Methods("POST")
	router.HandleFunc("/api/auth/reset-token-valid/{userId}/{token}", auth.ResetTokenValid).Methods("GET")
	router.HandleFunc("/api/auth", auth.UpdatePasswordUsingToken).Methods("PUT")
	router.HandleFunc("/api/auth/username/{username}", auth.UsernameExists).Methods("GET")
	router.HandleFunc("/api/auth/email/{email}", auth.EmailExists).Methods("GET")

	// User endpoints.
	router.HandleFunc("/api/users", users.GetUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", users.GetUser).Methods("GET")
	router.HandleFunc("/api/users/total/week", users.GetLastWeekTotalUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", users.UpdateUser).Methods("PUT")
	router.HandleFunc("/api/users/xp/{id}", users.UpdateUserXP).Methods("PUT")
	router.HandleFunc("/api/users/{id}", users.DeleteUser).Methods("DELETE")

	// Badge endpoints.
	router.HandleFunc("/api/badges", badges.GetBadges).Methods("GET")
	router.HandleFunc("/api/badges/{userId}", badges.GetUserBadges).Methods("GET")

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

	// Checkout endpoints.
	router.HandleFunc("/api/checkout/create-checkout-session", checkout.HandleCreateCheckoutSession).Methods("POST")
	router.HandleFunc("/api/checkout/webhook", checkout.HandleWebhook).Methods("POST")

	// Order endpoints.
	router.HandleFunc("/api/orders", orders.GetOrders).Methods("POST")
	router.HandleFunc("/api/orders/user/{email}", orders.GetUserOrders).Methods("GET")
	router.HandleFunc("/api/orders/status/{id}", orders.UpdateOrderStatus).Methods("PUT")
	router.HandleFunc("/api/orders/{id}", orders.DeleteOrder).Methods("DELETE")
	router.HandleFunc("/api/orders/email/{email}", orders.CancelOrder).Methods("DELETE")

	// Avatar endpoints.
	router.HandleFunc("/api/avatars", avatars.GetAvatars).Methods("GET")

	// Merch endpoints.
	router.HandleFunc("/api/merch", merch.GetMerch).Methods("GET")

	// Discount endpoints.
	router.HandleFunc("/api/discounts", discounts.GetDiscounts).Methods("GET")
	router.HandleFunc("/api/discounts/{code}", discounts.GetDiscount).Methods("GET")

	return router
}
