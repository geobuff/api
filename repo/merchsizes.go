package repo

type MerchSize struct {
	ID       int    `json:"id"`
	MerchID  int    `json:"merchId"`
	Size     string `json:"size"`
	Quantity int    `json:"quantity"`
}

func getMerchSizes(merchID int) ([]MerchSize, error) {
	rows, err := Connection.Query("SELECT * FROM merchsizes WHERE merchid = $1 ORDER BY id;", merchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sizes = []MerchSize{}
	for rows.Next() {
		var size MerchSize
		if err = rows.Scan(&size.ID, &size.MerchID, &size.Size, &size.Quantity); err != nil {
			return nil, err
		}
		sizes = append(sizes, size)
	}
	return sizes, rows.Err()
}

func ReduceMerchItemQuantity(sizeID, decrease int) error {
	statement := "UPDATE merchsizes SET quantity = quantity - $1 WHERE id = $2 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, decrease, sizeID).Scan(&id)
}

func MerchExists(items []CartItemDto) (bool, error) {
	for _, item := range items {
		var quantity int
		err := Connection.QueryRow("SELECT quantity FROM merchsizes WHERE id = $1;", item.SizeID).Scan(&quantity)
		if err != nil || quantity < item.Quantity {
			return false, err
		}
	}
	return true, nil
}
