package repo

type OrderItem struct {
	ID       int `json:"id"`
	OrderID  int `json:"orderId"`
	MerchID  int `json:"merchId"`
	SizeID   int `json:"sizeId"`
	Quantity int `json:"quantity"`
}

type OrderItemDto struct {
	MerchID  int    `json:"merchId"`
	ItemName string `json:"itemName"`
	SizeID   int    `json:"sizeId"`
	SizeName string `json:"sizeName"`
	ImageUrl string `json:"imageUrl"`
	Quantity int    `json:"quantity"`
}

func insertOrderItem(item CheckoutItemDto, orderId int) error {
	statement := "INSERT INTO orderItems (orderid, merchid, sizeid, quantity) VALUES ($1, $2, $3, $4) RETURNING id;"
	var id int
	return Connection.QueryRow(statement, orderId, item.ID, item.SizeID, item.Quantity).Scan(&id)
}

func GetOrderItems(orderID int) ([]OrderItemDto, error) {
	statement := "select i.merchid, m.name, s.id, s.size, mi.imageurl, i.quantity from orderItems i join merchsizes s on s.id = i.sizeid join merch m on m.id = i.merchid join merchimages mi on mi.id = i.merchid AND mi.isprimary WHERE i.orderId = $1;"
	rows, err := Connection.Query(statement, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items = []OrderItemDto{}
	for rows.Next() {
		var item OrderItemDto
		if err = rows.Scan(&item.MerchID, &item.ItemName, &item.SizeID, &item.SizeName, &item.ImageUrl, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
