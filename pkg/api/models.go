package api

type Visit struct {
	Page      string `json:"url"`
	VisitorID string `json:"visitor_id"`
}

type PageVisits struct {
	UniqueVisitors uint64 `json:"unique_visitors"`
}
