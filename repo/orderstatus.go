package repo

const (
	ORDER_STATUS_PENDING int = iota + 1
	ORDER_STATUS_PAYMENT_RECEIVED
	ORDER_STATUS_SHIPPED
)

func getOrderStatus(id int) (string, error) {
	statement := "SELECT status from orderStatus WHERE id = $1;"
	var result string
	err := Connection.QueryRow(statement, id).Scan(&result)
	return result, err
}
