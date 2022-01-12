package dailytriviaplays

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

func IncrementTriviaPlays(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	err = repo.IncrementTriviaPlays(id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}
