package src

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

func GetDiscounts(writer http.ResponseWriter, request *http.Request) {
	if code, err := IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	discounts, err := repo.GetDiscounts()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(discounts)
}

func GetDiscount(writer http.ResponseWriter, request *http.Request) {
	switch discount, err := repo.GetDiscountByCode(mux.Vars(request)["code"]); err {
	case sql.ErrNoRows:
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusNoContent)
	case nil:
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(discount)
	default:
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
	}
}
