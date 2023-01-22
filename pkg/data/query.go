package data

import (
	"database/sql"

	"github.com/pkg/errors"
)

func Query() ([]string, error) {
	if db == nil {
		return nil, errors.New("database not initialized")
	}

	stmt, err := db.Prepare("SELECT subject FROM doc ORDER BY subject")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to prepare select statement")
	}

	rows, err := stmt.Query()
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrapf(err, "failed to execute select statement")
	}
	defer rows.Close()

	list := make([]string, 0)
	for rows.Next() {
		var val string
		if err := rows.Scan(&val); err != nil {
			return nil, errors.Wrapf(err, "failed to scan row")
		}
		list = append(list, val)
	}

	return list, nil
}
