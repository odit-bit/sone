package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/odit-bit/sone/users/internal/domain"
)

var _ domain.UserGoogleRepository = (*userGlRepo)(nil)

type userGlRepo struct {
	tableName string
	db        *sql.DB
}

func NewUserGlRepo(ctx context.Context, db *sql.DB, tableName string) (*userGlRepo, error) {
	repo := userGlRepo{
		tableName: tableName,
		db:        db,
	}
	if err := repo.migrate(ctx); err != nil {
		return nil, err
	}
	return &repo, nil
}

func (u *userGlRepo) migrate(ctx context.Context) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		is_email_verified INTEGER NOT NULL
	);
	`, u.tableName,
	)
	_, err := u.db.ExecContext(ctx, query)
	return err
}

// Get implements domain.UserGoogleRepository.
func (u *userGlRepo) Get(ctx context.Context, id string) (domain.UserGl, bool, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = ? LIMIT 1 ;`, u.tableName)
	var user domain.UserGl
	var verified int
	row := u.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&verified,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, false, nil
		}
		return user, false, err
	}
	if verified == 1 {
		user.IsEmailVerified = true
	} else {
		user.IsEmailVerified = false
	}
	return user, true, nil

}

// Save implements domain.UserGoogleRepository.
func (u *userGlRepo) Save(ctx context.Context, user domain.UserGl) error {

	query := fmt.Sprintf(`
		INSERT INTO %v (id,name,email,is_email_verified)
		VAlUES (?, ?, ?, ?)
		ON CONFLICT (id) DO 
		UPDATE 
		SET 
			name = excluded.name,
			email = excluded.email,
			is_email_verified = excluded.is_email_verified ;
	`, u.tableName,
	)

	var verified int
	if user.IsEmailVerified {
		verified = 1
	}
	_, err := u.db.ExecContext(ctx, query, user.Id, user.Name, user.Email, verified)
	return err
}
