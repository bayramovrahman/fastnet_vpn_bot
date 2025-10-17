package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/bayramovrahman/fastnet_vpn_bot/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

func (m *postgresDBRepo) GetUserById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, username, first_name, last_name, email, password, is_verified, is_admin, access_level, signup_ip, signup_country, created_at, updated_at
						from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.IsVerified,
		&user.IsAdmin,
		&user.AccessLevel,
		&user.SignupIP,
		&user.SignupCountry,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

// UpdateUser updates a user in the database
func (m *postgresDBRepo) UpdateUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5`

	_, err := m.DB.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.AccessLevel,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// Authenticate authenticates a user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, password from users where email = $1`

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

// GetUserLoginSecurity gets user login security settings by user ID
func (m *postgresDBRepo) GetUserLoginSecurity(userID int) (models.UserLoginSecurity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, user_id, email_verification, phone_verification, multi_factor_auth, 
			  COALESCE(verification_code, ''), COALESCE(code_expires_at, '0001-01-01'), 
			  COALESCE(phone_number, ''), COALESCE(last_verification_sent_at, '0001-01-01'),
			  failed_attempts, COALESCE(locked_until, '0001-01-01'), created_at, updated_at
			  FROM user_login_security WHERE user_id = $1`

	row := m.DB.QueryRowContext(ctx, query, userID)

	var security models.UserLoginSecurity
	err := row.Scan(
		&security.ID,
		&security.UserID,
		&security.EmailVerification,
		&security.PhoneVerification,
		&security.MultiFactorAuth,
		&security.VerificationCode,
		&security.CodeExpiresAt,
		&security.PhoneNumber,
		&security.LastVerificationSentAt,
		&security.FailedAttempts,
		&security.LockedUntil,
		&security.CreatedAt,
		&security.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		security = models.UserLoginSecurity{
			UserID:            userID,
			EmailVerification: false,
			PhoneVerification: false,
			MultiFactorAuth:   false,
			FailedAttempts:    0,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		err = m.CreateUserLoginSecurity(security)
		if err != nil {
			return security, err
		}
		
		return m.GetUserLoginSecurity(userID)
	}

	if err != nil {
		return security, err
	}

	return security, nil
}

// UpdateUserLoginSecurity updates user login security settings
func (m *postgresDBRepo) UpdateUserLoginSecurity(security models.UserLoginSecurity) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE user_login_security 
			  SET email_verification = $1, phone_verification = $2, multi_factor_auth = $3,
			      verification_code = $4, code_expires_at = $5, phone_number = $6,
			      last_verification_sent_at = $7, failed_attempts = $8, locked_until = $9,
			      updated_at = $10
			  WHERE user_id = $11`

	_, err := m.DB.ExecContext(ctx, query,
		security.EmailVerification,
		security.PhoneVerification,
		security.MultiFactorAuth,
		security.VerificationCode,
		security.CodeExpiresAt,
		security.PhoneNumber,
		security.LastVerificationSentAt,
		security.FailedAttempts,
		security.LockedUntil,
		time.Now(),
		security.UserID,
	)

	return err
}

// CreateUserLoginSecurity creates new user login security settings
func (m *postgresDBRepo) CreateUserLoginSecurity(security models.UserLoginSecurity) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO user_login_security 
						(user_id, email_verification, phone_verification, multi_factor_auth, failed_attempts, created_at, updated_at)
			  		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, query,
		security.UserID,
		security.EmailVerification,
		security.PhoneVerification,
		security.MultiFactorAuth,
		security.FailedAttempts,
		time.Now(),
		time.Now(),
	)

	return err
}
