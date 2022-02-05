package orders

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/geobuff/api/auth"
	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

type OrderPageDto struct {
	Orders  []repo.OrderDto `json:"orders"`
	HasMore bool            `json:"hasMore"`
}

type UpdateOrderStatusDto struct {
	StatusID int `json:"statusId"`
}

func GetOrders(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var filter repo.OrdersFilterDto
	err = json.Unmarshal(requestBody, &filter)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	if code, err := auth.IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	orders, err := repo.GetOrders(filter)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	switch _, err := repo.GetFirstOrderID(filter.StatusID, (filter.Page+1)*filter.Limit); err {
	case sql.ErrNoRows:
		entriesDto := OrderPageDto{orders, false}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entriesDto)
	case nil:
		entriesDto := OrderPageDto{orders, true}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entriesDto)
	default:
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
	}
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

func DeleteOrder(writer http.ResponseWriter, request *http.Request) {
	if code, err := auth.IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	if err = repo.DeleteOrder(id); err != nil {
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
