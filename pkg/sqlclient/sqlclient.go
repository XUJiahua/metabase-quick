package sqlclient

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xujiahua/metabase-quick/pkg/metabase/model"
	"reflect"
)

type Client struct {
	db *sql.DB
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
		db: db,
	}, nil
}

// reference metabase code
//(defmethod sql-jdbc.sync/database-type->base-type :mysql
//    :BIGINT     :type/BigInteger
//    :INT        :type/Integer
//    :INTEGER    :type/Integer
//    :MEDIUMINT  :type/Integer
//    :TINYINT    :type/Integer
//    :SMALLINT   :type/Integer
//    :BIT        :type/Boolean
//    :DECIMAL    :type/Decimal
//    :NUMERIC    :type/Decimal
//    :DOUBLE     :type/Float
//    :REAL       :type/Float
//    :FLOAT      :type/Float
//    :CHAR       :type/Text
//    :LONGTEXT   :type/Text
//    :TEXT       :type/Text
//    :TINYTEXT   :type/Text
//    :VARCHAR    :type/Text
//    :MEDIUMTEXT :type/Text
//    :BINARY     :type/*
//    :BLOB       :type/*
//    :LONGBLOB   :type/*
//    :MEDIUMBLOB :type/*
//    :ENUM       :type/*
//    :SET        :type/*
//    :TINYBLOB   :type/*
//    :VARBINARY  :type/*
//    :DATE       :type/Date
//    :YEAR       :type/Date
//    :TIMESTAMP  :type/DateTimeWithLocalTZ ; stored as UTC in the database
//    :DATETIME   :type/DateTime
//    :TIME       :type/Time

// mysql datatype in go driver
// https://github.com/go-sql-driver/mysql/blob/46351a8892976898935f653f5333782579a96fa5/fields.go#L16
// TODO: not complete
func toMetabaseType(typ string) string {
	switch typ {
	case "SMALLINT", "INT", "TINYINT", "MEDIUMINT":
		return "type/Integer"
	case "BIGINT":
		return "type/BigInteger"
	case "VARCHAR", "NVARCHAR", "VARCHAR2", "CHAR", "TEXT":
		return "type/Text"
	case "DECIMAL":
		return "type/Decimal"
	case "DOUBLE", "FLOAT":
		return "type/Float"
	case "BOOL":
		return "type/Boolean"
	default:
		return "type/*"
	}

	return ""
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
			DisplayName: columnType.Name(),
			BaseType:    toMetabaseType(columnType.DatabaseTypeName()),
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
