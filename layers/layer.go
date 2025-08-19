package layers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/charmbracelet/log"
)

type Team struct {
	Commander      string
	Faction        string
	Name           string
	NumHelicopters int
	NumTanks       int
	Tickets        int
}

type Layer struct {
	ClassName    string `json:"levelName"`
	GameMode     string `json:"gamemode"`
	GameModeType string `json:"type"`
	LayerId      string `json:"rawName"`
	Name         string `json:"Name"`
	Size         string `json:"mapSize"`
	Version      string `json:"layerVersion"`
}

type MapResponse struct {
	Maps []Layer `json:"Maps"`
}

type Layers struct {
	Layers []Layer
	pulled bool
}

func New() *Layers {
	return &Layers{
		pulled: false,
	}
}

func (layers *Layers) FetchLayers(force bool) {
	if layers.pulled && !force {
		log.Warn("Layer refresh request but layers already fetched.")
		return
	}

	if force {
		log.Warn("Force updating layer information.")
	}

	resp, err := http.Get("https://raw.githubusercontent.com/Squad-Wiki/squad-wiki-pipeline-map-data/master/completed_output/_Current%20Version/finished.json")

	if err != nil {
		log.Error("There was an error loading layers: " + err.Error())
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Error("Error reading response: " + err.Error())
		return
	}

	var responseData MapResponse

	err = json.Unmarshal(body, &responseData)

	if err != nil {
		log.Errorf("An error occurred parsing the JSON: %s", err.Error())
		return
	}

	layers.pulled = true
}

func (layers *Layers) FilterLayers(predicate func(Layer) bool) []Layer {
	result := make([]Layer, 0, len(layers.Layers))

	for _, layer := range layers.Layers {
		if predicate(layer) {
			result = append(result, layer)
		}
	}

	return result
}

func (layers *Layers) GetLayerById(id string) Layer {
	foundLayers := layers.FilterLayers(func(layer Layer) bool {
		return layer.LayerId == id
	})

	if len(foundLayers) == 0 {
		return Layer{}
	}

	return foundLayers[0]
}
