package badges

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/geobuff/api/repo"
)

// GetBadges returns all badges.
func GetBadges(writer http.ResponseWriter, request *http.Request) {
	badges, err := repo.GetBadges()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(badges)
}
