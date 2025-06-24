package timeline

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jamesjarvis/mappyboi/v2/pkg/types"
)

type Timeline struct {
	Filepath string
}

func (p *Timeline) String() string {
	return p.Filepath
}

func (p *Timeline) Parse() (*types.LocationHistory, error) {
	var data []GoogleLocation
	fmt.Println("Parsing Google Timeline data from", p.Filepath)

	file, err := os.Open(p.Filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode json file '%s': %w", p.Filepath, err)
	}

	return timelineDataToData(data)
}
