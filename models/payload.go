package models

type CrawlerPayload struct {
	URL      string   `json:"url"`
	Keywords []string `json:"keywords"`
	Depth    int32    `json:"depth"`
}

type CrawlerRequest struct {
	Requests []CrawlerPayload `json:"requests"`
}
