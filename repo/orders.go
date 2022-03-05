package repo

import (
	"database/sql"
	"time"
)

type Order struct {
	ID         int           `json:"id"`
	StatusID   int           `json:"statusId"`
	ShippingId int           `json:"shippingId"`
	DiscountId sql.NullInt64 `json:"discountId"`
	Email      string        `json:"email"`
	FirstName  string        `json:"firstName"`
	LastName   string        `json:"lastName"`
	Address    string        `json:"address"`
	Added      time.Time     `json:"added"`
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
	Items      []CheckoutItemDto `json:"items"`
	Customer   Customer          `json:"customer"`
	ShippingId int               `json:"shippingId"`
	DiscountId sql.NullInt64     `json:"discountId"`
}

type OrderDto struct {
	ID             int            `json:"id"`
	StatusID       int            `json:"statusId"`
	Status         string         `json:"status"`
	ShippingOption string         `json:"shippingOption"`
	Discount       string         `json:"discount"`
	FirstName      string         `json:"firstName"`
	LastName       string         `json:"lastName"`
	Address        string         `json:"address"`
	Added          time.Time      `json:"added"`
	Items          []OrderItemDto `json:"items"`
}

type OrdersFilterDto struct {
	StatusID int `json:"statusId"`
	Page     int `json:"page"`
	Limit    int `json:"limit"`
}

func GetOrders(filter OrdersFilterDto) ([]OrderDto, error) {
	statement := "SELECT o.id, s.name, d.name, o.firstname, o.lastname, o.address, o.added FROM orders o JOIN shippingoptions s ON s.id = o.shippingid LEFT JOIN discounts d ON d.id = o.discountid WHERE o.statusid = $1 LIMIT $2 OFFSET $3;"
	rows, err := Connection.Query(statement, filter.StatusID, filter.Limit, filter.Limit*filter.Page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders = []OrderDto{}
	for rows.Next() {
		var order OrderDto
		if err = rows.Scan(&order.ID, &order.StatusID, &order.ShippingOption, &order.Discount, &order.FirstName, &order.LastName, &order.Address, &order.Added); err != nil {
			return nil, err
		}

		items, err := GetOrderItems(order.ID)
		if err != nil {
			return nil, err
		}

		status, err := getOrderStatus(order.StatusID)
		if err != nil {
			return nil, err
		}

		order.Items = items
		order.Status = status
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

var GetFirstOrderID = func(statusID, offset int) (int, error) {
	statement := "SELECT id FROM orders WHERE statusid = $1 LIMIT 1 OFFSET $2;"
	var id int
	err := Connection.QueryRow(statement, statusID, offset).Scan(&id)
	return id, err
}

func GetNonPendingOrders(email string) ([]OrderDto, error) {
	statement := "SELECT o.id, s.name, d.name, o.firstname, o.lastname, o.address, o.added FROM orders o JOIN shippingoptions s ON s.id = o.shippingid LEFT JOIN discounts d ON d.id = o.discountid WHERE o.email = $1 AND o.statusid != $2;"
	rows, err := Connection.Query(statement, email, ORDER_STATUS_PENDING)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders = []OrderDto{}
	for rows.Next() {
		var order OrderDto
		if err = rows.Scan(&order.ID, &order.StatusID, &order.ShippingOption, &order.Discount, &order.FirstName, &order.LastName, &order.Address, &order.Added); err != nil {
			return nil, err
		}

		items, err := GetOrderItems(order.ID)
		if err != nil {
			return nil, err
		}

		status, err := getOrderStatus(order.StatusID)
		if err != nil {
			return nil, err
		}

		order.Items = items
		order.Status = status
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func InsertOrder(order CreateCheckoutDto) (int, error) {
	statement := "INSERT INTO orders (statusid, shippingid, discountid, email, firstname, lastname, address, added) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, ORDER_STATUS_PENDING, order.ShippingId, order.DiscountId, order.Customer.Email, order.Customer.FirstName, order.Customer.LastName, order.Customer.Address, time.Now()).Scan(&id)
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

func UpdateOrderStatus(orderID, statusID int) error {
	statement := "UPDATE orders set statusId = $1 where id = $2 returning id;"
	var id int
	return Connection.QueryRow(statement, statusID, orderID).Scan(&id)
}
