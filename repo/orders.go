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
	Suburb    string         `json:"suburb"`
	City      string         `json:"city"`
	Postcode  string         `json:"postcode"`
	Added     time.Time      `json:"added"`
	Discount  sql.NullString `json:"discount"`
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
	SizeID   int    `json:"sizeId"`
	SizeName string `json:"sizeName"`
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

type OrderDto struct {
	Id        int            `json:"id"`
	Items     []OrderItemDto `json:"items"`
	Status    string         `json:"status"`
	FirstName string         `json:"firstName"`
	LastName  string         `json:"lastName"`
	Address   string         `json:"address"`
	Suburb    string         `json:"suburb"`
	City      string         `json:"city"`
	Postcode  string         `json:"postcode"`
	Added     time.Time      `json:"added"`
	Discount  sql.NullString `json:"discount"`
}

type OrderItemDto struct {
	ItemName string `json:"itemName"`
	SizeName string `json:"sizeName"`
	ImageUrl string `json:"imageUrl"`
	Quantity int    `json:"quantity"`
}

func InsertOrder(order CreateCheckoutDto) (int, error) {
	statement := "INSERT INTO orders (statusid, email, firstname, lastname, address, suburb, city, postcode, added, discount) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, ORDER_STATUS_PENDING, order.Customer.Email, order.Customer.FirstName, order.Customer.LastName, order.Customer.Address, order.Customer.Suburb, order.Customer.City, order.Customer.Postcode, time.Now(), nil).Scan(&id)
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

func insertOrderItem(item CheckoutItemDto, orderId int) error {
	statement := "INSERT INTO orderItems (orderid, merchid, sizeid, quantity) VALUES ($1, $2, $3, $4) RETURNING id;"
	var id int
	return Connection.QueryRow(statement, orderId, item.ID, item.SizeID, item.Quantity).Scan(&id)
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
		if err = rows.Scan(&order.ID, &order.StatusID, &order.Email, &order.FirstName, &order.LastName, &order.Address, &order.Suburb, &order.City, &order.Postcode, &order.Added, &order.Discount); err != nil {
			return nil, err
		}

		items, err := getItems(order.ID)
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
			Status:    status,
			FirstName: order.FirstName,
			LastName:  order.LastName,
			Address:   order.Address,
			Suburb:    order.Suburb,
			City:      order.City,
			Postcode:  order.Postcode,
			Added:     order.Added,
			Discount:  order.Discount,
		}
		orders = append(orders, temp)
	}
	return orders, rows.Err()
}

func getItems(id int) ([]OrderItemDto, error) {
	rows, err := Connection.Query("select m.name, s.size, mi.imageurl, i.quantity from orderItems i join merchsizes s on s.id = i.sizeid join merch m on m.id = i.merchid join merchimages mi on mi.id = i.merchid AND mi.isprimary WHERE i.orderId = $1;", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items = []OrderItemDto{}
	for rows.Next() {
		var item OrderItemDto
		if err = rows.Scan(&item.ItemName, &item.SizeName, &item.ImageUrl, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func getStatus(id int) (string, error) {
	statement := "SELECT status from orderStatus WHERE id = $1;"
	var result string
	err := Connection.QueryRow(statement, id).Scan(&result)
	return result, err
}

func UpdateStatusLatestOrder(email string) error {
	statement := "UPDATE orders set statusId = $1 where id = (select id from orders where email = $2 order by added desc LIMIT 1) returning id;"
	var id int
	return Connection.QueryRow(statement, ORDER_STATUS_PAYMENT_RECEIVED, email).Scan(&id)
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

func GetOrdersByStatusId(statusId int) ([]OrderDto, error) {
	rows, err := Connection.Query("SELECT * FROM orders WHERE statusid = $1;", statusId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders = []OrderDto{}
	for rows.Next() {
		var order Order
		if err = rows.Scan(&order.ID, &order.StatusID, &order.Email, &order.FirstName, &order.LastName, &order.Address, &order.Suburb, &order.City, &order.Postcode, &order.Added, &order.Discount); err != nil {
			return nil, err
		}

		items, err := getItems(order.ID)
		if err != nil {
			return nil, err
		}

		status, err := getStatus(order.StatusID)
		if err != nil {
			return nil, err
		}

		temp := OrderDto{
			Items:     items,
			Status:    status,
			FirstName: order.FirstName,
			LastName:  order.LastName,
			Address:   order.Address,
			Suburb:    order.Suburb,
			City:      order.City,
			Postcode:  order.Postcode,
			Added:     order.Added,
			Discount:  order.Discount,
		}
		orders = append(orders, temp)
	}
	return orders, rows.Err()
}

func UpdateOrderStatus(orderID, statusID int) error {
	statement := "UPDATE orders set statusId = $1 where id = $2 returning id;"
	var id int
	return Connection.QueryRow(statement, statusID, id).Scan(&id)
}
