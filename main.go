package main

import (
	"fmt"
	"net/http"

	"github.com/geobuff/geobuff-api/auth"
	"github.com/geobuff/geobuff-api/config"
	"github.com/geobuff/geobuff-api/database"
	"github.com/geobuff/geobuff-api/isocodes"
	"github.com/geobuff/geobuff-api/scores"
	"github.com/geobuff/geobuff-api/users"
	"github.com/geobuff/geobuff-api/world"
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
	jwtMiddleware := auth.GetJwtMiddleware()
	router := mux.NewRouter()

	// User endpoints.
	router.Handle("/api/users", jwtMiddleware.Handler(users.GetUsers)).Methods("GET")
	router.Handle("/api/users/{id}", jwtMiddleware.Handler(users.GetUser)).Methods("GET")
	router.HandleFunc("/api/users", users.CreateUser).Methods("POST")
	router.Handle("/api/users/{id}", jwtMiddleware.Handler(users.DeleteUser)).Methods("DELETE")

	// Score endpoints.
	router.Handle("/api/scores/{id}", jwtMiddleware.Handler(scores.GetScores)).Methods("GET")
	router.Handle("/api/scores", jwtMiddleware.Handler(scores.CreateScore)).Methods("POST")
	router.Handle("/api/scores", jwtMiddleware.Handler(scores.UpdateScore)).Methods("PUT")
	router.Handle("/api/scores/{id}", jwtMiddleware.Handler(scores.DeleteScore)).Methods("DELETE")

	// ISO-code endpoints.
	router.HandleFunc("/api/isocodes", isocodes.GetCodes).Methods("GET")

	// World endpoints.
	router.HandleFunc("/api/world/countries", world.GetCountries).Methods("GET")
	router.HandleFunc("/api/world/countries/alternatives", world.GetAlternativeNamings).Methods("GET")
	router.HandleFunc("/api/world/countries/prefixes", world.GetPrefixes).Methods("GET")
	router.HandleFunc("/api/world/countries/map", world.GetCountriesMap).Methods("GET")
	router.HandleFunc("/api/world/leaderboard", world.GetEntries).Methods("GET")
	router.HandleFunc("/api/world/leaderboard/{id}", world.GetEntry).Methods("GET")
	router.Handle("/api/world/leaderboard", jwtMiddleware.Handler(world.CreateEntry)).Methods("POST")
	router.Handle("/api/world/leaderboard/{id}", jwtMiddleware.Handler(world.UpdateEntry)).Methods("PUT")
	router.Handle("/api/world/leaderboard/{id}", jwtMiddleware.Handler(world.DeleteEntry)).Methods("DELETE")

	return router
}
