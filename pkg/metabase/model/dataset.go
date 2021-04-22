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
	Data       *Data     `json:"data"`
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
