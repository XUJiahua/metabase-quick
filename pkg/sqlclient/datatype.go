package sqlclient

// mysql datatype in go driver:
// https://github.com/go-sql-driver/mysql/blob/master/fields.go#L16-L98
// extract via https://regex101.com/ using regex /\".*\"/gm
//"BIT"
//"TEXT"
//"BLOB"
//"DATE"
//"DATETIME"
//"DECIMAL"
//"DOUBLE"
//"ENUM"
//"FLOAT"
//"GEOMETRY"
//"MEDIUMINT"
//"JSON"
//"INT"
//"LONGTEXT"
//"LONGBLOB"
//"BIGINT"
//"MEDIUMTEXT"
//"MEDIUMBLOB"
//"NULL"
//"SET"
//"SMALLINT"
//"BINARY"
//"CHAR"
//"TIME"
//"TIMESTAMP"
//"TINYINT"
//"TINYTEXT"
//"TINYBLOB"
//"VARBINARY"
//"VARCHAR"
//"VARBINARY"
//"VARCHAR"
//"YEAR"

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

func toMetabaseDataType(mysqlDataType string) string {
	switch mysqlDataType {
	case "SMALLINT", "INT", "TINYINT", "MEDIUMINT", "INTEGER":
		return "type/Integer"
	case "BIGINT":
		return "type/BigInteger"
	case "VARCHAR", "CHAR", "TEXT", "LONGTEXT", "TINYTEXT", "MEDIUMTEXT":
		return "type/Text"
	case "DECIMAL", "NUMERIC":
		return "type/Decimal"
	case "DOUBLE", "FLOAT", "REAL":
		return "type/Float"
	case "BIT":
		return "type/Boolean"
	case "DATE", "YEAR":
		return "type/Date"
	case "TIMESTAMP":
		return "type/DateTimeWithLocalTZ"
	case "DATETIME":
		return "type/DateTime"
	case "TIME":
		return "type/Time"
	}

	return "type/*"
}
