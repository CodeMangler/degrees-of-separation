package moviebuff

import (
	"../graph"
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "http://data.moviebuff.com/"

var httpClient = &http.Client{}

type mbEntity struct {
	URL    string         `json:"url"`
	Name   string         `json:"name"`
	Type   string         `json:"type"`
	Movies []mbConnection `json:"movies,omitempty"`
	Cast   []mbConnection `json:"cast,omitempty"`
}

type mbConnection struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Role string `json:"type"`
}

var entities = make(map[string]*mbEntity)

func fetchEntity(id string) (*mbEntity, error) {
	entityURL := baseURL + id

	if entity, ok := entities[id]; ok {
		fmt.Printf("----------->>>>> Found in Cache: %v\n", id)
		return entity, nil
	}
	fmt.Printf("----------- Fetching: %v\n", entityURL)
	response, errHTTP := httpClient.Get(entityURL)
	if errHTTP != nil {
		return nil, errHTTP
	}

	entity := &mbEntity{}
	dec := json.NewDecoder(response.Body)
	errDecode := dec.Decode(&entity)
	if errDecode != nil {
		return nil, errDecode
	}
	entities[id] = entity
	return entity, nil
}

// Fetch fetches moviebuff content given an ID/URL, and populates Neighbours of the Node.
func Fetch(n *graph.Node) {
	entity, _ := fetchEntity(n.ID)
	var connections []mbConnection
	if entity.Type == "Person" {
		connections = entity.Movies
	} else {
		connections = entity.Cast
	}

	for _, connection := range connections {
		n.Connect(graph.NewNode(connection.URL, Fetch))
	}
}
