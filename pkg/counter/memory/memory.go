package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"challenge/pkg/counter"
)

type MemoryCounter struct {
	pages sync.Map // using sync map as keys will only be written once and read multiple times on multiple goroutines
}

func NewCounter() counter.Counter { return new(MemoryCounter) }

// AddVisit register the visit to a page
func (mc *MemoryCounter) AddVisit(ctx context.Context, url, visitorID string) error {
	// context is ignored as this implementation is non-blocking
	// search if page is available
	page, err := mc.loadPage(url)
	if err != nil {
		if errors.Is(err, counter.ErrNotFound) {
			// create and register new page counter for unknown page
			page = new(pageVisits)
			mc.pages.Store(url, page)
		} else {
			return fmt.Errorf("error registering visit to page %s", url)
		}
	}
	// Add vistor to page
	page.addVisit(visitorID)
	return nil
}

// Visits will query the distinct visits to a page
func (mc *MemoryCounter) Visits(ctx context.Context, url string) (uint64, error) {
	// context is ignored as this implementation is non-blocking
	page, err := mc.loadPage(url)
	if err != nil {
		return 0, fmt.Errorf("error recovering visits for page %s: %w", url, err)
	}
	return page.visits(), nil
}

func (mc *MemoryCounter) loadPage(url string) (*pageVisits, error) {
	loaded, ok := mc.pages.Load(url)
	if !ok {
		return nil, counter.ErrNotFound
	}
	page, ok := loaded.(*pageVisits)
	// just in case, in this code path this failure condition should newver happen
	if !ok {
		return nil, counter.ErrBadData
	}
	return page, nil
}

// pageVisits is an auxiliar structure containing the counter and unique visitor for one page
type pageVisits struct {
	counter  uint64   // using an uint64 to count unique visitors in a thread safe manner with sync/atomic functions
	visitors sync.Map // using sync map as keys will only be written once and read in multiple goroutines times to filter previous visits
}

func (pv *pageVisits) addVisit(visitorID string) {
	// save visitor id to avoid count multiple times the same visitor ID
	_, loaded := pv.visitors.LoadOrStore(visitorID, nil)
	if !loaded {
		// new visitor received, incrementing counter
		atomic.AddUint64(&pv.counter, 1)
	}
}

func (pv *pageVisits) visits() uint64 {
	return atomic.LoadUint64(&pv.counter)
}
