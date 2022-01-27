package counter

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("no data found")
	ErrBadData  = errors.New("incorrect data format")
)

type Counter interface {
	AddVisit(ctx context.Context, url string, visitorID string) error // Register new visit
	Visits(cxt context.Context, url string) (uint64, error)           // query visits for a page
}
