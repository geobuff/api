package repo

type Map struct {
	ID        int    `json:"id"`
	Key       string `json:"key"`
	ClassName string `json:"className"`
	Label     string `json:"label"`
	ViewBox   string `json:"viewBox"`
}

type MapElement struct {
	ID         int    `json:"id"`
	MapID      int    `json:"mapId"`
	TypeID     int    `json:"typeId"`
	ElementID  string `json:"elementId"`
	Name       string `json:"name"`
	D          string `json:"d"`
	Points     string `json:"points"`
	X          string `json:"x"`
	Y          string `json:"y"`
	Width      string `json:"width"`
	Height     string `json:"height"`
	Cx         string `json:"cx"`
	Cy         string `json:"cy"`
	R          string `json:"r"`
	Transform  string `json:"transform"`
	XlinkHref  string `json:"xlinkHref"`
	ClipPath   string `json:"clipPath"`
	ClipPathId string `json:"clipPathId"`
	X1         string `json:"x1"`
	Y1         string `json:"y1"`
	X2         string `json:"x2"`
	Y2         string `json:"y2"`
}

type MapElementType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MapDto struct {
	ID        int             `json:"id"`
	Key       string          `json:"key"`
	ClassName string          `json:"className"`
	Label     string          `json:"label"`
	ViewBox   string          `json:"viewBox"`
	Elements  []MapElementDto `json:"elements"`
}

type MapElementDto struct {
	MapID      int    `json:"mapId"`
	Type       string `json:"type"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	D          string `json:"d"`
	Points     string `json:"points"`
	X          string `json:"x"`
	Y          string `json:"y"`
	Width      string `json:"width"`
	Height     string `json:"height"`
	Cx         string `json:"cx"`
	Cy         string `json:"cy"`
	R          string `json:"r"`
	Transform  string `json:"transform"`
	XlinkHref  string `json:"xlinkHref"`
	ClipPath   string `json:"clipPath"`
	ClipPathId string `json:"clipPathId"`
	X1         string `json:"x1"`
	Y1         string `json:"y1"`
	X2         string `json:"x2"`
	Y2         string `json:"y2"`
}

type GetMapsDto struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type HighlightedRegionDto struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func GetMaps() ([]GetMapsDto, error) {
	rows, err := Connection.Query("SELECT label, classname from maps;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var maps = []GetMapsDto{}
	for rows.Next() {
		var m GetMapsDto
		if err = rows.Scan(&m.Label, &m.Value); err != nil {
			return nil, err
		}
		maps = append(maps, m)
	}

	return maps, rows.Err()
}

func GetMap(className string) (*MapDto, error) {
	statement := "SELECT * from maps WHERE classname = $1;"
	var m MapDto
	err := Connection.QueryRow(statement, className).Scan(&m.ID, &m.Key, &m.ClassName, &m.Label, &m.ViewBox)
	if err != nil {
		return nil, err
	}

	rows, err := Connection.Query("SELECT e.mapid, t.name, e.elementid, e.name, e.d, e.points, e.x, e.y, e.width, e.height, e.cx, e.cy, e.r, e.transform, e.xlinkhref, e.clippath, e.clippathid, e.x1, e.y1, e.x2, e.y2 FROM mapElements e JOIN mapElementType t ON t.id = e.typeid WHERE e.mapId = $1;", m.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var elements = []MapElementDto{}
	for rows.Next() {
		var e MapElementDto
		if err = rows.Scan(&e.MapID, &e.Type, &e.ID, &e.Name, &e.D, &e.Points, &e.X, &e.Y, &e.Width, &e.Height, &e.Cx, &e.Cy, &e.R, &e.Transform, &e.XlinkHref, &e.ClipPath, &e.ClipPathId, &e.X1, &e.Y1, &e.X2, &e.Y2); err != nil {
			return nil, err
		}
		elements = append(elements, e)
	}

	m.Elements = elements
	return &m, rows.Err()
}

func GetMapHighlightedRegions(className string) ([]HighlightedRegionDto, error) {
	statement := "SELECT id from maps WHERE classname = $1;"
	var mapId int
	err := Connection.QueryRow(statement, className).Scan(&mapId)
	if err != nil {
		return nil, err
	}

	rows, err := Connection.Query("SELECT name, name FROM mapElements WHERE mapId = $1 AND elementId != '';", mapId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var regions = []HighlightedRegionDto{}
	for rows.Next() {
		var region HighlightedRegionDto
		if err = rows.Scan(&region.Label, &region.Value); err != nil {
			return nil, err
		}
		regions = append(regions, region)
	}

	return regions, rows.Err()
}
