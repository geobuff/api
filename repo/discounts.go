package repo

import "database/sql"

// Discount is the database object for a discount entry.
type Discount struct {
	ID      int           `json:"id"`
	MerchID sql.NullInt64 `json:"merchId"`
	Code    string        `json:"code"`
	Amount  float64       `json:"amount"`
}

var GetDiscounts = func() ([]Discount, error) {
	rows, err := Connection.Query("SELECT * from discounts;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var discounts = []Discount{}
	for rows.Next() {
		var discount Discount
		if err = rows.Scan(&discount.ID, &discount.MerchID, &discount.Code, &discount.Amount); err != nil {
			return nil, err
		}
		discounts = append(discounts, discount)
	}
	return discounts, rows.Err()
}

// GetDiscount returns a discount by code.
var GetDiscount = func(code string) (Discount, error) {
	statement := "SELECT * from discounts WHERE code = $1;"
	var discount Discount
	err := Connection.QueryRow(statement, code).Scan(&discount.ID, &discount.MerchID, &discount.Code, &discount.Amount)
	return discount, err
}
