package main

import (
	"fmt"
	"net/http"

	"github.com/geobuff/api/badges"
	"github.com/geobuff/api/capitals"
	"github.com/geobuff/api/config"
	"github.com/geobuff/api/countries"
	"github.com/geobuff/api/database"
	"github.com/geobuff/api/mappings"
	"github.com/geobuff/api/quizzes"
	"github.com/geobuff/api/scores"
	"github.com/geobuff/api/users"
	"github.com/geobuff/auth"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	fmt.Println("successfully connected to database")

	driver, err := postgres.WithInstance(database.Connection, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)

	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
	fmt.Println("successfully ran database migrations")

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
		AllowedHeaders: []string{"*"},
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

	// Mapping endpoints.
	router.HandleFunc("/api/mappings/{key}", mappings.GetMapping).Methods("GET")

	// User endpoints.
	router.Handle("/api/users", jwtMiddleware.Handler(users.GetUsers)).Methods("GET")
	router.Handle("/api/users/{id}", jwtMiddleware.Handler(users.GetUser)).Methods("GET")
	router.HandleFunc("/api/users/id/{username}", users.GetUserID).Methods("GET")
	router.HandleFunc("/api/users", users.CreateUser).Methods("POST")
	router.HandleFunc("/api/users/{id}", users.UpdateUser).Methods("PUT")
	router.Handle("/api/users/{id}", jwtMiddleware.Handler(users.DeleteUser)).Methods("DELETE")

	// Badge endpoints.
	router.HandleFunc("/api/badges", badges.GetBadges).Methods("GET")

	// Score endpoints.
	router.Handle("/api/scores/{userId}", jwtMiddleware.Handler(scores.GetScores)).Methods("GET")
	router.Handle("/api/scores/{userId}/{quizId}", jwtMiddleware.Handler(scores.GetScore)).Methods("GET")
	router.Handle("/api/scores", jwtMiddleware.Handler(scores.CreateScore)).Methods("POST")
	router.Handle("/api/scores/{id}", jwtMiddleware.Handler(scores.UpdateScore)).Methods("PUT")
	router.Handle("/api/scores/{id}", jwtMiddleware.Handler(scores.DeleteScore)).Methods("DELETE")

	// World Countries endpoints.
	router.HandleFunc("/api/world-countries/leaderboard/all", countries.GetEntries).Methods("POST")
	router.HandleFunc("/api/world-countries/leaderboard/{userId}", countries.GetEntry).Methods("GET")
	router.Handle("/api/world-countries/leaderboard", jwtMiddleware.Handler(countries.CreateEntry)).Methods("POST")
	router.Handle("/api/world-countries/leaderboard/{id}", jwtMiddleware.Handler(countries.UpdateEntry)).Methods("PUT")
	router.Handle("/api/world-countries/leaderboard/{id}", jwtMiddleware.Handler(countries.DeleteEntry)).Methods("DELETE")

	// World Capitals endpoints.
	router.HandleFunc("/api/world-capitals/leaderboard/all", capitals.GetEntries).Methods("POST")
	router.HandleFunc("/api/world-capitals/leaderboard/{userId}", capitals.GetEntry).Methods("GET")
	router.Handle("/api/world-capitals/leaderboard", jwtMiddleware.Handler(capitals.CreateEntry)).Methods("POST")
	router.Handle("/api/world-capitals/leaderboard/{id}", jwtMiddleware.Handler(capitals.UpdateEntry)).Methods("PUT")
	router.Handle("/api/world-capitals/leaderboard/{id}", jwtMiddleware.Handler(capitals.DeleteEntry)).Methods("DELETE")

	return router
}
