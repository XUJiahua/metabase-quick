package model

type Database struct {
	Description              interface{} `json:"description"`
	Features                 []string    `json:"features"`
	CacheFieldValuesSchedule string      `json:"cache_field_values_schedule"`
	Timezone                 string      `json:"timezone"`
	AutoRunQueries           bool        `json:"auto_run_queries"`
	MetadataSyncSchedule     string      `json:"metadata_sync_schedule"`
	Name                     string      `json:"name"`
	Caveats                  interface{} `json:"caveats"`
	Tables                   []*Table    `json:"tables"`
	IsFullSync               bool        `json:"is_full_sync"`
	UpdatedAt                string      `json:"updated_at"`
	NativePermissions        string      `json:"native_permissions"`
	Details                  struct {
		Db string `json:"db"`
	} `json:"details"`
	IsSample         bool        `json:"is_sample"`
	Id               int         `json:"id"`
	IsOnDemand       bool        `json:"is_on_demand"`
	Options          interface{} `json:"options"`
	Engine           string      `json:"engine"`
	Refingerprint    interface{} `json:"refingerprint"`
	CreatedAt        string      `json:"created_at"`
	PointsOfInterest interface{} `json:"points_of_interest"`
}

//{
//"description": null,
//"features": [
//"basic-aggregations",
//"standard-deviation-aggregations",
//"expression-aggregations",
//"foreign-keys",
//"right-join",
//"left-join",
//"native-parameters",
//"nested-queries",
//"expressions",
//"case-sensitivity-string-filter-options",
//"binning",
//"inner-join",
//"advanced-math-expressions"
//],
//"cache_field_values_schedule": "0 0 9 * * ? *",
//"timezone": "UTC",
//"auto_run_queries": true,
//"metadata_sync_schedule": "0 38 * * * ? *",
//"name": "Sample Dataset",
//"caveats": null,
//"is_full_sync": true,
//"updated_at": "2021-04-21T02:43:59.167",
//"native_permissions": "write",
//"details": {
//"db": "/Users/jiahua/opensource/metabase/resources/sample-dataset.db;USER=GUEST;PASSWORD=guest"
//},
//"is_sample": true,
//"id": 1,
//"is_on_demand": false,
//"options": null,
//"engine": "h2",
//"refingerprint": null,
//"created_at": "2021-04-20T05:02:06.893",
//"points_of_interest": null
//}

func NewDatabase(name string, id int) *Database {
	return &Database{
		Features: []string{
			"basic-aggregations",
			"standard-deviation-aggregations",
			"expression-aggregations",
			"foreign-keys",
			"right-join",
			"left-join",
			"native-parameters",
			"nested-queries",
			"expressions",
			"case-sensitivity-string-filter-options",
			"binning",
			"inner-join",
			"advanced-math-expressions",
		},
		Name:              name,
		Id:                id,
		NativePermissions: "write",
		IsSample:          true,
	}
}

func (d *Database) AddTable(tableName string, id int) {
	t := NewTable(tableName, id, d)
	d.Tables = append(d.Tables, t)
}
