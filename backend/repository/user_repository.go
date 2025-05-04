package repository

import (
	"database/sql"
	"errors"
	"time"

	"library-management-system/internal/models"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *Database
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *Database) *UserRepository {
	return &UserRepository{db: db}
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, phone, address, 
		       created_at, updated_at, last_login, account_status, 
		       failed_login_attempts, two_factor_enabled
		FROM users 
		WHERE id = ?`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.Role, &user.Phone, &user.Address, &user.CreatedAt,
		&user.UpdatedAt, &user.LastLogin, &user.AccountStatus,
		&user.FailedLoginAttempts, &user.TwoFactorEnabled,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, phone, address, 
		       created_at, updated_at, last_login, account_status, 
		       failed_login_attempts, two_factor_enabled
		FROM users 
		WHERE email = ?`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.Role, &user.Phone, &user.Address, &user.CreatedAt,
		&user.UpdatedAt, &user.LastLogin, &user.AccountStatus,
		&user.FailedLoginAttempts, &user.TwoFactorEnabled,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// Create adds a new user to the database
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (
			email, password_hash, full_name, role, phone, address, 
			account_status, two_factor_enabled
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(
		query,
		user.Email, user.PasswordHash, user.FullName, user.Role,
		user.Phone, user.Address, user.AccountStatus, user.TwoFactorEnabled,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

// Update updates an existing user
func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET email = ?, full_name = ?, role = ?, phone = ?, address = ?,
			account_status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err := r.db.Exec(
		query,
		user.Email, user.FullName, user.Role, user.Phone,
		user.Address, user.AccountStatus, user.ID,
	)
	return err
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(userID int64, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err := r.db.Exec(query, passwordHash, userID)
	return err
}

// UpdateLoginAttempts updates a user's failed login attempts
func (r *UserRepository) UpdateLoginAttempts(userID int64, attempts int) error {
	query := `
		UPDATE users
		SET failed_login_attempts = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err := r.db.Exec(query, attempts, userID)
	return err
}

// UpdateLastLogin updates a user's last login timestamp
func (r *UserRepository) UpdateLastLogin(userID int64) error {
	query := `
		UPDATE users
		SET last_login = CURRENT_TIMESTAMP, failed_login_attempts = 0, 
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err := r.db.Exec(query, userID)
	return err
}

// Delete removes a user from the database
func (r *UserRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// List retrieves users with pagination
func (r *UserRepository) List(page, pageSize int) ([]*models.User, error) {
	// Safeguard against pagination parameter manipulation
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	query := `
		SELECT id, email, full_name, role, phone, address, 
		       created_at, updated_at, last_login, account_status
		FROM users
		ORDER BY id DESC
		LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		var lastLogin sql.NullTime

		err := rows.Scan(
			&user.ID, &user.Email, &user.FullName, &user.Role,
			&user.Phone, &user.Address, &user.CreatedAt, &user.UpdatedAt,
			&lastLogin, &user.AccountStatus,
		)
		if err != nil {
			return nil, err
		}

		if lastLogin.Valid {
			user.LastLogin = &lastLogin.Time
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Count returns the total number of users
func (r *UserRepository) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users`
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// SetResetToken sets a password reset token for a user
func (r *UserRepository) SetResetToken(userID int64, token string, expiresAt time.Time) error {
	query := `
		UPDATE users
		SET reset_token = ?, reset_token_expires = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err := r.db.Exec(query, token, expiresAt, userID)
	return err
}

// GetByResetToken retrieves a user by reset token
func (r *UserRepository) GetByResetToken(token string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, 
		       reset_token_expires, account_status
		FROM users
		WHERE reset_token = ? AND reset_token_expires > CURRENT_TIMESTAMP
		AND account_status = 'active'`

	var user models.User
	var expires time.Time

	err := r.db.QueryRow(query, token).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.Role, &expires, &user.AccountStatus,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrInvalidResetToken
		}
		return nil, err
	}

	return &user, nil
}

// ClearResetToken clears a user's reset token
func (r *UserRepository) ClearResetToken(userID int64) error {
	query := `
		UPDATE users
		SET reset_token = NULL, reset_token_expires = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err := r.db.Exec(query, userID)
	return err
}

// SetTwoFactorSecret sets a user's two-factor authentication secret
func (r *UserRepository) SetTwoFactorSecret(userID int64, secret string, enabled bool) error {
	query := `
		UPDATE users
		SET two_factor_secret = ?, two_factor_enabled = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`

	_, err := r.db.Exec(query, secret, enabled, userID)
	return err
}

// GetTwoFactorSecret retrieves a user's two-factor authentication secret
func (r *UserRepository) GetTwoFactorSecret(userID int64) (string, error) {
	query := `SELECT two_factor_secret FROM users WHERE id = ?`

	var secret sql.NullString
	err := r.db.QueryRow(query, userID).Scan(&secret)
	if err != nil {
		return "", err
	}

	if !secret.Valid {
		return "", models.ErrTwoFactorNotEnabled
	}

	return secret.String, nil
}
