package parser

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"log"

	"github.com/jamesjarvis/mappyboi/pkg/models"
	pool "github.com/jamesjarvis/massivelyconcurrentsystems/pool"
)

type parserWork struct {
	parser Parser
	wg     *sync.WaitGroup

	Error error
	Data  *models.Data
}

func ParseAll(parsers ...Parser) (*models.Data, error) {
	data := &models.Data{}

	workerFunc := func(p *parserWork) error {
		defer p.wg.Done()
		tempData, err := p.parser.Parse()
		if err != nil {
			p.Error = fmt.Errorf("failed to parse '%s': %w", p.parser.String(), err)
			return p.Error
		}

		p.Data = tempData

		// log.Printf("Parsed %d points from %s...\n", len(tempData.GoLocations), p.parser.String())
		return nil
	}

	numWorkers := runtime.NumCPU() - 2

	dispatcher := pool.NewSingleDispatcher(workerFunc, pool.NewConfig(pool.SetNumConsumers(numWorkers)))
	dispatcher.Start()
	defer dispatcher.Close()

	work := []*parserWork{}

	var wg sync.WaitGroup

	for _, p := range parsers {
		wg.Add(1)
		pw := &parserWork{
			parser: p,
			wg:     &wg,
		}
		work = append(work, pw)
		err := dispatcher.Put(context.Background(), pw)
		if err != nil {
			return nil, fmt.Errorf("failed to schedule parser '%s': %w", p.String(), err)
		}
	}

	wg.Wait()

	// Collect results from work.
	for _, pw := range work {
		if pw.Error != nil {
			return nil, pw.Error
		}

		data.GoLocations = append(data.GoLocations, pw.Data.GoLocations...)
	}

	log.Printf("Collected %d points from %d parsers\n", len(data.GoLocations), len(work))

	return data, nil
}
