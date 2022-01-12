package repo

type Continent struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var GetContinents = func() ([]Continent, error) {
	rows, err := Connection.Query("SELECT * FROM continents;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var continents = []Continent{}
	for rows.Next() {
		var continent Continent
		if err = rows.Scan(&continent.ID, &continent.Name); err != nil {
			return nil, err
		}
		continents = append(continents, continent)
	}
	return continents, rows.Err()
}
