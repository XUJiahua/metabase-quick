package sqldb

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestNew(t *testing.T) {
	s, err := New("localhost:3306")
	assert.Equal(t, nil, err)
	s.ImportTable()
	err = s.Start()
	assert.Equal(t, nil, err)
}
