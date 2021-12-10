package discounts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

// GetDiscount gets a discount for a given code.
func GetDiscount(writer http.ResponseWriter, request *http.Request) {
	switch discount, err := repo.GetDiscount(mux.Vars(request)["code"]); err {
	case sql.ErrNoRows:
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusNoContent)
	case nil:
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(discount)
	default:
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
	}
}