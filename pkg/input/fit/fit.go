package fit

import (
	"os"

	"github.com/jamesjarvis/mappyboi/v2/pkg/types"
	"github.com/tormoder/fit"
)

type FitFile struct {
	Filepath   string
	Compressed bool
}

func (p *FitFile) String() string {
	return p.Filepath
}

func (p *FitFile) Parse() (*types.LocationHistory, error) {
	file, err := os.Open(p.Filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode .fit file.
	fitFile, err := fit.Decode(file)
	if err != nil {
		return nil, err
	}

	// Get the actual activity
	activity, err := fitFile.Activity()
	if err != nil {
		return nil, err
	}

	// Add to data object.
	data := &types.LocationHistory{
		Data: make([]types.Location, 0, len(activity.Records)),
	}

	for _, record := range activity.Records {
		data.Insert(fitRecordToGoLocation(record))
	}

	return data, nil
}
