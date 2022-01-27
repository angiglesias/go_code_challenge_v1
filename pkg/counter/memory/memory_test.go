package memory

import (
	"challenge/pkg/counter"
	"context"
	"errors"
	"sync/atomic"
	"testing"
)

const (
	PageUrl        = "/products/42"
	UnknownPageUrl = "/products/43"
	KnownVisitor   = "07b69940-ca58-4ef0-8b5d-1b787d85b919"
	UnknownVisitor = "53847aa8-f5ab-4b8a-844c-f1bb1241ab55"
)

func TestAddVisit(t *testing.T) {
	counter := new(MemoryCounter)

	err := counter.AddVisit(context.TODO(), PageUrl, KnownVisitor)
	if err != nil {
		t.Errorf("AddVisit unexpected error %v", err)
		t.FailNow()
	}
	// check if page entry was created and validate counter
	page, err := counter.loadPage(PageUrl)
	if err != nil {
		t.Errorf("No entry for page '%s' was found", PageUrl)
		t.FailNow()
	}
	count := atomic.LoadUint64(&page.counter)
	if count != 1 {
		t.Errorf("Expected counter with value of 1, got %d", count)
	}

	// repeat operation to check counter is not incremented
	err = counter.AddVisit(context.TODO(), PageUrl, KnownVisitor)
	if err != nil {
		t.Errorf("AddVisit unexpected error %v", err)
		t.FailNow()
	}
	count = atomic.LoadUint64(&page.counter)
	if count != 1 {
		t.Errorf("Expected counter with value of 1, got %d", count)
	}

	// repeat operation to check counter is incremented with new visitor
	err = counter.AddVisit(context.TODO(), PageUrl, UnknownVisitor)
	if err != nil {
		t.Errorf("AddVisit unexpected error %v", err)
		t.FailNow()
	}
	count = atomic.LoadUint64(&page.counter)
	if count != 2 {
		t.Errorf("Expected counter with value of 2, got %d", count)
	}
}

func TestVisits(t *testing.T) {
	memcounter := new(MemoryCounter)
	// register page
	page := new(pageVisits)
	// register visitor manually
	page.visitors.Store(KnownVisitor, nil)
	atomic.StoreUint64(&page.counter, 1)
	memcounter.pages.Store(PageUrl, page)

	// query visits for page
	visits, err := memcounter.Visits(context.TODO(), PageUrl)
	if err != nil {
		t.Errorf("Unexpected error querying visits for page '%s': %v", PageUrl, err)
		t.FailNow()
	}
	if visits != 1 {
		t.Errorf("Expected counter with value of 1, got %d", visits)
	}

	// query visits for unkown page
	_, err = memcounter.Visits(context.TODO(), UnknownPageUrl)
	if err == nil {
		t.Error("Expected a failure querying an unknown page")
		t.FailNow()
	}
	if !errors.Is(err, counter.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, go instead: %v", err)
	}

}
