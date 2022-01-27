package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"challenge/pkg/api"
	"challenge/pkg/counter"
	"challenge/pkg/logging"
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
	logging.Debugf("[API:RegisterVisit] Processing new request")
	// Accept only POST requests
	if req.Method != http.MethodPost {
		logging.Errorf("[API:RegisterVisit] Method %s is not allowed for this operation", req.Method)
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// check request content
	contenType := req.Header.Get("Content-Type")
	if !strings.Contains(contenType, "application/json") {
		logging.Errorf("[API:RegisterVisit] Content-Type needs to be 'application/json'")
		rw.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	// process request body
	var visit api.Visit
	err := json.NewDecoder(req.Body).Decode(&visit)
	if err != nil {
		logging.Errorf("[API:RegisterVisit] Error decoding request payload: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	// TODO sanitize page url
	if err := ca.counterSvc.AddVisit(req.Context(), visit.Page, visit.VisitorID); err != nil {
		logging.Errorf("[API:RegisterVisit] Error counting visit to page %s: %v", visit.Page, err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (ca *CounterAPIHandler) GetVisits(rw http.ResponseWriter, req *http.Request) {
	logging.Debugf("[API:GetVisits] Processing new request")
	// accept only GET requests
	if req.Method != http.MethodGet {
		logging.Errorf("[API:GetVisits] Method %s is not allowed for this operation", req.Method)
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// parse URL query to get page url
	query := req.URL.Query()
	url := query.Get("url")
	// return an bad request response because no url was received on the query
	if url == "" {
		logging.Errorf("[API:GetVisits] Page URL is empty")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	// recover visits
	visits, err := ca.counterSvc.Visits(req.Context(), url)
	if err != nil {
		if errors.Is(err, counter.ErrNotFound) {
			logging.Warnf("[API:GetVisits] No data available for page '%s'", url)
			http.NotFound(rw, req)
			return
		} else {
			logging.Errorf("[API:GetVisits] Error consulting data for page '%s'", url)
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
