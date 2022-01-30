package orders

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/geobuff/api/auth"
	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

type UpdateOrderStatusDto struct {
	StatusID int `json:"statusId"`
}

func GetOrdersByStatusId(writer http.ResponseWriter, request *http.Request) {
	statusId, err := strconv.Atoi(mux.Vars(request)["statusId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	if code, err := auth.IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	orders, err := repo.GetOrdersByStatusId(statusId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(orders)
}

func GetUserOrders(writer http.ResponseWriter, request *http.Request) {
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

func UpdateOrderStatus(writer http.ResponseWriter, request *http.Request) {
	if code, err := auth.IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var dto UpdateOrderStatusDto
	err = json.Unmarshal(requestBody, &dto)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	err = repo.UpdateOrderStatus(id, dto.StatusID)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

func CancelOrder(writer http.ResponseWriter, request *http.Request) {
	err := repo.RemoveLatestPendingOrder(mux.Vars(request)["email"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}
