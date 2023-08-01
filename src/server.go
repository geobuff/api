package src

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/didip/tollbooth"
	"github.com/geobuff/api/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	ts utils.ITranslationService
	es utils.IEmailService
	vs utils.IValidationService
}

func NewServer(ts utils.ITranslationService, es utils.IEmailService, vs utils.IValidationService) *Server {
	return &Server{
		ts,
		es,
		vs,
	}
}

func getMockServer() *Server {
	return NewServer(utils.NewTranslationService(), utils.NewEmailService(), utils.NewValidationService())
}

func (s *Server) Start() error {
	max, _ := strconv.ParseFloat(os.Getenv("RATE_LIMITER_MAX"), 64)
	return http.ListenAndServe(":8080", tollbooth.LimitHandler(tollbooth.NewLimiter(max, nil), (handler(s.router()))))
}

func handler(router http.Handler) http.Handler {
	corsOptions := cors.New(cors.Options{
		AllowedOrigins: strings.Split(os.Getenv("CORS_ORIGINS"), ","),
		AllowedMethods: strings.Split(os.Getenv("CORS_METHODS"), ","),
		AllowedHeaders: strings.Split(os.Getenv("CORS_HEADERS"), ","),
	})

	return corsOptions.Handler(router)
}

func (s *Server) router() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("PING SUCCESSFUL"))
	})

	// Quiz endpoints.
	router.HandleFunc("/api/quizzes/all", s.getQuizzes).Methods("POST")
	router.HandleFunc("/api/quizzes/{id}", GetQuiz).Methods("GET")
	router.HandleFunc("/api/quizzes/route/{route}", s.getQuizByRoute).Methods("GET")
	router.HandleFunc("/api/quizzes", CreateQuiz).Methods("POST")
	router.HandleFunc("/api/quizzes/{id}", UpdateQuiz).Methods("PUT")
	router.HandleFunc("/api/quizzes/{id}", DeleteQuiz).Methods("DELETE")

	// Quiz Type endpoints.
	router.HandleFunc("/api/quiztype", GetTypes).Methods("GET")

	// Quiz Plays endpoints.
	router.HandleFunc("/api/quiz-plays", GetAllQuizPlays).Methods("GET")
	router.HandleFunc("/api/quiz-plays/{quizId}", GetQuizPlays).Methods("GET")
	router.HandleFunc("/api/quiz-plays-top-five", GetTopFiveQuizPlays).Methods("GET")
	router.HandleFunc("/api/quiz-plays/{quizId}", IncrementQuizPlays).Methods("PUT")

	// Trivia endpoints.
	router.HandleFunc("/api/trivia/all", s.getAllTrivia).Methods("POST")
	router.HandleFunc("/api/trivia/{date}", s.getTriviaByDate).Methods("GET")
	router.HandleFunc("/api/trivia", GenerateTrivia).Methods("POST")
	router.HandleFunc("/api/trivia/{date}", RegenerateTrivia).Methods("PUT")
	router.HandleFunc("/api/trivia/{date}", DeleteTrivia).Methods("DELETE")
	router.HandleFunc("/api/trivia/old/{newTriviaCount}", DeleteOldTrivia).Methods("DELETE")

	// Trivia Plays endpoints.
	router.HandleFunc("/api/trivia-plays/week", GetLastWeekTriviaPlays).Methods("GET")
	router.HandleFunc("/api/trivia-plays/{id}", IncrementTriviaPlays).Methods("PUT")

	// Trivia Question Type endpoints.
	router.HandleFunc("/api/trivia-question-types", GetTriviaQuestionTypes).Methods("GET")

	// Trivia Question Category endpoints.
	router.HandleFunc("/api/trivia-question-categories", GetTriviaQuestionCategories).Methods("GET")

	// Manual Trivia Question endpoints.
	router.HandleFunc("/api/manual-trivia-questions/all", GetManualTriviaQuestions).Methods("POST")
	router.HandleFunc("/api/manual-trivia-questions", CreateManualTriviaQuestion).Methods("POST")
	router.HandleFunc("/api/manual-trivia-questions/{id}", UpdateManualTriviaQuestion).Methods("PUT")
	router.HandleFunc("/api/manual-trivia-questions/{id}", DeleteManualTriviaQuestion).Methods("DELETE")

	// Mapping endpoints.
	router.HandleFunc("/api/mappings", GetMappingGroups).Methods("GET")
	router.HandleFunc("/api/mappings/{key}", s.getMappingEntries).Methods("GET")
	router.HandleFunc("/api/mappings-no-flags", GetMappingsWithoutFlags).Methods("GET")
	router.HandleFunc("/api/mappings/{key}", EditMapping).Methods("PUT")
	router.HandleFunc("/api/mappings/{key}", DeleteMapping).Methods("DELETE")

	// Map endpoints.
	router.HandleFunc("/api/maps", GetMaps).Methods("GET")
	router.HandleFunc("/api/maps/highlighted/{className}", GetMapHighlightedRegions).Methods("GET")
	router.HandleFunc("/api/maps/{className}", GetMap).Methods("GET")
	router.HandleFunc("/api/maps/preview", GetMapPreview).Methods("POST")
	router.HandleFunc("/api/maps", CreateMap).Methods("POST")

	// Flag endpoints.
	router.HandleFunc("/api/flags", GetFlagGroups).Methods("GET")
	router.HandleFunc("/api/flags/{key}", GetFlagEntries).Methods("GET")
	router.HandleFunc("/api/flags/url/{code}", GetFlagUrl).Methods("GET")
	router.HandleFunc("/api/flags", CreateFlags).Methods("POST")

	// Continent endpoints.
	router.HandleFunc("/api/continents", GetContinents).Methods("GET")

	// Auth endpoints.
	router.HandleFunc("/api/auth/login", Login).Methods("POST")
	router.HandleFunc("/api/auth/register", s.register).Methods("POST")
	router.HandleFunc("/api/auth/refresh", RefreshToken).Methods("POST")
	router.HandleFunc("/api/auth/send-reset-token", s.sendResetToken).Methods("POST")
	router.HandleFunc("/api/auth/reset-token-valid/{userId}/{token}", ResetTokenValid).Methods("GET")
	router.HandleFunc("/api/auth", UpdatePasswordUsingToken).Methods("PUT")
	router.HandleFunc("/api/auth/username/{username}", UsernameExists).Methods("GET")
	router.HandleFunc("/api/auth/email/{email}", EmailExists).Methods("GET")

	// User endpoints.
	router.HandleFunc("/api/users/all", GetUsers).Methods("POST")
	router.HandleFunc("/api/users/{id}", GetUser).Methods("GET")
	router.HandleFunc("/api/users/email/{email}", GetUserByEmail).Methods("GET")
	router.HandleFunc("/api/users/total/week", GetLastWeekTotalUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", s.updateUser).Methods("PUT")
	router.HandleFunc("/api/users/xp/{id}", s.updateUserXP).Methods("PUT")
	router.HandleFunc("/api/users/{id}", DeleteUser).Methods("DELETE")

	// Badge endpoints.
	router.HandleFunc("/api/badges", GetBadges).Methods("GET")
	router.HandleFunc("/api/badges/{userId}", GetUserBadges).Methods("GET")

	// Temp Score endpoints.
	router.HandleFunc("/api/tempscores/{id}", GetTempScore).Methods("GET")
	router.HandleFunc("/api/tempscores", CreateTempScore).Methods("POST")

	// Leaderboard endpoints.
	router.HandleFunc("/api/leaderboard/all/{quizId}", GetEntries).Methods("POST")
	router.HandleFunc("/api/leaderboard/{userId}", GetUserEntries).Methods("GET")
	router.HandleFunc("/api/leaderboard/{quizId}/{userId}", GetEntry).Methods("GET")
	router.HandleFunc("/api/leaderboard", CreateEntry).Methods("POST")
	router.HandleFunc("/api/leaderboard/{id}", UpdateEntry).Methods("PUT")
	router.HandleFunc("/api/leaderboard/{id}", DeleteEntry).Methods("DELETE")

	// Shipping option endpoints.
	router.HandleFunc("/api/shipping-options", GetShippingOptions).Methods("GET")

	// Checkout endpoints.
	router.HandleFunc("/api/checkout/create-checkout-session", HandleCreateCheckoutSession).Methods("POST")
	router.HandleFunc("/api/checkout/webhook", HandleWebhook).Methods("POST")

	// Order endpoints.
	router.HandleFunc("/api/orders", GetOrders).Methods("POST")
	router.HandleFunc("/api/orders/user/{email}", GetUserOrders).Methods("GET")
	router.HandleFunc("/api/orders/status/{id}", UpdateOrderStatus).Methods("PUT")
	router.HandleFunc("/api/orders/{id}", DeleteOrder).Methods("DELETE")
	router.HandleFunc("/api/orders/email/{email}", CancelOrder).Methods("DELETE")

	// Avatar endpoints.
	router.HandleFunc("/api/avatars", s.getAvatars).Methods("GET")

	// Merch endpoints.
	router.HandleFunc("/api/merch", GetMerch).Methods("GET")
	router.HandleFunc("/api/merch/exists", MerchExists).Methods("POST")

	// Discount endpoints.
	router.HandleFunc("/api/discounts", GetDiscounts).Methods("GET")
	router.HandleFunc("/api/discounts/{code}", GetDiscount).Methods("GET")

	// Community Quiz endpoints.
	router.HandleFunc("/api/community-quizzes/all", GetCommunityQuizzes).Methods("POST")
	router.HandleFunc("/api/community-quizzes/{id}", GetCommunityQuiz).Methods("GET")
	router.HandleFunc("/api/community-quizzes/user/{userId}", GetUserCommunityQuizzes).Methods("GET")
	router.HandleFunc("/api/community-quizzes", CreateCommunityQuiz).Methods("POST")
	router.HandleFunc("/api/community-quizzes/{id}", UpdateCommunityQuiz).Methods("PUT")
	router.HandleFunc("/api/community-quizzes/{id}", DeleteCommunityQuiz).Methods("DELETE")

	// Community Quiz Play endpoints.
	router.HandleFunc("/api/community-quiz-plays/{id}", IncrementCommunityQuizPlays).Methods("PUT")

	// SEO endpoints.
	router.HandleFunc("/api/seo/dynamic-routes", GetDynamicRoutes).Methods("GET")

	return router
}
