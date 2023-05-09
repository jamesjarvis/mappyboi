package parser

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"log"

	"github.com/jamesjarvis/mappyboi/v2/pkg/types"
	pool "github.com/jamesjarvis/massivelyconcurrentsystems/pool"
)

type parserWork struct {
	parser Parser
	wg     *sync.WaitGroup

	Error error
	Data  types.LocationHistory
}

func ParseAll(parsers ...Parser) (types.LocationHistory, error) {
	data := types.LocationHistory{}

	workerFunc := func(p *parserWork) error {
		defer p.wg.Done()
		tempData, err := p.parser.Parse()
		if err != nil {
			p.Error = fmt.Errorf("failed to parse '%s': %w", p.parser.String(), err)
			return p.Error
		}

		p.Data = tempData

		log.Printf("Parsed %d points from %s...\n", len(tempData.Data), p.parser.String())
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
			return types.LocationHistory{}, fmt.Errorf("failed to schedule parser '%s': %w", p.String(), err)
		}
	}

	wg.Wait()

	// Collect results from work.
	for _, pw := range work {
		if pw.Error != nil {
			return types.LocationHistory{}, pw.Error
		}

		data.Insert(pw.Data.Data...)
	}

	log.Printf("Collected %d points from %d parsers\n", len(data.Data), len(work))

	return data, nil
}
