package repo

type MapElementType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetMapElementTypeId(name string) (int, error) {
	var id int
	err := Connection.QueryRow("SELECT id FROM mapelementtype WHERE name = $1;", name).Scan(&id)
	return id, err
}
