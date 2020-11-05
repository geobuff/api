package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/geobuff/geobuff-api/auth"
	"github.com/geobuff/geobuff-api/database"
	"github.com/geobuff/geobuff-api/isocodes"
	"github.com/geobuff/geobuff-api/world"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	var err error
	database.DBConnection, err = sql.Open("postgres", database.GetConnectionString())
	if err != nil {
		panic(err)
	}

	err = database.DBConnection.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to database!")

	jwtMiddleware := auth.GetJwtMiddleware()
	router := mux.NewRouter()

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

	// ISO-code endpoints.
	router.HandleFunc("/api/isocodes", isocodes.GetCodes).Methods("GET")

	handler := cors.Default().Handler(router)
	http.ListenAndServe(":8080", handler)
}
