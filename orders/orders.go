package orders

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

func GetOrders(writer http.ResponseWriter, request *http.Request) {
	orders, err := repo.GetOrders(mux.Vars(request)["email"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(orders)
}
