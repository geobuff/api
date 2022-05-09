package seo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/geobuff/api/repo"
)

func GetDynamicRoutes(writer http.ResponseWriter, request *http.Request) {
	var routes []string
	quizRoutes, err := repo.GetQuizRoutes()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	merchRoutes, err := repo.GetMerchRoutes()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	routes = append(routes, quizRoutes...)
	routes = append(routes, merchRoutes...)

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(routes)
}
