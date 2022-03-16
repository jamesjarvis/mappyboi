package parser

import (
	"fmt"

	"log"

	"github.com/jamesjarvis/mappyboi/pkg/models"
)

func ParseAll(parsers ...Parser) (*models.Data, error) {
	data := &models.Data{}
	var tempData *models.Data
	var err error
	for _, p := range parsers {
		tempData, err = p.Parse()
		if err != nil {
			return nil, fmt.Errorf("failed to parse '%s': %w", p.String(), err)
		}
		data.GoLocations = append(data.GoLocations, tempData.GoLocations...)
		log.Printf("%s parsed, %d/%d points...\n", p.String(), len(tempData.GoLocations), len(data.GoLocations))
	}
	return data, nil
}
