package main

import (
	"fmt"
	"net/http"

	"github.com/geobuff/api/auth"
	"github.com/geobuff/api/badges"
	"github.com/geobuff/api/capitals"
	"github.com/geobuff/api/config"
	"github.com/geobuff/api/countries"
	"github.com/geobuff/api/mappings"
	"github.com/geobuff/api/quizzes"
	"github.com/geobuff/api/repo"
	"github.com/geobuff/api/scores"
	"github.com/geobuff/api/users"
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

	err = repo.OpenConnection()
	if err != nil {
		panic(err)
	}
	fmt.Println("successfully connected to database")

	driver, err := postgres.WithInstance(repo.Connection, &postgres.Config{})
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
	router := mux.NewRouter()

	// Quiz endpoints.
	router.HandleFunc("/api/quizzes", quizzes.GetQuizzes).Methods("GET")
	router.HandleFunc("/api/quizzes/{id}", quizzes.GetQuiz).Methods("GET")

	// Mapping endpoints.
	router.HandleFunc("/api/mappings/{key}", mappings.GetMapping).Methods("GET")

	// Auth endpoints.
	router.HandleFunc("/api/auth/login", auth.Login).Methods("POST")
	router.HandleFunc("/api/auth/register", auth.Register).Methods("POST")

	// User endpoints.
	router.HandleFunc("/api/users", users.GetUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", users.GetUser).Methods("GET")
	router.HandleFunc("/api/users/{id}", users.UpdateUser).Methods("PUT")
	router.HandleFunc("/api/users/{id}", users.DeleteUser).Methods("DELETE")

	// Badge endpoints.
	router.HandleFunc("/api/badges", badges.GetBadges).Methods("GET")

	// Score endpoints.
	router.HandleFunc("/api/scores/{userId}", scores.GetScores).Methods("GET")
	router.HandleFunc("/api/scores/{userId}/{quizId}", scores.GetScore).Methods("GET")
	router.HandleFunc("/api/scores", scores.CreateScore).Methods("POST")
	router.HandleFunc("/api/scores/{id}", scores.UpdateScore).Methods("PUT")
	router.HandleFunc("/api/scores/{id}", scores.DeleteScore).Methods("DELETE")

	// World Countries endpoints.
	router.HandleFunc("/api/world-countries/leaderboard/all", countries.GetEntries).Methods("POST")
	router.HandleFunc("/api/world-countries/leaderboard/{userId}", countries.GetEntry).Methods("GET")
	router.HandleFunc("/api/world-countries/leaderboard", countries.CreateEntry).Methods("POST")
	router.HandleFunc("/api/world-countries/leaderboard/{id}", countries.UpdateEntry).Methods("PUT")
	router.HandleFunc("/api/world-countries/leaderboard/{id}", countries.DeleteEntry).Methods("DELETE")

	// World Capitals endpoints.
	router.HandleFunc("/api/world-capitals/leaderboard/all", capitals.GetEntries).Methods("POST")
	router.HandleFunc("/api/world-capitals/leaderboard/{userId}", capitals.GetEntry).Methods("GET")
	router.HandleFunc("/api/world-capitals/leaderboard", capitals.CreateEntry).Methods("POST")
	router.HandleFunc("/api/world-capitals/leaderboard/{id}", capitals.UpdateEntry).Methods("PUT")
	router.HandleFunc("/api/world-capitals/leaderboard/{id}", capitals.DeleteEntry).Methods("DELETE")

	return router
}
