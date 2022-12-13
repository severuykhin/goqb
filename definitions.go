package goqb

import "context"

type Executor interface {
	Update(ctx context.Context, where Where, fields FieldMap) error
	Insert(ctx context.Context, fields FieldMap) error
	Delete(ctx context.Context, where Where) error
	Find(ctx context.Context, fields Fields, params FindParams, scanFunc func(*rows) error) error
}
