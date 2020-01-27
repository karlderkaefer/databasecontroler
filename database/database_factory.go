package database

import (
	"fmt"
)

const (
	Oracle11      = 11
	Oracle12      = 12
	MySQL         = 57
	SqlServer2017 = 17
	Db2105        = 105
)

func GetDatabase(m int) (Database, error) {
	switch m {
	case Oracle11:
		db := new(Oracle)
		db.version = Oracle11
		return db, nil
	case Oracle12:
		db := new(Oracle)
		db.version = Oracle12
		return db, nil
	case MySQL:
		return new(Mysql), nil
	case SqlServer2017:
		return new(Sqlserver), nil
	case Db2105:
		return new(Db2), nil
	default:
		return nil, fmt.Errorf("Database %d not recognized\n", m)
	}
}

func ParseVersion(db string) int {
	switch db {
	case "oracle11":
		return Oracle11
	case "oracle12":
		return Oracle12
	case "mysql":
		return MySQL
	case "sqlserver2017":
		return SqlServer2017
	case "db2":
		return Db2105
	default:
		return 0
	}
}
