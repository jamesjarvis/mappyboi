package google

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jamesjarvis/mappyboi/v2/pkg/types"
)

type LocationHistory struct {
	Filepath string
}

func (p *LocationHistory) String() string {
	return p.Filepath
}

func (p *LocationHistory) Parse() (types.LocationHistory, error) {
	var data GoogleData

	file, err := os.Open(p.Filepath)
	if err != nil {
		return types.LocationHistory{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return types.LocationHistory{}, fmt.Errorf("failed to decode json file '%s': %w", p.Filepath, err)
	}

	return GoogleDataToData(data)
}
