package main

import (
	"fmt"
	"net/http"

	"github.com/geobuff/api/capitals"
	"github.com/geobuff/api/config"
	"github.com/geobuff/api/countries"
	"github.com/geobuff/api/database"
	"github.com/geobuff/api/quizzes"
	"github.com/geobuff/api/scores"
	"github.com/geobuff/api/uk"
	"github.com/geobuff/api/us"
	"github.com/geobuff/api/users"
	"github.com/geobuff/auth"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {
	err := config.Load("config.json")
	if err != nil {
		panic(err)
	}

	err = database.OpenConnection()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to database!")

	err = serve()
	if err != nil {
		panic(err)
	}
}

var serve = func() error {
	return http.ListenAndServe(":8080", handler(router()))
}

func handler(router http.Handler) http.Handler {
	corsOptions := cors.New(cors.Options{
		AllowedOrigins: config.Values.Cors.Origins,
		AllowedMethods: config.Values.Cors.Methods,
		Debug:          true,
	})

	return corsOptions.Handler(router)
}

func router() http.Handler {
	jwtMiddleware := auth.GetJwtMiddleware(config.Values.Auth0.Audience, config.Values.Auth0.Issuer)
	router := mux.NewRouter()

	// Quiz endpoints.
	router.HandleFunc("/api/quizzes", quizzes.GetQuizzes).Methods("GET")
	router.HandleFunc("/api/quizzes/{id}", quizzes.GetQuiz).Methods("GET")

	// User endpoints.
	router.Handle("/api/users", jwtMiddleware.Handler(users.GetUsers)).Methods("GET")
	router.Handle("/api/users/{id}", jwtMiddleware.Handler(users.GetUser)).Methods("GET")
	router.HandleFunc("/api/users/id/{username}", users.GetUserID).Methods("GET")
	router.HandleFunc("/api/users", users.CreateUser).Methods("POST")
	router.Handle("/api/users/{id}", jwtMiddleware.Handler(users.DeleteUser)).Methods("DELETE")

	// Score endpoints.
	router.Handle("/api/scores/{userId}", jwtMiddleware.Handler(scores.GetScores)).Methods("GET")
	router.Handle("/api/scores", jwtMiddleware.Handler(scores.CreateScore)).Methods("POST")
	router.Handle("/api/scores/{id}", jwtMiddleware.Handler(scores.UpdateScore)).Methods("PUT")
	router.Handle("/api/scores/{id}", jwtMiddleware.Handler(scores.DeleteScore)).Methods("DELETE")

	// Countries endpoints.
	router.HandleFunc("/api/countries", countries.GetCountries).Methods("GET")
	router.HandleFunc("/api/countries/leaderboard", countries.GetEntries).Methods("POST")
	router.HandleFunc("/api/countries/leaderboard/{userId}", countries.GetEntry).Methods("GET")
	router.Handle("/api/countries/leaderboard", jwtMiddleware.Handler(countries.CreateEntry)).Methods("POST")
	router.Handle("/api/countries/leaderboard/{id}", jwtMiddleware.Handler(countries.UpdateEntry)).Methods("PUT")
	router.Handle("/api/countries/leaderboard/{id}", jwtMiddleware.Handler(countries.DeleteEntry)).Methods("DELETE")

	// Capitals endpoints.
	router.HandleFunc("/api/capitals", capitals.GetCapitals).Methods("GET")
	router.HandleFunc("/api/capitals/leaderboard", capitals.GetEntries).Methods("POST")
	router.HandleFunc("/api/capitals/leaderboard/{userId}", capitals.GetEntry).Methods("GET")
	router.Handle("/api/capitals/leaderboard", jwtMiddleware.Handler(capitals.CreateEntry)).Methods("POST")
	router.Handle("/api/capitals/leaderboard/{id}", jwtMiddleware.Handler(capitals.UpdateEntry)).Methods("PUT")
	router.Handle("/api/capitals/leaderboard/{id}", jwtMiddleware.Handler(capitals.DeleteEntry)).Methods("DELETE")

	// UK endpoints.
	router.HandleFunc("/api/uk", uk.GetCounties).Methods("GET")

	// US endpoints.
	router.HandleFunc("/api/us", us.GetStates).Methods("GET")

	return router
}
