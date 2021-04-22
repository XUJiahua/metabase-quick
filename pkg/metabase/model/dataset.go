package model

import "time"

type DataSetRequest struct {
	Type   string `json:"type"`
	Native struct {
		Query        string `json:"query"`
		TemplateTags struct {
		} `json:"template-tags"`
	} `json:"native"`
	Database   int           `json:"database"`
	Parameters []interface{} `json:"parameters"`
}

//{
//"type": "native",
//"native": {
//"query": "select USER_ID, count(1)\nfrom ORDERS\ngroup by USER_ID\norder by 2 desc\nlimit 10",
//"template-tags": {}
//},
//"database": 1,
//"parameters": []
//}

type DataSetResponse struct {
	Data struct {
		Rows [][]interface{} `json:"rows"`
		Cols []struct {
			DisplayName string   `json:"display_name"`
			Source      string   `json:"source"`
			FieldRef    []string `json:"field_ref"`
			Name        string   `json:"name"`
			BaseType    string   `json:"base_type"`
		} `json:"cols"`
		NativeForm struct {
			Query  string      `json:"query"`
			Params interface{} `json:"params"`
		} `json:"native_form"`
		ResultsTimezone string `json:"results_timezone"`
		ResultsMetadata struct {
			Checksum string `json:"checksum"`
			Columns  []struct {
				Name        string      `json:"name"`
				DisplayName string      `json:"display_name"`
				BaseType    string      `json:"base_type"`
				FieldRef    []string    `json:"field_ref"`
				SpecialType interface{} `json:"special_type"`
				Fingerprint struct {
					Global struct {
						DistinctCount int     `json:"distinct-count"`
						Nil           float64 `json:"nil%"`
					} `json:"global"`
					Type struct {
						TypeNumber struct {
							Min float64 `json:"min"`
							Q1  float64 `json:"q1"`
							Q3  float64 `json:"q3"`
							Max float64 `json:"max"`
							Sd  float64 `json:"sd"`
							Avg float64 `json:"avg"`
						} `json:"type/Number"`
					} `json:"type"`
				} `json:"fingerprint"`
			} `json:"columns"`
		} `json:"results_metadata"`
		Insights interface{} `json:"insights"`
	} `json:"data"`
	DatabaseId int       `json:"database_id"`
	StartedAt  time.Time `json:"started_at"`
	JsonQuery  struct {
		Type   string `json:"type"`
		Native struct {
			Query        string `json:"query"`
			TemplateTags struct {
			} `json:"template-tags"`
		} `json:"native"`
		Database   int `json:"database"`
		Middleware struct {
			JsIntToString                 bool `json:"js-int-to-string?"`
			AddDefaultUserlandConstraints bool `json:"add-default-userland-constraints?"`
		} `json:"middleware"`
	} `json:"json_query"`
	AverageExecutionTime interface{} `json:"average_execution_time"`
	Status               string      `json:"status"`
	Context              string      `json:"context"`
	RowCount             int         `json:"row_count"`
	RunningTime          int         `json:"running_time"`
}
