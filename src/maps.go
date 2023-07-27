package src

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

type SvgDto struct {
	SVG string `json:"svg"`
}

type CreateMapDto struct {
	SVGMap   repo.MapDto            `json:"svgMap"`
	Mappings repo.CreateMappingsDto `json:"mappings"`
	Quiz     repo.CreateQuizDto     `json:"quiz"`
}

func GetMaps(writer http.ResponseWriter, request *http.Request) {
	maps, err := repo.GetMaps()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(maps)
}

func GetMap(writer http.ResponseWriter, request *http.Request) {
	svgMap, err := repo.GetMap(mux.Vars(request)["className"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(svgMap)
}

func GetMapHighlightedRegions(writer http.ResponseWriter, request *http.Request) {
	regions, err := repo.GetMapHighlightedRegions(mux.Vars(request)["className"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(regions)
}

func GetMapPreview(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var svgDto SvgDto
	err = json.Unmarshal(requestBody, &svgDto)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	preview, err := getMapPreview(svgDto.SVG)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(preview)
}

func getMapPreview(svg string) (repo.MapDto, error) {
	result := repo.MapDto{
		ID:        0,
		Key:       "preview",
		ClassName: "Preview",
		Label:     "Map Preview",
		ViewBox:   "",
		Elements:  []repo.MapElementDto{},
	}

	var width string
	var height string
	var currentElement repo.MapElementDto
	scanner := bufio.NewScanner(strings.NewReader(svg))
	for scanner.Scan() {
		text := scanner.Text()

		if strings.Contains(text, "width=") {
			widthSub := text[(strings.Index(text, `width="`) + 7):]
			width = widthSub[:strings.Index(widthSub, `"`)]
		}

		if strings.Contains(text, "height=") {
			heightSub := text[(strings.Index(text, `height="`) + 8):]
			height = heightSub[:strings.Index(heightSub, `"`)]
		}

		if strings.Contains(text, "id=") {
			idSub := text[(strings.Index(text, `id="`) + 4):]
			currentElement.ID = idSub[:strings.Index(idSub, `"`)]
			currentElement.Type = "path"
			result.Elements = append(result.Elements, currentElement)
		} else if strings.Contains(text, "d=") {
			sub := text[(strings.Index(text, `d="`) + 3):]
			currentElement.D = sub[:strings.Index(sub, `"`)]
		} else if strings.Contains(text, "title=") {
			titleSub := text[(strings.Index(text, `title="`) + 7):]
			currentElement.Name = titleSub[:strings.Index(titleSub, `"`)]
		}
	}
	result.ViewBox = fmt.Sprintf("0 0 %s %s", width, height)

	if err := scanner.Err(); err != nil {
		return repo.MapDto{}, err
	}

	sort.Slice(result.Elements, func(i, j int) bool {
		return result.Elements[i].Name < result.Elements[j].Name
	})

	return result, nil
}

func CreateMap(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var payload CreateMapDto
	err = json.Unmarshal(requestBody, &payload)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	err = repo.CreateMap(payload.SVGMap)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	err = repo.CreateMappings(payload.Mappings)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	_, err = repo.CreateQuiz(payload.Quiz)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}
