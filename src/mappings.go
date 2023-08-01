package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

func GetMappingGroups(writer http.ResponseWriter, request *http.Request) {
	mapping, err := repo.GetMappingGroups()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(mapping)
}

func (s *Server) getMappingEntries(writer http.ResponseWriter, request *http.Request) {
	entries, err := repo.GetMappingEntries(mux.Vars(request)["key"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}
	language := request.Header.Get("Content-Language")
	if language != "" && language != "en" {
		translatedEntries := make([]repo.MappingEntryDto, len(entries))
		for index, entry := range entries {
			translatedEntry, err := s.translateEntry(entry, language)
			if err != nil {
				http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
				return
			}

			translatedEntries[index] = translatedEntry
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(translatedEntries)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(entries)
}

func (s *Server) translateEntry(entry repo.MappingEntryDto, language string) (repo.MappingEntryDto, error) {
	name, err := s.ts.TranslateText(language, entry.Name)
	if err != nil {
		return repo.MappingEntryDto{}, err
	}

	svgName, err := s.ts.TranslateText(language, entry.SVGName)
	if err != nil {
		return repo.MappingEntryDto{}, err
	}

	return repo.MappingEntryDto{
		ID:               entry.ID,
		GroupID:          entry.GroupID,
		Name:             name,
		Code:             entry.Code,
		FlagUrl:          entry.FlagUrl,
		SVGName:          svgName,
		AlternativeNames: entry.AlternativeNames,
		Prefixes:         entry.Prefixes,
		Grouping:         entry.Grouping,
	}, nil
}

func GetMappingsWithoutFlags(writer http.ResponseWriter, request *http.Request) {
	keys, err := repo.GetMappingsWithoutFlags()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(keys)
}

func EditMapping(writer http.ResponseWriter, request *http.Request) {
	if code, err := IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var mapping repo.UpdateMappingDto
	err = json.Unmarshal(requestBody, &mapping)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	if err = repo.UpdateMapping(mux.Vars(request)["key"], mapping); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

func DeleteMapping(writer http.ResponseWriter, request *http.Request) {
	if code, err := IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	if err := repo.DeleteMapping(mux.Vars(request)["key"]); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
	}
}
