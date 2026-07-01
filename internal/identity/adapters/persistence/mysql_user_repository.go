package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
)

const duplicateEntryErrorNumber uint16 = 1062

var (
	// ErrDuplicateUserEmail is returned when a user email already exists.
	ErrDuplicateUserEmail = errors.New("duplicate user email")

	// ErrUserNotFound is returned when a user cannot be found.
	ErrUserNotFound = errors.New("user not found")
)

// MySQLUserRepository persists identity users in MySQL.
type MySQLUserRepository struct {
	db *sql.DB
}

// NewMySQLUserRepository constructs a MySQL-backed user repository.
func NewMySQLUserRepository(db *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{db: db}
}

// ExistsByEmail reports whether a user exists for an email address.
func (r *MySQLUserRepository) ExistsByEmail(ctx context.Context, email valueobject.Email) (bool, error) {
	const query = `
SELECT EXISTS(
    SELECT 1
    FROM identity_users
    WHERE email = ?
)`

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, email.String()).Scan(&exists); err != nil {
		return false, fmt.Errorf("query user exists by email: %w", err)
	}

	return exists, nil
}

// FindByEmail retrieves a user by email address.
func (r *MySQLUserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*domain.User, error) {
	const query = `
SELECT id, email, password_hash, registered_at
FROM identity_users
WHERE email = ?`

	var row userRow
	err := r.db.QueryRowContext(ctx, query, email.String()).Scan(
		&row.id,
		&row.email,
		&row.passwordHash,
		&row.registeredAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("query user by email: %w", err)
	}

	user, err := row.toDomain()
	if err != nil {
		return nil, fmt.Errorf("map user row to domain: %w", err)
	}

	return user, nil
}

// Save persists a new user.
func (r *MySQLUserRepository) Save(ctx context.Context, user *domain.User) error {
	const query = `
INSERT INTO identity_users (
    id,
    email,
    password_hash,
    registered_at
) VALUES (?, ?, ?, ?)`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID().String(),
		user.Email().String(),
		user.PasswordHash().String(),
		user.RegisteredAt(),
	)
	if isDuplicateEntry(err) {
		return ErrDuplicateUserEmail
	}

	if err != nil {
		return fmt.Errorf("insert identity user: %w", err)
	}

	return nil
}

func isDuplicateEntry(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == duplicateEntryErrorNumber
}
