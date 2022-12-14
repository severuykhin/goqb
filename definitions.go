package goqb

import "context"

type Executor interface {
	Update(ctx context.Context, where Where, fields FieldMap) error
	Insert(ctx context.Context, fields FieldMap) error
	Delete(ctx context.Context, where Where) error
	Find(ctx context.Context, fields Fields, params FindParams, scanFunc func(Rows) error) error
}

type Rows interface {
	Next() bool
	Scan(dest ...any) error
}
