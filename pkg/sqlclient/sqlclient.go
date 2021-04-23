package sqlclient

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xujiahua/metabase-quick/pkg/metabase/model"
	"reflect"
	"strings"
)

type Client struct {
	db     *sql.DB
	DBName string
}

func New(dbAddr, user, pass, dbName string) (*Client, error) {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	//dsn := "root:@tcp(127.0.0.1:3306)/default?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, dbAddr, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &Client{
		db:     db,
		DBName: dbName,
	}, nil
}

func (c Client) Query(query string) ([][]interface{}, []*model.Column, error) {
	rows, err := c.db.Query(query)
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
			DisplayName: strings.ToUpper(columnType.Name()),
			BaseType:    toMetabaseDataType(columnType.DatabaseTypeName()),
			Source:      "native",
		})
	}

	nCol := len(columnTypes)
	var _rows [][]interface{}
	for rows.Next() {
		values := make([]interface{}, nCol)
		for i, columnType := range columnTypes {
			// ScanType() supported in mysql:
			// https://github.com/go-sql-driver/mysql/blob/master/fields.go#L101
			values[i] = reflect.New(columnType.ScanType()).Interface()
		}

		err = rows.Scan(values...)
		if err != nil {
			return nil, nil, err
		}

		for i := range columnTypes {
			// unwrap the value if type is sql.NullBool/...
			if valuer, ok := values[i].(driver.Valuer); ok {
				values[i], _ = valuer.Value()
			} else if bs, ok := values[i].(*sql.RawBytes); ok {
				// nilness check
				if bs == nil {
					values[i] = nil
				} else {
					// treat as string
					values[i] = string(*bs)
				}
			}
			// otherwise, keep values[i] with no change
		}

		_rows = append(_rows, values)
	}

	return _rows, columns, nil
}

func (c Client) GetTables() ([]string, error) {
	rows, cols, err := c.Query("show tables;")
	if err != nil {
		return nil, err
	}

	if len(cols) != 1 {
		return nil, errors.New("unexpect to happen")
	}

	var tables []string
	for _, row := range rows {
		if table, ok := row[0].(string); ok {
			tables = append(tables, table)
		} else {
			return nil, errors.New("unexpect to happen")
		}
	}
	return tables, nil
}
