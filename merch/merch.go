package merch

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/geobuff/api/repo"
)

// GetMerch returns all the merch entries.
func GetMerch(writer http.ResponseWriter, request *http.Request) {
	merch, err := repo.GetMerch()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(merch)
}
