package repo

type MerchImage struct {
	ID        int    `json:"id"`
	MerchID   int    `json:"merchId"`
	ImageUrl  string `json:"imageUrl"`
	IsPrimary bool   `json:"isPrimary"`
}

func getMerchImages(merchID int) ([]MerchImage, error) {
	rows, err := Connection.Query("SELECT * FROM merchImages WHERE merchid = $1;", merchID)
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
	return images, rows.Err()
}
