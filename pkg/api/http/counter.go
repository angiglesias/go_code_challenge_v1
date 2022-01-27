package http

import (
	"challenge/pkg/api"
	"challenge/pkg/counter"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type CounterAPIHandler struct {
	counterSvc counter.Counter
}

func NewCounterAPI(counter counter.Counter) *CounterAPIHandler {
	return &CounterAPIHandler{counterSvc: counter}
}

func (ca *CounterAPIHandler) Setup(mux *http.ServeMux) {
	mux.HandleFunc("/visits/new", ca.RegisterVisit)
	mux.HandleFunc("/visits/stats", ca.GetVisits)
}

func (ca *CounterAPIHandler) RegisterVisit(rw http.ResponseWriter, req *http.Request) {
	// Accept only POST requests
	if req.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// check request content
	contenType := req.Header.Get("Content-Type")
	if !strings.Contains(contenType, "application/json") {
		rw.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	// process request body
	var visit api.Visit
	err := json.NewDecoder(req.Body).Decode(&visit)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	// TODO sanitize page url
	if err := ca.counterSvc.AddVisit(req.Context(), visit.Page, visit.VisitorID); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (ca *CounterAPIHandler) GetVisits(rw http.ResponseWriter, req *http.Request) {
	// accept only GET requests
	if req.Method != http.MethodGet {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// parse URL query to get page url
	query := req.URL.Query()
	url := query.Get("url")
	// return an bad request response because no url was received on the query
	if url == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	// recover visits
	visits, err := ca.counterSvc.Visits(req.Context(), url)
	if err != nil {
		if errors.Is(err, counter.ErrNotFound) {
			http.NotFound(rw, req)
			return
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	// build response json
	resp := api.PageVisits{UniqueVisitors: visits}
	// write response
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(resp)
}
