package parser

import (
	"github.com/jamesjarvis/mappyboi/v2/pkg/types"
)

type Parser interface {
	String() string
	Parse() (*types.LocationHistory, error)
}
