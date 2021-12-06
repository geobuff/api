package repo

import (
	"database/sql"
	"time"
)

type Order struct {
	ID        int             `json:"id"`
	StatusID  int             `json:"statusId"`
	Email     string          `json:"email"`
	FirstName string          `json:"firstName"`
	LastName  string          `json:"lastName"`
	Address   string          `json:"address"`
	Suburb    string          `json:"suburb"`
	City      string          `json:"city"`
	Postcode  string          `json:"postcode"`
	Discount  *sql.NullString `json:"discount"`
}

type OrderStatus struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type OrderItem struct {
	ID       int `json:"id"`
	OrderID  int `json:"orderId"`
	MerchID  int `json:"merchId"`
	SizeID   int `json:"sizeId"`
	Quantity int `json:"quantity"`
}

type CheckoutItemDto struct {
	ID       int    `json:"id"`
	Size     string `json:"size"`
	Quantity int    `json:"quantity"`
}

type Customer struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Address   string `json:"address"`
	Suburb    string `json:"suburb"`
	City      string `json:"city"`
	Postcode  string `json:"postcode"`
}

type CreateCheckoutDto struct {
	Items    []CheckoutItemDto `json:"items"`
	Customer Customer          `json:"customer"`
}

var InsertOrder = func(order CreateCheckoutDto) (int, error) {
	statement := "INSERT INTO orders (statusid, email, firstname, lastname, address, suburb, city, postcode, added, discount) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, 1, order.Customer.Email, order.Customer.FirstName, order.Customer.LastName, order.Customer.Address, order.Customer.Suburb, order.Customer.City, order.Customer.Postcode, time.Now(), nil).Scan(&id)
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

// TODO - Get correct size and merch id
func insertOrderItem(item CheckoutItemDto, orderId int) error {
	statement := "INSERT INTO orderItems (orderid, merchid, sizeid, quantity) VALUES ($1, $2, $3, $4) RETURNING id;"
	var id int
	return Connection.QueryRow(statement, orderId, 1, 1, item.Quantity).Scan(&id)
}
