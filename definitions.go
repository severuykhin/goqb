package goqb

import "context"

type Executor interface {
	Update(ctx context.Context, where Where, fields FieldMap) error
	Insert(ctx context.Context, fields FieldMap) error
	Delete(ctx context.Context, where Where) error
	Select(fields Fields) SelectBuilder
}

type Rows interface {
	Next() bool
	Scan(dest ...any) error
}

type Query interface {
	ToSql() (string, []interface{}, error)
}
