package metabase

import "time"

type DataSetRequest struct {
	Database int `json:"database"`
	Query    struct {
		SourceTable int             `json:"source-table"`
		Aggregation [][]string      `json:"aggregation"`
		Breakout    [][]interface{} `json:"breakout"`
	} `json:"query"`
	Type       string        `json:"type"`
	Parameters []interface{} `json:"parameters"`
}

type DataSetResponse struct {
	Data struct {
		Rows [][]interface{} `json:"rows"`
		Cols []struct {
			Description    string        `json:"description"`
			TableId        int           `json:"table_id"`
			SpecialType    *string       `json:"special_type"`
			Name           string        `json:"name"`
			Settings       interface{}   `json:"settings"`
			Source         string        `json:"source"`
			FieldRef       []interface{} `json:"field_ref"`
			ParentId       interface{}   `json:"parent_id"`
			Id             int           `json:"id"`
			VisibilityType string        `json:"visibility_type"`
			DisplayName    string        `json:"display_name"`
			Fingerprint    *struct {
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
					} `json:"type/Number,omitempty"`
					TypeDateTime struct {
						Earliest time.Time `json:"earliest"`
						Latest   time.Time `json:"latest"`
					} `json:"type/DateTime,omitempty"`
				} `json:"type,omitempty"`
			} `json:"fingerprint"`
			BaseType string `json:"base_type"`
			Unit     string `json:"unit,omitempty"`
		} `json:"cols"`
		NativeForm struct {
			Query  string      `json:"query"`
			Params interface{} `json:"params"`
		} `json:"native_form"`
		ResultsTimezone string `json:"results_timezone"`
		RowsTruncated   int    `json:"rows_truncated"`
		ResultsMetadata struct {
			Checksum string `json:"checksum"`
			Columns  []struct {
				Name        string        `json:"name"`
				DisplayName string        `json:"display_name"`
				BaseType    string        `json:"base_type"`
				SpecialType *string       `json:"special_type"`
				Id          int           `json:"id"`
				FieldRef    []interface{} `json:"field_ref"`
				Description string        `json:"description"`
				Fingerprint *struct {
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
						} `json:"type/Number,omitempty"`
						TypeDateTime struct {
							Earliest time.Time `json:"earliest"`
							Latest   time.Time `json:"latest"`
						} `json:"type/DateTime,omitempty"`
					} `json:"type,omitempty"`
				} `json:"fingerprint"`
				Unit string `json:"unit,omitempty"`
			} `json:"columns"`
		} `json:"results_metadata"`
		Insights interface{} `json:"insights"`
	} `json:"data"`
	DatabaseId int       `json:"database_id"`
	StartedAt  time.Time `json:"started_at"`
	JsonQuery  struct {
		Database int `json:"database"`
		Query    struct {
			SourceTable int `json:"source-table"`
		} `json:"query"`
		Type       string `json:"type"`
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
