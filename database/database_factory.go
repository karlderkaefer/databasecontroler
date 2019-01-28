package database

import (
	"fmt"
)

const (
	Oracle11 = 11
	Oracle12 = 12
	MySQL = 57
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
	default:
		return 0
	}
}
