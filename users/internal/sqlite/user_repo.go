package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/odit-bit/sone/users/internal/domain"
)

var _ domain.UserRepository = (*UserRepo)(nil)

type UserRepo struct {
	tableName string
	db        *sql.DB
}

func New(tableName string, db *sql.DB) (*UserRepo, error) {
	u := &UserRepo{
		tableName: tableName,
		db:        db,
	}
	err := u.migrate(context.Background())
	return u, err
}

func (u *UserRepo) migrate(ctx context.Context) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		hash_password TEXT NOT NULL,
		token TEXT NULL
	);
	`, u.tableName,
	)
	_, err := u.db.ExecContext(ctx, query)
	return err
}

// Find implements domain.UserRepository.
func (u *UserRepo) Find(ctx context.Context, filter domain.FilterOption) (domain.User, error) {

	query := fmt.Sprintf("SELECT id, name, hash_password FROM %s WHERE %s = ? LIMIT 1", u.tableName, filter.Field)

	var user domain.User
	err := u.db.QueryRowContext(ctx, query, filter.Value).Scan(
		&user.ID,
		&user.Name,
		&user.HashPassword,
	)
	return user, err
}

func (u *UserRepo) FindID(ctx context.Context, id string) (domain.User, error) {

	query := fmt.Sprintf("SELECT id, name, hash_password FROM %s WHERE id = ? LIMIT 1", u.tableName)

	var user domain.User
	err := u.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.HashPassword,
	)
	return user, err
}

// Save implements domain.UserRepository.
func (u *UserRepo) Save(ctx context.Context, user *domain.User) error {

	// query := fmt.Sprintf("INSERT INTO %s (id, name, hash_password) Values (?,?,?)", u.tableName)
	query := fmt.Sprintf(`	INSERT INTO %v (id,name,hash_password,token)
		VAlUES (?, ?, ?, ?)
		ON CONFLICT (id) DO 
		UPDATE 
		SET 
			name = excluded.name,
			hash_password = excluded.hash_password,
			token = excluded.token ;
			`, u.tableName)
	_, err := u.db.ExecContext(ctx, query, user.ID, user.Name, user.HashPassword, user.Token)
	return err
}

//	func (u *UserRepo) queryTable(query string) string {
//		return fmt.Sprintf(query, u.tableName)
//	}

// type findQueryBuilder struct {
// 	tableName  string
// 	query      *strings.Builder
// 	fieldCount int
// 	isInit     bool
// 	isFinish   bool
// }

// func (b *findQueryBuilder) Field(field string, value any) {
// 	if b.isInit {
// 		b.query.WriteString("SELECT id, name, hash_password FROM ")
// 		b.query.WriteString(b.tableName)
// 		b.query.WriteString(" WHERE ")
// 		b.isInit = true
// 	}
// 	if b.fieldCount != 0 {
// 		b.query.WriteString(",")
// 	}
// 	b.query.WriteString(field)
// 	b.fieldCount++
// }

// func (b *findQueryBuilder) String() string {
// 	if !b.isFinish {
// 		b.query.WriteString(" = ")
// 		for b.fieldCount != 0 {
// 			b.query.WriteString("?")
// 			b.query.WriteString(",")
// 			b.fieldCount--
// 		}
// 		b.isFinish = true
// 	}
// 	return b.query.String()
// }
