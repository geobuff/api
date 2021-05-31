package repo

type Merch struct {
	ID          int          `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Price       float64      `json:"price"`
	Disabled    bool         `json:"disabled"`
	Sizes       []MerchSize  `json:"sizes"`
	Images      []MerchImage `json:"images"`
}

type MerchSize struct {
	Size    string `json:"size"`
	SoldOut bool   `json:"soldOut"`
}

type MerchImage struct {
	ImageUrl  string `json:"imageUrl"`
	IsPrimary bool   `json:"isPrimary"`
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
		if err = rows.Scan(&entry.ID, &entry.Name, &entry.Description, &entry.Price, &entry.Disabled); err != nil {
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

		query = "SELECT imageurl, isprimary FROM merchImages WHERE merchid = $1;"
		rows, err = Connection.Query(query, entry.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var images = []MerchImage{}
		for rows.Next() {
			var image MerchImage
			if err = rows.Scan(&image.ImageUrl, &image.IsPrimary); err != nil {
				return nil, err
			}
			images = append(images, image)
		}
		entry.Images = images

		merch = append(merch, entry)
	}
	return merch, rows.Err()
}
