package pkg

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/go-gota/gota/dataframe"
	"os"
)

func Load(filename string, hasHeader bool) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	dataFrame := dataframe.ReadCSV(file, dataframe.HasHeader(hasHeader))
	spew.Dump(dataFrame)
}
