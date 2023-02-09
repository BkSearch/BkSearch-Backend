package elasticsearch

type DocumentElasticSearchView struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	Type     int     `json:"type"`
	Score    float64 `json:"score"`
	Accepted bool    `json:"Accepted"`
}
