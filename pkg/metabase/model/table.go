package model

type Table struct {
	Description          string      `json:"description"`
	EntityType           string      `json:"entity_type"`
	Schema               string      `json:"schema"`
	ShowInGettingStarted bool        `json:"show_in_getting_started"`
	Name                 string      `json:"name"`
	Caveats              interface{} `json:"caveats"`
	UpdatedAt            string      `json:"updated_at"`
	EntityName           interface{} `json:"entity_name"`
	Active               bool        `json:"active"`
	Id                   int         `json:"id"`
	DbId                 int         `json:"db_id"`
	VisibilityType       interface{} `json:"visibility_type"`
	FieldOrder           string      `json:"field_order"`
	DisplayName          string      `json:"display_name"`
	CreatedAt            string      `json:"created_at"`
	PointsOfInterest     interface{} `json:"points_of_interest"`
}

// sample data
//{
//"description": "This is a confirmed order for a product from a user.",
//"entity_type": "entity/TransactionTable",
//"schema": "PUBLIC",
//"show_in_getting_started": false,
//"name": "ORDERS",
//"caveats": null,
//"updated_at": "2021-04-22T06:38:03.046",
//"entity_name": null,
//"active": true,
//"id": 2,
//"db_id": 1,
//"visibility_type": null,
//"field_order": "database",
//"display_name": "Orders",
//"created_at": "2021-04-20T05:02:07.618",
//"points_of_interest": null
//}

func NewTable(tableName string, id int, db *Database) *Table {
	return &Table{
		EntityType:  "entity/TransactionTable",
		Schema:      db.Name,
		DbId:        db.Id,
		Name:        tableName,
		Id:          id,
		DisplayName: tableName,
		FieldOrder:  "database",
		Active:      true,
	}
}
