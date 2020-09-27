package model

type AggregatedStreamResponse struct {
	GroupKey string                 `json:"group_key"`
	Data     map[string]interface{} `json:"data"`
}
