package model

type Data struct {
	Rows       [][]interface{} `json:"rows"`
	Cols       []*Column       `json:"cols"`
	NativeForm struct {
		Query  string      `json:"query"`
		Params interface{} `json:"params"`
	} `json:"native_form"`
	ResultsTimezone string `json:"results_timezone"`
	ResultsMetadata struct {
		Checksum string    `json:"checksum"`
		Columns  []*Column `json:"columns"`
	} `json:"results_metadata"`
	Insights interface{} `json:"insights"`
}
