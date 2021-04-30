package main

import (
	"fmt"
	"net/http"

	"github.com/geobuff/api/auth"
	"github.com/geobuff/api/badges"
	"github.com/geobuff/api/config"
	"github.com/geobuff/api/leaderboard"
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

	// Score endpoints.
	router.HandleFunc("/api/scores/{userId}", scores.GetScores).Methods("GET")
	router.HandleFunc("/api/scores/{userId}/{quizId}", scores.GetScore).Methods("GET")
	router.HandleFunc("/api/scores", scores.CreateScore).Methods("POST")
	router.HandleFunc("/api/scores/{id}", scores.UpdateScore).Methods("PUT")
	router.HandleFunc("/api/scores/{id}", scores.DeleteScore).Methods("DELETE")

	// Leaderboard endpoints.
	router.HandleFunc("/api/leaderboard/all/{quizId}", leaderboard.GetEntries).Methods("POST")
	router.HandleFunc("/api/leaderboard/{userId}", leaderboard.GetUserEntries).Methods("GET")
	router.HandleFunc("/api/leaderboard/{quizId}/{userId}", leaderboard.GetEntry).Methods("GET")
	router.HandleFunc("/api/leaderboard", leaderboard.CreateEntry).Methods("POST")
	router.HandleFunc("/api/leaderboard/{id}", leaderboard.UpdateEntry).Methods("PUT")
	router.HandleFunc("/api/leaderboard/{id}", leaderboard.DeleteEntry).Methods("DELETE")

	return router
}
