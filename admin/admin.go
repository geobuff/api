package admin

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/geobuff/api/auth"
	"github.com/geobuff/api/repo"
)

func GetDashboardData(writer http.ResponseWriter, request *http.Request) {
	if code, err := auth.IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	data, err := repo.GetAdminDashboardData()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(data)
}
