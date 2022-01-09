package orders

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/geobuff/api/auth"
	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

func GetOrders(writer http.ResponseWriter, request *http.Request) {
	email := mux.Vars(request)["email"]
	user, err := repo.GetAuthUserUsingEmail(email)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	if code, err := auth.ValidUser(request, user.ID); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	orders, err := repo.GetNonPendingOrders(email)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(orders)
}

func CancelOrder(writer http.ResponseWriter, request *http.Request) {
	err := repo.RemoveLatestOrder(mux.Vars(request)["email"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}
