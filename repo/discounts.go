package repo

import "database/sql"

// Discount is the database object for a discount entry.
type Discount struct {
	ID      int           `json:"id"`
	MerchID sql.NullInt64 `json:"merchId"`
	Code    string        `json:"code"`
	Amount  float64       `json:"amount"`
}

// GetDiscount returns a discount by code.
var GetDiscount = func(code string) (Discount, error) {
	statement := "SELECT * from discounts WHERE code = $1;"
	var discount Discount
	err := Connection.QueryRow(statement, code).Scan(&discount.ID, &discount.MerchID, &discount.Code, &discount.Amount)
	return discount, err
}
