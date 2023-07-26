package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	crypto_rand "crypto/rand"

	"cloud.google.com/go/errorreporting"
	"github.com/didip/tollbooth"
	"github.com/geobuff/api/api"
	"github.com/geobuff/api/repo"
	"github.com/geobuff/api/types"
	"github.com/geobuff/api/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

var errorClient *errorreporting.Client

func main() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	err = loadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println("successfully loaded .env config")

	environment := os.Getenv("ENV")
	if environment == types.DEV || environment == types.PROD {
		ctx := context.Background()
		errorClient, err = errorreporting.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"), errorreporting.Config{
			OnError: func(err error) {
				log.Printf("Could not log error: %v", err)
			},
		})

		if err != nil {
			log.Fatal(err)
		}
		defer errorClient.Close()
	}

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

	err = utils.Init()
	if err != nil {
		panic(err)
	}
	fmt.Println("successfully initialized validator")

	utils.InitCache()
	fmt.Println("successfully initialized translation cache")

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
	max, _ := strconv.ParseFloat(os.Getenv("RATE_LIMITER_MAX"), 64)
	return http.ListenAndServe(":8080", tollbooth.LimitHandler(tollbooth.NewLimiter(max, nil), (handler(router()))))
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
	router.HandleFunc("/api/quizzes/all", api.GetQuizzes).Methods("POST")
	router.HandleFunc("/api/quizzes/{id}", api.GetQuiz).Methods("GET")
	router.HandleFunc("/api/quizzes/route/{route}", api.GetQuizByRoute).Methods("GET")
	router.HandleFunc("/api/quizzes", api.CreateQuiz).Methods("POST")
	router.HandleFunc("/api/quizzes/{id}", api.UpdateQuiz).Methods("PUT")
	router.HandleFunc("/api/quizzes/{id}", api.DeleteQuiz).Methods("DELETE")

	// Quiz Type endpoints.
	router.HandleFunc("/api/quiztype", api.GetTypes).Methods("GET")

	// Quiz Plays endpoints.
	router.HandleFunc("/api/quiz-plays", api.GetAllQuizPlays).Methods("GET")
	router.HandleFunc("/api/quiz-plays/{quizId}", api.GetQuizPlays).Methods("GET")
	router.HandleFunc("/api/quiz-plays-top-five", api.GetTopFiveQuizPlays).Methods("GET")
	router.HandleFunc("/api/quiz-plays/{quizId}", api.IncrementQuizPlays).Methods("PUT")

	// Trivia endpoints.
	router.HandleFunc("/api/trivia/all", api.GetAllTrivia).Methods("POST")
	router.HandleFunc("/api/trivia/{date}", api.GetTriviaByDate).Methods("GET")
	router.HandleFunc("/api/trivia", api.GenerateTrivia).Methods("POST")
	router.HandleFunc("/api/trivia/{date}", api.RegenerateTrivia).Methods("PUT")
	router.HandleFunc("/api/trivia/{date}", api.DeleteTrivia).Methods("DELETE")
	router.HandleFunc("/api/trivia/old/{newTriviaCount}", api.DeleteOldTrivia).Methods("DELETE")

	// Trivia Plays endpoints.
	router.HandleFunc("/api/trivia-plays/week", api.GetLastWeekTriviaPlays).Methods("GET")
	router.HandleFunc("/api/trivia-plays/{id}", api.IncrementTriviaPlays).Methods("PUT")

	// Trivia Question Type endpoints.
	router.HandleFunc("/api/trivia-question-types", api.GetTriviaQuestionTypes).Methods("GET")

	// Trivia Question Category endpoints.
	router.HandleFunc("/api/trivia-question-categories", api.GetTriviaQuestionCategories).Methods("GET")

	// Manual Trivia Question endpoints.
	router.HandleFunc("/api/manual-trivia-questions/all", api.GetManualTriviaQuestions).Methods("POST")
	router.HandleFunc("/api/manual-trivia-questions", api.CreateManualTriviaQuestion).Methods("POST")
	router.HandleFunc("/api/manual-trivia-questions/{id}", api.UpdateManualTriviaQuestion).Methods("PUT")
	router.HandleFunc("/api/manual-trivia-questions/{id}", api.DeleteManualTriviaQuestion).Methods("DELETE")

	// Mapping endpoints.
	router.HandleFunc("/api/mappings", api.GetMappingGroups).Methods("GET")
	router.HandleFunc("/api/mappings/{key}", api.GetMappingEntries).Methods("GET")
	router.HandleFunc("/api/mappings-no-flags", api.GetMappingsWithoutFlags).Methods("GET")
	router.HandleFunc("/api/mappings/{key}", api.EditMapping).Methods("PUT")
	router.HandleFunc("/api/mappings/{key}", api.DeleteMapping).Methods("DELETE")

	// Map endpoints.
	router.HandleFunc("/api/maps", api.GetMaps).Methods("GET")
	router.HandleFunc("/api/maps/highlighted/{className}", api.GetMapHighlightedRegions).Methods("GET")
	router.HandleFunc("/api/maps/{className}", api.GetMap).Methods("GET")
	router.HandleFunc("/api/maps/preview", api.GetMapPreview).Methods("POST")
	router.HandleFunc("/api/maps", api.CreateMap).Methods("POST")

	// Flag endpoints.
	router.HandleFunc("/api/flags", api.GetFlagGroups).Methods("GET")
	router.HandleFunc("/api/flags/{key}", api.GetFlagEntries).Methods("GET")
	router.HandleFunc("/api/flags/url/{code}", api.GetFlagUrl).Methods("GET")
	router.HandleFunc("/api/flags", api.CreateFlags).Methods("POST")

	// Continent endpoints.
	router.HandleFunc("/api/continents", api.GetContinents).Methods("GET")

	// Auth endpoints.
	router.HandleFunc("/api/auth/login", api.Login).Methods("POST")
	router.HandleFunc("/api/auth/register", api.Register).Methods("POST")
	router.HandleFunc("/api/auth/refresh", api.RefreshToken).Methods("POST")
	router.HandleFunc("/api/auth/send-reset-token", api.SendResetToken).Methods("POST")
	router.HandleFunc("/api/auth/reset-token-valid/{userId}/{token}", api.ResetTokenValid).Methods("GET")
	router.HandleFunc("/api/auth", api.UpdatePasswordUsingToken).Methods("PUT")
	router.HandleFunc("/api/auth/username/{username}", api.UsernameExists).Methods("GET")
	router.HandleFunc("/api/auth/email/{email}", api.EmailExists).Methods("GET")

	// User endpoints.
	router.HandleFunc("/api/users/all", api.GetUsers).Methods("POST")
	router.HandleFunc("/api/users/{id}", api.GetUser).Methods("GET")
	router.HandleFunc("/api/users/email/{email}", api.GetUserByEmail).Methods("GET")
	router.HandleFunc("/api/users/total/week", api.GetLastWeekTotalUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", api.UpdateUser).Methods("PUT")
	router.HandleFunc("/api/users/xp/{id}", api.UpdateUserXP).Methods("PUT")
	router.HandleFunc("/api/users/{id}", api.DeleteUser).Methods("DELETE")

	// Badge endpoints.
	router.HandleFunc("/api/badges", api.GetBadges).Methods("GET")
	router.HandleFunc("/api/badges/{userId}", api.GetUserBadges).Methods("GET")

	// Temp Score endpoints.
	router.HandleFunc("/api/tempscores/{id}", api.GetTempScore).Methods("GET")
	router.HandleFunc("/api/tempscores", api.CreateTempScore).Methods("POST")

	// Leaderboard endpoints.
	router.HandleFunc("/api/leaderboard/all/{quizId}", api.GetEntries).Methods("POST")
	router.HandleFunc("/api/leaderboard/{userId}", api.GetUserEntries).Methods("GET")
	router.HandleFunc("/api/leaderboard/{quizId}/{userId}", api.GetEntry).Methods("GET")
	router.HandleFunc("/api/leaderboard", api.CreateEntry).Methods("POST")
	router.HandleFunc("/api/leaderboard/{id}", api.UpdateEntry).Methods("PUT")
	router.HandleFunc("/api/leaderboard/{id}", api.DeleteEntry).Methods("DELETE")

	// Shipping option endpoints.
	router.HandleFunc("/api/shipping-options", api.GetShippingOptions).Methods("GET")

	// Checkout endpoints.
	router.HandleFunc("/api/checkout/create-checkout-session", api.HandleCreateCheckoutSession).Methods("POST")
	router.HandleFunc("/api/checkout/webhook", api.HandleWebhook).Methods("POST")

	// Order endpoints.
	router.HandleFunc("/api/orders", api.GetOrders).Methods("POST")
	router.HandleFunc("/api/orders/user/{email}", api.GetUserOrders).Methods("GET")
	router.HandleFunc("/api/orders/status/{id}", api.UpdateOrderStatus).Methods("PUT")
	router.HandleFunc("/api/orders/{id}", api.DeleteOrder).Methods("DELETE")
	router.HandleFunc("/api/orders/email/{email}", api.CancelOrder).Methods("DELETE")

	// Avatar endpoints.
	router.HandleFunc("/api/avatars", api.GetAvatars).Methods("GET")

	// Merch endpoints.
	router.HandleFunc("/api/merch", api.GetMerch).Methods("GET")
	router.HandleFunc("/api/merch/exists", api.MerchExists).Methods("POST")

	// Discount endpoints.
	router.HandleFunc("/api/discounts", api.GetDiscounts).Methods("GET")
	router.HandleFunc("/api/discounts/{code}", api.GetDiscount).Methods("GET")

	// Community Quiz endpoints.
	router.HandleFunc("/api/community-quizzes/all", api.GetCommunityQuizzes).Methods("POST")
	router.HandleFunc("/api/community-quizzes/{id}", api.GetCommunityQuiz).Methods("GET")
	router.HandleFunc("/api/community-quizzes/user/{userId}", api.GetUserCommunityQuizzes).Methods("GET")
	router.HandleFunc("/api/community-quizzes", api.CreateCommunityQuiz).Methods("POST")
	router.HandleFunc("/api/community-quizzes/{id}", api.UpdateCommunityQuiz).Methods("PUT")
	router.HandleFunc("/api/community-quizzes/{id}", api.DeleteCommunityQuiz).Methods("DELETE")

	// Community Quiz Play endpoints.
	router.HandleFunc("/api/community-quiz-plays/{id}", api.IncrementCommunityQuizPlays).Methods("PUT")

	// SEO endpoints.
	router.HandleFunc("/api/seo/dynamic-routes", api.GetDynamicRoutes).Methods("GET")

	return router
}
