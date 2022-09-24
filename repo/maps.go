package repo

type Map struct {
	ID        int    `json:"id"`
	Key       string `json:"key"`
	ClassName string `json:"className"`
	Label     string `json:"label"`
	ViewBox   string `json:"viewBox"`
}

type MapDto struct {
	ID        int             `json:"id"`
	Key       string          `json:"key"`
	ClassName string          `json:"className"`
	Label     string          `json:"label"`
	ViewBox   string          `json:"viewBox"`
	Elements  []MapElementDto `json:"elements"`
}

type GetMapsDto struct {
	Name  string `json:"name"`
	Key   string `json:"key"`
	Label string `json:"label"`
	Value string `json:"value"`
}

type CreateMapDto struct {
	ID        int                   `json:"id"`
	Key       string                `json:"key"`
	ClassName string                `json:"className"`
	Label     string                `json:"label"`
	ViewBox   string                `json:"viewBox"`
	Elements  []CreateMapElementDto `json:"elements"`
}

type HighlightedRegionDto struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func GetMaps() ([]GetMapsDto, error) {
	rows, err := Connection.Query("SELECT label, key, label, classname from maps;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var maps = []GetMapsDto{}
	for rows.Next() {
		var m GetMapsDto
		if err = rows.Scan(&m.Name, &m.Key, &m.Label, &m.Value); err != nil {
			return nil, err
		}
		maps = append(maps, m)
	}

	return maps, rows.Err()
}

func GetMap(className string) (MapDto, error) {
	statement := "SELECT * from maps WHERE classname = $1;"
	var m MapDto
	err := Connection.QueryRow(statement, className).Scan(&m.ID, &m.Key, &m.ClassName, &m.Label, &m.ViewBox)
	if err != nil {
		return MapDto{}, err
	}

	elements, err := GetMapElements(m.ID)
	if err != nil {
		return MapDto{}, err
	}

	m.Elements = elements
	return m, nil
}

func GetMapUsingKey(key string) (MapDto, error) {
	statement := "SELECT * from maps WHERE key = $1;"
	var m MapDto
	err := Connection.QueryRow(statement, key).Scan(&m.ID, &m.Key, &m.ClassName, &m.Label, &m.ViewBox)
	if err != nil {
		return MapDto{}, err
	}

	elements, err := GetMapElements(m.ID)
	if err != nil {
		return MapDto{}, err
	}

	m.Elements = elements
	return m, nil
}

func GetMapHighlightedRegions(className string) ([]HighlightedRegionDto, error) {
	statement := "SELECT id from maps WHERE classname = $1;"
	var mapId int
	err := Connection.QueryRow(statement, className).Scan(&mapId)
	if err != nil {
		return nil, err
	}

	regions, err := GetHighlightedElements(mapId)
	if err != nil {
		return nil, err
	}

	return regions, nil
}

func CreateMap(svgMap MapDto) error {
	var id int
	statement := "INSERT INTO maps (key, classname, label, viewbox) values ($1, $2, $3, $4) RETURNING id;"
	if err := Connection.QueryRow(statement, svgMap.Key, svgMap.ClassName, svgMap.Label, svgMap.ViewBox).Scan(&id); err != nil {
		return err
	}

	for _, element := range svgMap.Elements {
		if err := CreateMapElement(id, element); err != nil {
			return err
		}
	}
	return nil
}

func GetMapId(key string) (int, error) {
	statement := "SELECT id from maps WHERE key = $1;"
	var id int
	err := Connection.QueryRow(statement, key).Scan(&id)
	return id, err
}

func DeleteMap(mapId int) error {
	var id int
	return Connection.QueryRow("DELETE FROM maps where id = $1 RETURNING id;", mapId).Scan(&id)
}
