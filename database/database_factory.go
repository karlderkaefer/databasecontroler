package database

import "fmt"

const (
	oracle11      = 11
	oracle12      = 12
	mysql         = 57
	sqlserver2017 = 17
	db2105        = 105
)

// GetDatabase is used get the API of the requested database
func GetDatabase(m int) (Database, error) {
	switch m {
	case oracle11:
		db := new(oracleDatabase)
		db.version = oracle11
		return db, nil
	case oracle12:
		db := new(oracleDatabase)
		db.version = oracle12
		return db, nil
	case mysql:
		return new(mysqlDatabase), nil
	case sqlserver2017:
		return new(sqlserverDatabase), nil
	case db2105:
		return new(db2), nil
	default:
		return nil, fmt.Errorf("database %d not recognized", m)
	}
}

// ParseVersion will parse the input and return a suitable database.
// returns 0 if database is not registered.
func ParseVersion(db string) int {
	switch db {
	case "oracle11":
		return oracle11
	case "oracle12":
		return oracle12
	case "mysql":
		return mysql
	case "sqlserver2017":
		return sqlserver2017
	case "db2":
		return db2105
	default:
		return 0
	}
}
