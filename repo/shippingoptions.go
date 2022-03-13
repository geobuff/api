package repo

type ShippingOption struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"imageUrl"`
}

func GetShippingOptions() ([]ShippingOption, error) {
	rows, err := Connection.Query("SELECT * from shippingoptions;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options = []ShippingOption{}
	for rows.Next() {
		var option ShippingOption
		if err = rows.Scan(&option.ID, &option.Name, &option.Description, &option.Price, &option.ImageURL); err != nil {
			return nil, err
		}
		options = append(options, option)
	}
	return options, rows.Err()
}

func GetShippingOption(id int) (ShippingOption, error) {
	statement := "SELECT * from shippingoptions WHERE id = $1;"
	var option ShippingOption
	err := Connection.QueryRow(statement, id).Scan(&option.ID, &option.Name, &option.Description, &option.Price, &option.ImageURL)
	return option, err
}
