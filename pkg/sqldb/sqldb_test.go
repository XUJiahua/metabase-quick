package sqldb

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestNew(t *testing.T) {
	s, err := New("localhost:3306")
	assert.Equal(t, nil, err)
	_, err = s.ImportTable("/Users/jiahua/goworkshop/metabase-quick/dataset/sample-dataset/orders.csv", true)
	assert.Equal(t, nil, err)
	err = s.Start()
	assert.Equal(t, nil, err)
}
