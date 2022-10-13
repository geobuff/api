package repo

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

type MapElementDto struct {
	EntryID    int    `json:"entryId"`
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

type CreateMapElementDto struct {
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

type UpdateMapElementDto struct {
	Name      string `json:"name"`
	ElementID string `json:"elementId"`
}

func GetMapElements(mapId int) ([]MapElementDto, error) {
	rows, err := Connection.Query("SELECT e.id, e.mapid, t.name, e.elementid, e.name, e.d, e.points, e.x, e.y, e.width, e.height, e.cx, e.cy, e.r, e.transform, e.xlinkhref, e.clippath, e.clippathid, e.x1, e.y1, e.x2, e.y2 FROM mapElements e JOIN mapElementType t ON t.id = e.typeid WHERE e.mapId = $1;", mapId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var elements = []MapElementDto{}
	for rows.Next() {
		var e MapElementDto
		if err = rows.Scan(&e.EntryID, &e.MapID, &e.Type, &e.ID, &e.Name, &e.D, &e.Points, &e.X, &e.Y, &e.Width, &e.Height, &e.Cx, &e.Cy, &e.R, &e.Transform, &e.XlinkHref, &e.ClipPath, &e.ClipPathId, &e.X1, &e.Y1, &e.X2, &e.Y2); err != nil {
			return nil, err
		}
		elements = append(elements, e)
	}

	return elements, rows.Err()
}

func GetHighlightedElements(mapId int) ([]HighlightedRegionDto, error) {
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

func CreateMapElement(mapId int, element MapElementDto) error {
	typeId, err := GetMapElementTypeId(element.Type)
	if err != nil {
		return err
	}

	var id int
	statement := "INSERT INTO mapelements (mapid, typeid, elementid, name, d, points, x, y, width, height, cx, cy, r, transform, xlinkhref, clippath, clippathid, x1, y1, x2, y2) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21) RETURNING id;"
	return Connection.QueryRow(statement, mapId, typeId, element.ID, element.Name, element.D, element.Points, element.X, element.Y, element.Width, element.Height, element.Cx, element.Cy, element.R, element.Transform, element.XlinkHref, element.ClipPath, element.ClipPathId, element.X1, element.Y1, element.X2, element.Y2).Scan(&id)
}

func DeleteMapElements(mapId int) error {
	var id int
	return Connection.QueryRow("DELETE FROM mapelements where mapid = $1 RETURNING id;", mapId).Scan(&id)
}

func UpdateMapElement(entryID int, entry UpdateMapElementDto) error {
	var id int
	return Connection.QueryRow("UPDATE mapelements SET name = $2, elementid = $3 WHERE id = $1 RETURNING id;", entryID, entry.Name, entry.ElementID).Scan(&id)
}
