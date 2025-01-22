package models

type SearchRequest struct {
    Term  string `json:"term"`
    From  int    `json:"from"`
    Size  int    `json:"size"`
    Field string `json:"field"`
}

type ZincSearchQuery struct {
    SearchType  string   `json:"search_type"`
    Query       Query    `json:"query"`
    From        int     `json:"from"`
    MaxResults  int     `json:"max_results"`
    Source      []string `json:"_source"`
}

type Query struct {
    Term  string `json:"term"`
    Field string `json:"field"`
}