package goqb

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type goqb struct {
	connection *sqlx.DB
	tableName  string
}

type Where map[string]interface{}

type Fields []string

type FieldMap map[string]interface{}

type FindParams struct {
	Where   Where
	Limit   int
	Offset  int
	OrderBy string
}

func New(connection *sqlx.DB, tableName string) Executor {
	return &goqb{
		connection: connection,
		tableName:  tableName,
	}
}

func (qb *goqb) Update(
	ctx context.Context,
	where Where,
	fields FieldMap,
) error {

	query := sq.
		Update(qb.tableName).
		SetMap(fields).
		Where(where)

	sql, values, err := query.ToSql()
	if err != nil {
		return err
	}

	sql = qb.connection.Rebind(sql)

	res, err := qb.connection.Exec(sql, values...)

	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	// if rowsAffected <= 0 {
	// 	return errors.ErrInternal.
	// 		WithMessage(fmt.Sprintf("db:no rows affected while updating %s", tableName))

	// }

	return nil
}

func (qb *goqb) Insert(
	ctx context.Context,
	fields FieldMap,
) error {

	columns := []string{}
	values := []interface{}{}

	for column, value := range fields {
		columns = append(columns, column)
		values = append(values, value)
	}

	query := sq.Insert(qb.tableName).Columns(columns...).Values(values...)

	sql, values, err := query.ToSql()
	if err != nil {
		return err
	}

	sql = qb.connection.Rebind(sql)

	_, err = qb.connection.Exec(sql, values...)

	if err != nil {
		return err
	}

	return nil
}

func (qb *goqb) Delete(
	ctx context.Context,
	where Where,
) error {
	query := sq.
		Delete(qb.tableName).
		Where(where)

	sql, values, err := query.ToSql()
	if err != nil {
		return err
	}

	sql = qb.connection.Rebind(sql)

	_, err = qb.connection.Exec(sql, values...)

	if err != nil {
		return err
	}

	return nil
}

func (qb *goqb) Find(
	ctx context.Context,
	fields Fields,
	params FindParams,
	scanFunc func(Rows) error,
) error {

	query := sq.Select(fields...).From(qb.tableName)

	if params.Where != nil {
		query = query.Where(params.Where)
	}

	if params.Limit > 0 {
		query = query.Limit(uint64(params.Limit))
	}

	if params.Offset > 0 {
		query = query.Offset(uint64(params.Offset))
	}

	if len(params.OrderBy) > 0 {
		query = query.OrderBy(params.OrderBy)
	}

	sqlQuery, values, err := query.ToSql()
	if err != nil {
		return err
	}

	sqlQuery = qb.connection.Rebind(sqlQuery)

	sqlRows, err := qb.connection.Query(sqlQuery, values...)

	defer func() {
		if sqlRows != nil {
			err := sqlRows.Close()
			if err != nil {
				fmt.Println("QB ERR: ", err)
				// c.logger.Error(ErrCodeCloseRowsError, err.Error(), "sql", sqlQuery)
			}
		} else {
			fmt.Println("QB ERR: nil rows")
			// c.logger.Error(ErrCodeFetchRowsError, err.Error(), "sql", sqlQuery)
		}
	}()

	if err != nil {
		return err
	}

	err = scanFunc(&rows{sqlRows: sqlRows})

	if err != nil {
		return err
	}

	return nil
}
