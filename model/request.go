package model

type AggregatedCreateRequest struct {
	AggregatedID       string   `json:"aggregated_id"`
	GroupByKeyID       []string `json:"group_by_key_id"`
	AggregatedFunction []struct {
		PropertyName string `json:"property_name"`
		FieldName    string `json:"field_name"`
		Function     string `json:"function"`
	} `json:"aggregated_function"`
}
