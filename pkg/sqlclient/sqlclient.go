package sqlclient

import (
	"fmt"
	"github.com/xujiahua/metabase-quick/pkg/metabase/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}

func New(dbAddr, user, pass, dbName string) (*Client, error) {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	//dsn := "root:@tcp(127.0.0.1:3306)/default?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, dbAddr, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Client{
		db: db,
	}, nil
}

func (c Client) RowsAndCols(query string) ([][]interface{}, []*model.Column, error) {
	rows, err := c.db.Raw(query).Rows()
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, nil, err
	}
	var columns []*model.Column
	for _, columnType := range columnTypes {
		columns = append(columns, &model.Column{
			Name:        columnType.Name(),
			DisplayName: columnType.Name(),
			BaseType:    columnType.DatabaseTypeName(),
		})
	}

	nCol := len(columnTypes)
	var _rows [][]interface{}
	for rows.Next() {
		values := make([]interface{}, nCol)
		valuePtrs := make([]interface{}, nCol)
		for i := 0; i < nCol; i++ {
			valuePtrs[i] = &values[i]
		}
		err = rows.Scan(valuePtrs...)
		if err != nil {
			return nil, nil, err
		}

		_rows = append(_rows, values)
	}

	return _rows, columns, nil
}
