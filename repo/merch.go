package repo

import (
	"database/sql"
	"fmt"
)

type Merch struct {
	ID                int             `json:"id"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	SizeGuideImageUrl sql.NullString  `json:"sizeGuideImageUrl"`
	Price             sql.NullFloat64 `json:"price"`
	ExternalLink      sql.NullString  `json:"externalLink"`
	Route             string          `json:"route"`
}

type MerchDto struct {
	ID                int             `json:"id"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	SizeGuideImageUrl sql.NullString  `json:"sizeGuideImageUrl"`
	Price             sql.NullFloat64 `json:"price"`
	ExternalLink      sql.NullString  `json:"externalLink"`
	Route             string          `json:"route"`
	Sizes             []MerchSize     `json:"sizes"`
	Images            []MerchImage    `json:"images"`
	SoldOut           bool            `json:"soldOut"`
}

type CartItemDto struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	SizeID      int     `json:"sizeId"`
	SizeName    string  `json:"sizeName"`
	ImageURL    string  `json:"imageUrl"`
	Quantity    int     `json:"quantity"`
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
		if err = rows.Scan(&entry.ID, &entry.Name, &entry.Description, &entry.SizeGuideImageUrl, &entry.Price, &entry.ExternalLink, &entry.Route); err != nil {
			return nil, err
		}

		sizes, err := getMerchSizes(entry.ID)
		if err != nil {
			return nil, err
		}

		images, err := getMerchImages(entry.ID)
		if err != nil {
			return nil, err
		}

		entry.Images = images
		entry.Sizes = sizes
		entry.SoldOut = isSoldOut(sizes)
		merch = append(merch, entry)
	}
	return merch, rows.Err()
}

func GetMerchItem(id int) (*MerchDto, error) {
	statement := "SELECT * from merch WHERE id = $1;"
	var entry MerchDto
	if err := Connection.QueryRow(statement, id).Scan(&entry.ID, &entry.Name, &entry.Description, &entry.SizeGuideImageUrl, &entry.Price, &entry.ExternalLink, &entry.Route); err != nil {
		return nil, err
	}

	sizes, err := getMerchSizes(entry.ID)
	if err != nil {
		return nil, err
	}

	images, err := getMerchImages(entry.ID)
	if err != nil {
		return nil, err
	}

	entry.Images = images
	entry.Sizes = sizes
	entry.SoldOut = isSoldOut(sizes)

	return &entry, nil
}

func isSoldOut(sizes []MerchSize) bool {
	for _, size := range sizes {
		if size.Quantity > 0 {
			return false
		}
	}
	return true
}

func GetMerchRoutes() ([]string, error) {
	rows, err := Connection.Query("SELECT route FROM merch;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes = []string{}
	for rows.Next() {
		var route string
		if err = rows.Scan(&route); err != nil {
			return nil, err
		}
		routes = append(routes, fmt.Sprintf("merch/%s", route))
	}
	return routes, rows.Err()
}
