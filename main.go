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
	"github.com/geobuff/api/repo"
	"github.com/geobuff/api/src"
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
	router.HandleFunc("/api/quizzes/all", src.GetQuizzes).Methods("POST")
	router.HandleFunc("/api/quizzes/{id}", src.GetQuiz).Methods("GET")
	router.HandleFunc("/api/quizzes/route/{route}", src.GetQuizByRoute).Methods("GET")
	router.HandleFunc("/api/quizzes", src.CreateQuiz).Methods("POST")
	router.HandleFunc("/api/quizzes/{id}", src.UpdateQuiz).Methods("PUT")
	router.HandleFunc("/api/quizzes/{id}", src.DeleteQuiz).Methods("DELETE")

	// Quiz Type endpoints.
	router.HandleFunc("/api/quiztype", src.GetTypes).Methods("GET")

	// Quiz Plays endpoints.
	router.HandleFunc("/api/quiz-plays", src.GetAllQuizPlays).Methods("GET")
	router.HandleFunc("/api/quiz-plays/{quizId}", src.GetQuizPlays).Methods("GET")
	router.HandleFunc("/api/quiz-plays-top-five", src.GetTopFiveQuizPlays).Methods("GET")
	router.HandleFunc("/api/quiz-plays/{quizId}", src.IncrementQuizPlays).Methods("PUT")

	// Trivia endpoints.
	router.HandleFunc("/api/trivia/all", src.GetAllTrivia).Methods("POST")
	router.HandleFunc("/api/trivia/{date}", src.GetTriviaByDate).Methods("GET")
	router.HandleFunc("/api/trivia", src.GenerateTrivia).Methods("POST")
	router.HandleFunc("/api/trivia/{date}", src.RegenerateTrivia).Methods("PUT")
	router.HandleFunc("/api/trivia/{date}", src.DeleteTrivia).Methods("DELETE")
	router.HandleFunc("/api/trivia/old/{newTriviaCount}", src.DeleteOldTrivia).Methods("DELETE")

	// Trivia Plays endpoints.
	router.HandleFunc("/api/trivia-plays/week", src.GetLastWeekTriviaPlays).Methods("GET")
	router.HandleFunc("/api/trivia-plays/{id}", src.IncrementTriviaPlays).Methods("PUT")

	// Trivia Question Type endpoints.
	router.HandleFunc("/api/trivia-question-types", src.GetTriviaQuestionTypes).Methods("GET")

	// Trivia Question Category endpoints.
	router.HandleFunc("/api/trivia-question-categories", src.GetTriviaQuestionCategories).Methods("GET")

	// Manual Trivia Question endpoints.
	router.HandleFunc("/api/manual-trivia-questions/all", src.GetManualTriviaQuestions).Methods("POST")
	router.HandleFunc("/api/manual-trivia-questions", src.CreateManualTriviaQuestion).Methods("POST")
	router.HandleFunc("/api/manual-trivia-questions/{id}", src.UpdateManualTriviaQuestion).Methods("PUT")
	router.HandleFunc("/api/manual-trivia-questions/{id}", src.DeleteManualTriviaQuestion).Methods("DELETE")

	// Mapping endpoints.
	router.HandleFunc("/api/mappings", src.GetMappingGroups).Methods("GET")
	router.HandleFunc("/api/mappings/{key}", src.GetMappingEntries).Methods("GET")
	router.HandleFunc("/api/mappings-no-flags", src.GetMappingsWithoutFlags).Methods("GET")
	router.HandleFunc("/api/mappings/{key}", src.EditMapping).Methods("PUT")
	router.HandleFunc("/api/mappings/{key}", src.DeleteMapping).Methods("DELETE")

	// Map endpoints.
	router.HandleFunc("/api/maps", src.GetMaps).Methods("GET")
	router.HandleFunc("/api/maps/highlighted/{className}", src.GetMapHighlightedRegions).Methods("GET")
	router.HandleFunc("/api/maps/{className}", src.GetMap).Methods("GET")
	router.HandleFunc("/api/maps/preview", src.GetMapPreview).Methods("POST")
	router.HandleFunc("/api/maps", src.CreateMap).Methods("POST")

	// Flag endpoints.
	router.HandleFunc("/api/flags", src.GetFlagGroups).Methods("GET")
	router.HandleFunc("/api/flags/{key}", src.GetFlagEntries).Methods("GET")
	router.HandleFunc("/api/flags/url/{code}", src.GetFlagUrl).Methods("GET")
	router.HandleFunc("/api/flags", src.CreateFlags).Methods("POST")

	// Continent endpoints.
	router.HandleFunc("/api/continents", src.GetContinents).Methods("GET")

	// Auth endpoints.
	router.HandleFunc("/api/auth/login", src.Login).Methods("POST")
	router.HandleFunc("/api/auth/register", src.Register).Methods("POST")
	router.HandleFunc("/api/auth/refresh", src.RefreshToken).Methods("POST")
	router.HandleFunc("/api/auth/send-reset-token", src.SendResetToken).Methods("POST")
	router.HandleFunc("/api/auth/reset-token-valid/{userId}/{token}", src.ResetTokenValid).Methods("GET")
	router.HandleFunc("/api/auth", src.UpdatePasswordUsingToken).Methods("PUT")
	router.HandleFunc("/api/auth/username/{username}", src.UsernameExists).Methods("GET")
	router.HandleFunc("/api/auth/email/{email}", src.EmailExists).Methods("GET")

	// User endpoints.
	router.HandleFunc("/api/users/all", src.GetUsers).Methods("POST")
	router.HandleFunc("/api/users/{id}", src.GetUser).Methods("GET")
	router.HandleFunc("/api/users/email/{email}", src.GetUserByEmail).Methods("GET")
	router.HandleFunc("/api/users/total/week", src.GetLastWeekTotalUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", src.UpdateUser).Methods("PUT")
	router.HandleFunc("/api/users/xp/{id}", src.UpdateUserXP).Methods("PUT")
	router.HandleFunc("/api/users/{id}", src.DeleteUser).Methods("DELETE")

	// Badge endpoints.
	router.HandleFunc("/api/badges", src.GetBadges).Methods("GET")
	router.HandleFunc("/api/badges/{userId}", src.GetUserBadges).Methods("GET")

	// Temp Score endpoints.
	router.HandleFunc("/api/tempscores/{id}", src.GetTempScore).Methods("GET")
	router.HandleFunc("/api/tempscores", src.CreateTempScore).Methods("POST")

	// Leaderboard endpoints.
	router.HandleFunc("/api/leaderboard/all/{quizId}", src.GetEntries).Methods("POST")
	router.HandleFunc("/api/leaderboard/{userId}", src.GetUserEntries).Methods("GET")
	router.HandleFunc("/api/leaderboard/{quizId}/{userId}", src.GetEntry).Methods("GET")
	router.HandleFunc("/api/leaderboard", src.CreateEntry).Methods("POST")
	router.HandleFunc("/api/leaderboard/{id}", src.UpdateEntry).Methods("PUT")
	router.HandleFunc("/api/leaderboard/{id}", src.DeleteEntry).Methods("DELETE")

	// Shipping option endpoints.
	router.HandleFunc("/api/shipping-options", src.GetShippingOptions).Methods("GET")

	// Checkout endpoints.
	router.HandleFunc("/api/checkout/create-checkout-session", src.HandleCreateCheckoutSession).Methods("POST")
	router.HandleFunc("/api/checkout/webhook", src.HandleWebhook).Methods("POST")

	// Order endpoints.
	router.HandleFunc("/api/orders", src.GetOrders).Methods("POST")
	router.HandleFunc("/api/orders/user/{email}", src.GetUserOrders).Methods("GET")
	router.HandleFunc("/api/orders/status/{id}", src.UpdateOrderStatus).Methods("PUT")
	router.HandleFunc("/api/orders/{id}", src.DeleteOrder).Methods("DELETE")
	router.HandleFunc("/api/orders/email/{email}", src.CancelOrder).Methods("DELETE")

	// Avatar endpoints.
	router.HandleFunc("/api/avatars", src.GetAvatars).Methods("GET")

	// Merch endpoints.
	router.HandleFunc("/api/merch", src.GetMerch).Methods("GET")
	router.HandleFunc("/api/merch/exists", src.MerchExists).Methods("POST")

	// Discount endpoints.
	router.HandleFunc("/api/discounts", src.GetDiscounts).Methods("GET")
	router.HandleFunc("/api/discounts/{code}", src.GetDiscount).Methods("GET")

	// Community Quiz endpoints.
	router.HandleFunc("/api/community-quizzes/all", src.GetCommunityQuizzes).Methods("POST")
	router.HandleFunc("/api/community-quizzes/{id}", src.GetCommunityQuiz).Methods("GET")
	router.HandleFunc("/api/community-quizzes/user/{userId}", src.GetUserCommunityQuizzes).Methods("GET")
	router.HandleFunc("/api/community-quizzes", src.CreateCommunityQuiz).Methods("POST")
	router.HandleFunc("/api/community-quizzes/{id}", src.UpdateCommunityQuiz).Methods("PUT")
	router.HandleFunc("/api/community-quizzes/{id}", src.DeleteCommunityQuiz).Methods("DELETE")

	// Community Quiz Play endpoints.
	router.HandleFunc("/api/community-quiz-plays/{id}", src.IncrementCommunityQuizPlays).Methods("PUT")

	// SEO endpoints.
	router.HandleFunc("/api/seo/dynamic-routes", src.GetDynamicRoutes).Methods("GET")

	return router
}
