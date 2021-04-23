package sqlclient

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestNew(t *testing.T) {
	client, err := New("127.0.0.1:3306", "root", "", "default")
	assert.Equal(t, nil, err)
	rows, columns, err := client.Query("show tables;")
	//rows, columns, err := client.Query("select * from ORDERS limit 1;")
	//rows, columns, err := client.Query("select USER_ID, count(1) from ORDERS group by USER_ID order by 2 desc limit 10;")
	assert.Equal(t, nil, err)
	spew.Dump(rows)
	spew.Dump(columns)
}

func TestClient_GetTables(t *testing.T) {
	client, err := New("127.0.0.1:3306", "root", "", "default")
	assert.Equal(t, nil, err)
	tables, err := client.GetTables()
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(tables))
	spew.Dump(tables)
}
