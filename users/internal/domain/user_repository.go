package domain

import "context"

// type FilterOption interface {
// 	mustImplemented()
// }

// type TokenFilter struct {
// 	Value string
// }

// func (tf *TokenFilter) mustImplemented() {}

type FilterOption struct {
	Field string
	Value any
}

const Id_Field = "id"
const Token_Field = "token"
const Name_Field = "name"

type UserRepository interface {
	Save(ctx context.Context, user *User) error
	Find(ctx context.Context, filter FilterOption) (User, error)
	FindID(ctx context.Context, id string) (User, error)
}
