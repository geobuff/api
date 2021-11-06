package repo

import (
	"database/sql"
)

type Merch struct {
	ID           int             `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Price        sql.NullFloat64 `json:"price"`
	ExternalLink sql.NullString  `json:"externalLink"`
	Sizes        []MerchSize     `json:"sizes"`
	Images       []MerchImage    `json:"images"`
}

type MerchSize struct {
	ID       int    `json:"id"`
	MerchID  int    `json:"merchId"`
	Size     string `json:"size"`
	Quantity int    `json:"quantity"`
}

type MerchImage struct {
	ID        int    `json:"id"`
	MerchID   int    `json:"merchId"`
	ImageUrl  string `json:"imageUrl"`
	IsPrimary bool   `json:"isPrimary"`
}

type MerchDto struct {
	ID           int             `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Price        sql.NullFloat64 `json:"price"`
	ExternalLink sql.NullString  `json:"externalLink"`
	Sizes        []MerchSize     `json:"sizes"`
	Images       []MerchImage    `json:"images"`
	SoldOut      bool            `json:"soldOut"`
}

var GetMerch = func() ([]MerchDto, error) {
	rows, err := Connection.Query("SELECT * FROM merch;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var merch = []MerchDto{}
	for rows.Next() {
		var entry MerchDto
		if err = rows.Scan(&entry.ID, &entry.Name, &entry.Description, &entry.Price, &entry.ExternalLink); err != nil {
			return nil, err
		}

		sizeQuery := "SELECT * FROM merchsizes WHERE merchid = $1;"
		rows, err := Connection.Query(sizeQuery, entry.ID)
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
		entry.Sizes = sizes
		entry.SoldOut = isSoldOut(sizes)

		imageQuery := "SELECT * FROM merchImages WHERE merchid = $1;"
		rows, err = Connection.Query(imageQuery, entry.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var images = []MerchImage{}
		for rows.Next() {
			var image MerchImage
			if err = rows.Scan(&image.ID, &image.MerchID, &image.ImageUrl, &image.IsPrimary); err != nil {
				return nil, err
			}
			images = append(images, image)
		}
		entry.Images = images

		merch = append(merch, entry)
	}
	return merch, rows.Err()
}

func isSoldOut(sizes []MerchSize) bool {
	for _, size := range sizes {
		if size.Quantity > 0 {
			return false
		}
	}
	return true
}
