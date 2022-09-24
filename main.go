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
	"github.com/geobuff/api/auth"
	"github.com/geobuff/api/avatars"
	"github.com/geobuff/api/badges"
	"github.com/geobuff/api/checkout"
	"github.com/geobuff/api/communityquizplays"
	"github.com/geobuff/api/communityquizzes"
	"github.com/geobuff/api/continents"
	"github.com/geobuff/api/discounts"
	"github.com/geobuff/api/flags"
	"github.com/geobuff/api/leaderboard"
	"github.com/geobuff/api/manualtriviaquestions"
	"github.com/geobuff/api/mappings"
	"github.com/geobuff/api/maps"
	"github.com/geobuff/api/merch"
	"github.com/geobuff/api/orders"
	"github.com/geobuff/api/quizplays"
	"github.com/geobuff/api/quiztype"
	"github.com/geobuff/api/quizzes"
	"github.com/geobuff/api/repo"
	"github.com/geobuff/api/seo"
	"github.com/geobuff/api/shippingoptions"
	"github.com/geobuff/api/tempscores"
	"github.com/geobuff/api/trivia"
	"github.com/geobuff/api/triviaplays"
	"github.com/geobuff/api/triviaquestioncategory"
	"github.com/geobuff/api/triviaquestiontype"
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
	if environment == DEV || environment == PROD {
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
	router.HandleFunc("/api/quizzes/all", quizzes.GetQuizzes).Methods("POST")
	router.HandleFunc("/api/quizzes/{id}", quizzes.GetQuiz).Methods("GET")
	router.HandleFunc("/api/quizzes/route/{route}", quizzes.GetQuizByRoute).Methods("GET")
	router.HandleFunc("/api/quizzes", quizzes.CreateQuiz).Methods("POST")
	router.HandleFunc("/api/quizzes/{id}", quizzes.UpdateQuiz).Methods("PUT")
	router.HandleFunc("/api/quizzes/{id}", quizzes.DeleteQuiz).Methods("DELETE")

	// Quiz Type endpoints.
	router.HandleFunc("/api/quiztype", quiztype.GetTypes).Methods("GET")

	// Quiz Plays endpoints.
	router.HandleFunc("/api/quiz-plays", quizplays.GetAllQuizPlays).Methods("GET")
	router.HandleFunc("/api/quiz-plays/{quizId}", quizplays.GetQuizPlays).Methods("GET")
	router.HandleFunc("/api/quiz-plays-top-five", quizplays.GetTopFiveQuizPlays).Methods("GET")
	router.HandleFunc("/api/quiz-plays/{quizId}", quizplays.IncrementQuizPlays).Methods("PUT")

	// Trivia endpoints.
	router.HandleFunc("/api/trivia/all", trivia.GetAllTrivia).Methods("POST")
	router.HandleFunc("/api/trivia/{date}", trivia.GetTriviaByDate).Methods("GET")
	router.HandleFunc("/api/trivia", trivia.GenerateTrivia).Methods("POST")
	router.HandleFunc("/api/trivia/{date}", trivia.RegenerateTrivia).Methods("PUT")
	router.HandleFunc("/api/trivia/{date}", trivia.DeleteTrivia).Methods("DELETE")
	router.HandleFunc("/api/trivia/old/{newTriviaCount}", trivia.DeleteOldTrivia).Methods("DELETE")

	// Trivia Plays endpoints.
	router.HandleFunc("/api/trivia-plays/week", triviaplays.GetLastWeekTriviaPlays).Methods("GET")
	router.HandleFunc("/api/trivia-plays/{id}", triviaplays.IncrementTriviaPlays).Methods("PUT")

	// Trivia Question Type endpoints.
	router.HandleFunc("/api/trivia-question-types", triviaquestiontype.GetTriviaQuestionTypes).Methods("GET")

	// Trivia Question Category endpoints.
	router.HandleFunc("/api/trivia-question-categories", triviaquestioncategory.GetTriviaQuestionCategories).Methods("GET")

	// Manual Trivia Question endpoints.
	router.HandleFunc("/api/manual-trivia-questions/all", manualtriviaquestions.GetManualTriviaQuestions).Methods("POST")
	router.HandleFunc("/api/manual-trivia-questions", manualtriviaquestions.CreateManualTriviaQuestion).Methods("POST")
	router.HandleFunc("/api/manual-trivia-questions/{id}", manualtriviaquestions.UpdateManualTriviaQuestion).Methods("PUT")
	router.HandleFunc("/api/manual-trivia-questions/{id}", manualtriviaquestions.DeleteManualTriviaQuestion).Methods("DELETE")

	// Mapping endpoints.
	router.HandleFunc("/api/mappings", mappings.GetMappingGroups).Methods("GET")
	router.HandleFunc("/api/mappings/{key}", mappings.GetMappingEntries).Methods("GET")
	router.HandleFunc("/api/mappings-no-flags", mappings.GetMappingsWithoutFlags).Methods("GET")
	router.HandleFunc("/api/mappings/{key}", mappings.EditMapping).Methods("PUT")
	router.HandleFunc("/api/mappings/{key}", mappings.DeleteMapping).Methods("DELETE")

	// Map endpoints.
	router.HandleFunc("/api/maps", maps.GetMaps).Methods("GET")
	router.HandleFunc("/api/maps/highlighted/{className}", maps.GetMapHighlightedRegions).Methods("GET")
	router.HandleFunc("/api/maps/{className}", maps.GetMap).Methods("GET")
	router.HandleFunc("/api/maps/preview", maps.GetMapPreview).Methods("POST")
	router.HandleFunc("/api/maps", maps.CreateMap).Methods("POST")

	// Flag endpoints.
	router.HandleFunc("/api/flags", flags.GetFlagGroups).Methods("GET")
	router.HandleFunc("/api/flags/{key}", flags.GetFlagEntries).Methods("GET")
	router.HandleFunc("/api/flags/url/{code}", flags.GetFlagUrl).Methods("GET")
	router.HandleFunc("/api/flags", flags.CreateFlags).Methods("POST")

	// Continent endpoints.
	router.HandleFunc("/api/continents", continents.GetContinents).Methods("GET")

	// Auth endpoints.
	router.HandleFunc("/api/auth/login", auth.Login).Methods("POST")
	router.HandleFunc("/api/auth/register", auth.Register).Methods("POST")
	router.HandleFunc("/api/auth/refresh", auth.RefreshToken).Methods("POST")
	router.HandleFunc("/api/auth/send-reset-token", auth.SendResetToken).Methods("POST")
	router.HandleFunc("/api/auth/reset-token-valid/{userId}/{token}", auth.ResetTokenValid).Methods("GET")
	router.HandleFunc("/api/auth", auth.UpdatePasswordUsingToken).Methods("PUT")
	router.HandleFunc("/api/auth/username/{username}", auth.UsernameExists).Methods("GET")
	router.HandleFunc("/api/auth/email/{email}", auth.EmailExists).Methods("GET")

	// User endpoints.
	router.HandleFunc("/api/users/all", users.GetUsers).Methods("POST")
	router.HandleFunc("/api/users/{id}", users.GetUser).Methods("GET")
	router.HandleFunc("/api/users/email/{email}", users.GetUserByEmail).Methods("GET")
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

	// Shipping option endpoints.
	router.HandleFunc("/api/shipping-options", shippingoptions.GetShippingOptions).Methods("GET")

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
	router.HandleFunc("/api/merch/exists", merch.MerchExists).Methods("POST")

	// Discount endpoints.
	router.HandleFunc("/api/discounts", discounts.GetDiscounts).Methods("GET")
	router.HandleFunc("/api/discounts/{code}", discounts.GetDiscount).Methods("GET")

	// Community Quiz endpoints.
	router.HandleFunc("/api/community-quizzes/all", communityquizzes.GetCommunityQuizzes).Methods("POST")
	router.HandleFunc("/api/community-quizzes/{id}", communityquizzes.GetCommunityQuiz).Methods("GET")
	router.HandleFunc("/api/community-quizzes/user/{userId}", communityquizzes.GetUserCommunityQuizzes).Methods("GET")
	router.HandleFunc("/api/community-quizzes", communityquizzes.CreateCommunityQuiz).Methods("POST")
	router.HandleFunc("/api/community-quizzes/{id}", communityquizzes.UpdateCommunityQuiz).Methods("PUT")
	router.HandleFunc("/api/community-quizzes/{id}", communityquizzes.DeleteCommunityQuiz).Methods("DELETE")

	// Community Quiz Play endpoints.
	router.HandleFunc("/api/community-quiz-plays/{id}", communityquizplays.IncrementCommunityQuizPlays).Methods("PUT")

	// SEO endpoints.
	router.HandleFunc("/api/seo/dynamic-routes", seo.GetDynamicRoutes).Methods("GET")

	return router
}
