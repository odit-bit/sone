package domain

import "context"

type UserGoogleRepository interface {
	Get(ctx context.Context, id string) (UserGl, bool, error)
	Save(ctx context.Context, user UserGl) error
}
