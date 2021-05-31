package repo

type Merch struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Price       float64     `json:"price"`
	Sizes       []MerchSize `json:"sizes"`
	ImageUrls   []string    `json:"imageUrls"`
}

type MerchSize struct {
	Size    string `json:"size"`
	SoldOut bool   `json:"soldOut"`
}

func GetMerch() ([]Merch, error) {
	rows, err := Connection.Query("SELECT * FROM merch;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var merch = []Merch{}
	for rows.Next() {
		var entry Merch
		if err = rows.Scan(&entry.ID, &entry.Name, &entry.Description, &entry.Price); err != nil {
			return nil, err
		}

		query := "SELECT size, soldout FROM merchsizes WHERE merchid = $1;"
		rows, err := Connection.Query(query, entry.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var sizes = []MerchSize{}
		for rows.Next() {
			var size MerchSize
			if err = rows.Scan(&size.Size, &size.SoldOut); err != nil {
				return nil, err
			}
			sizes = append(sizes, size)
		}
		entry.Sizes = sizes

		query = "SELECT imageurl FROM merchImages WHERE merchid = $1;"
		rows, err = Connection.Query(query, entry.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var images = []string{}
		for rows.Next() {
			var image string
			if err = rows.Scan(&image); err != nil {
				return nil, err
			}
			images = append(images, image)
		}
		entry.ImageUrls = images

		merch = append(merch, entry)
	}
	return merch, rows.Err()
}
