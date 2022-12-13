package goqb

import "database/sql"

type rows struct {
	sqlRows *sql.Rows
}

func newRows(sqlRows *sql.Rows) *rows {
	return &rows{
		sqlRows: sqlRows,
	}
}

func (r *rows) Next() bool {
	return r.sqlRows.Next()
}

func (r *rows) Scan(dest ...any) error {
	return r.sqlRows.Scan(dest...)
}
