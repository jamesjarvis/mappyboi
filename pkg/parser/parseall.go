package parser

import (
	"github.com/jamesjarvis/mappyboi/pkg/models"
)

func ParseAll(parsers ...Parser) (*models.Data, error) {
	data := &models.Data{}
	var tempData *models.Data
	var err error
	for _, p := range parsers {
		tempData, err = p.Parse()
		if err != nil {
			return nil, err
		}
		data.GoLocations = append(data.GoLocations, tempData.GoLocations...)
	}
	return data, nil
}
