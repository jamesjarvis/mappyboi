package polarsteps

import (
	"encoding/json"
	"os"

	"github.com/jamesjarvis/mappyboi/v2/pkg/types"
)

type PolarstepLocationFile struct {
	Filepath string
}

func (p *PolarstepLocationFile) String() string {
	return p.Filepath
}

func (p *PolarstepLocationFile) Parse() (*types.LocationHistory, error) {
	var data PolarstepData

	file, err := os.Open(p.Filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return polarstepDataToLocationHistory(data)
}
