package repo

import (
	"database/sql"
	"time"
)

type Order struct {
	ID        int            `json:"id"`
	StatusID  int            `json:"statusId"`
	Email     string         `json:"email"`
	FirstName string         `json:"firstName"`
	LastName  string         `json:"lastName"`
	Address   string         `json:"address"`
	Added     time.Time      `json:"added"`
	Discount  sql.NullString `json:"discount"`
}

type OrderStatus struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type CheckoutItemDto struct {
	ID       int    `json:"id"`
	SizeID   int    `json:"sizeId"`
	SizeName string `json:"sizeName"`
	Quantity int    `json:"quantity"`
}

type Customer struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Address   string `json:"address"`
}

type CreateCheckoutDto struct {
	Items    []CheckoutItemDto `json:"items"`
	Customer Customer          `json:"customer"`
}

type OrderDto struct {
	Id        int            `json:"id"`
	Items     []OrderItemDto `json:"items"`
	StatusId  int            `json:"statusId"`
	Status    string         `json:"status"`
	FirstName string         `json:"firstName"`
	LastName  string         `json:"lastName"`
	Address   string         `json:"address"`
	Added     time.Time      `json:"added"`
	Discount  sql.NullString `json:"discount"`
}

type OrdersFilterDto struct {
	StatusID int `json:"statusId"`
	Page     int `json:"page"`
	Limit    int `json:"limit"`
}

func InsertOrder(order CreateCheckoutDto) (int, error) {
	statement := "INSERT INTO orders (statusid, email, firstname, lastname, address, added, discount) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, ORDER_STATUS_PENDING, order.Customer.Email, order.Customer.FirstName, order.Customer.LastName, order.Customer.Address, time.Now(), nil).Scan(&id)
	if err != nil {
		return 0, err
	}

	for _, item := range order.Items {
		err := insertOrderItem(item, id)
		if err != nil {
			return 0, err
		}
	}

	return id, err
}

func GetNonPendingOrders(email string) ([]OrderDto, error) {
	rows, err := Connection.Query("SELECT * FROM orders WHERE email = $1 AND statusid != $2;", email, ORDER_STATUS_PENDING)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders = []OrderDto{}
	for rows.Next() {
		var order Order
		if err = rows.Scan(&order.ID, &order.StatusID, &order.Email, &order.FirstName, &order.LastName, &order.Address, &order.Added, &order.Discount); err != nil {
			return nil, err
		}

		items, err := GetOrderItems(order.ID)
		if err != nil {
			return nil, err
		}

		status, err := getStatus(order.StatusID)
		if err != nil {
			return nil, err
		}

		temp := OrderDto{
			Id:        order.ID,
			Items:     items,
			StatusId:  order.StatusID,
			Status:    status,
			FirstName: order.FirstName,
			LastName:  order.LastName,
			Address:   order.Address,
			Added:     order.Added,
			Discount:  order.Discount,
		}
		orders = append(orders, temp)
	}
	return orders, rows.Err()
}

func getStatus(id int) (string, error) {
	statement := "SELECT status from orderStatus WHERE id = $1;"
	var result string
	err := Connection.QueryRow(statement, id).Scan(&result)
	return result, err
}

func UpdateStatusLatestOrder(email string) (int, error) {
	statement := "UPDATE orders set statusId = $1 where id = (select id from orders where email = $2 order by added desc LIMIT 1) returning id;"
	var id int
	err := Connection.QueryRow(statement, ORDER_STATUS_PAYMENT_RECEIVED, email).Scan(&id)
	return id, err
}

func DeleteOrder(orderId int) error {
	statement := "DELETE FROM orders WHERE id = $1 returning id;"
	var id int
	return Connection.QueryRow(statement, orderId).Scan(&id)
}

func RemoveLatestPendingOrder(email string) error {
	statement := "SELECT id from orders where email = $1 AND statusid = $2 order by added desc LIMIT 1;"
	var orderId int
	err := Connection.QueryRow(statement, email, ORDER_STATUS_PENDING).Scan(&orderId)
	if err != nil {
		return err
	}

	Connection.QueryRow("DELETE from orderItems where orderid = $1;", orderId)
	var id int
	return Connection.QueryRow("DELETE from orders where id = $1;", orderId).Scan(&id)
}

func GetOrders(filter OrdersFilterDto) ([]OrderDto, error) {
	statement := "SELECT * FROM orders WHERE statusid = $1 LIMIT $2 OFFSET $3;"
	rows, err := Connection.Query(statement, filter.StatusID, filter.Limit, filter.Limit*filter.Page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders = []OrderDto{}
	for rows.Next() {
		var order Order
		if err = rows.Scan(&order.ID, &order.StatusID, &order.Email, &order.FirstName, &order.LastName, &order.Address, &order.Added, &order.Discount); err != nil {
			return nil, err
		}

		items, err := GetOrderItems(order.ID)
		if err != nil {
			return nil, err
		}

		status, err := getStatus(order.StatusID)
		if err != nil {
			return nil, err
		}

		temp := OrderDto{
			Items:     items,
			Id:        order.ID,
			StatusId:  order.StatusID,
			Status:    status,
			FirstName: order.FirstName,
			LastName:  order.LastName,
			Address:   order.Address,
			Added:     order.Added,
			Discount:  order.Discount,
		}
		orders = append(orders, temp)
	}
	return orders, rows.Err()
}

var GetFirstOrderID = func(statusID, offset int) (int, error) {
	statement := "SELECT id FROM orders WHERE statusid = $1 LIMIT 1 OFFSET $2;"
	var id int
	err := Connection.QueryRow(statement, statusID, offset).Scan(&id)
	return id, err
}

func UpdateOrderStatus(orderID, statusID int) error {
	statement := "UPDATE orders set statusId = $1 where id = $2 returning id;"
	var id int
	return Connection.QueryRow(statement, statusID, orderID).Scan(&id)
}
