package repos

import (
	"backend/models"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetMe(ctx context.Context, userID uuid.UUID) (models.User, error)
	RegisterUser(ctx context.Context, data *models.User) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	GetUser(ctx context.Context, email string) (models.User, error)
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error)
	GetUserByVerificationToken(ctx context.Context, vCode string) (models.User, error)
	GetUserByPasswordResetToken(ctx context.Context, resetToken string) (models.User, error)
	UpdateUser(ctx context.Context, data *models.User) (*models.User, error)
	PatchUserVerification(ctx context.Context, userID uuid.UUID, fields map[string]interface{}) (*models.User, error)
	PatchResetPassword(ctx context.Context, user *models.User, fields map[string]interface{}) (*models.User, error)
	SetPasswordResetToken(ctx context.Context, userID uuid.UUID, token string, expiry *time.Time) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	// PatchUser(ctx context.Context, userID uuid.UUID, fields map[string]interface{}) (*models.User, error)
	// GetTicketByID(ctx context.Context, ticketID uuid.UUID) (*Ticket, error)
	// GetTickets(ctx context.Context) ([]*Ticket, error)
	// PatchTicket(ctx context.Context, ticketID uuid.UUID, fields map[string]interface{}) (*Ticket, error)
	// DeleteTicket(ctx context.Context, ticketID uuid.UUID) (uuid.UUID, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetMe(ctx context.Context, userID uuid.UUID) (models.User, error) {

	var user models.User

	query := `SELECT * FROM users WHERE user_id = ?`

	rows, err := r.db.Raw(query, userID).Rows()
	if err != nil {
		return models.User{}, fmt.Errorf("error getting user: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		err = r.db.ScanRows(rows, &user)
		if err != nil {
			return models.User{}, fmt.Errorf("error scanning user: %v", err)
		}
	}

	return user, nil
}

func (r *UserRepository) RegisterUser(ctx context.Context, data *models.User) (models.User, error) {

	var user models.User

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query := `
	WITH new_user AS (
		INSERT INTO users (first_name, last_name, email, password, photo_url)
		VALUES (?, ?, ?, ?, ?)
		RETURNING *
	),
	new_user_verification AS (
		INSERT INTO user_verifications (user_id, verification_token, expires_on)
		VALUES ((SELECT user_id FROM new_user), ?, ?)
		RETURNING *
	)
	SELECT * FROM new_user LEFT JOIN new_user_verification ON new_user.user_id = new_user_verification.user_id;
	`

	rows, err := tx.Raw(query, data.FirstName, data.LastName, data.Email, data.Password, data.PhotoURL, data.Verification.VerificationToken, data.Verification.ExpiresOn).Rows()
	if err != nil {
		tx.Rollback()
		return models.User{}, fmt.Errorf("error while registering user: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = r.db.ScanRows(rows, &data)
		if err != nil {
			tx.Rollback()
			return models.User{}, fmt.Errorf("error while scanning registered user: %v", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return models.User{}, fmt.Errorf("error while committing transaction: %v", err)
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {

	user := models.User{}

	query := `SELECT u.user_id, u.first_name, u.last_name, u.email, u.password, u.role, u.provider, u.photo_url, u.created_on, u.modified_on, uv.* FROM users u LEFT JOIN user_verifications uv ON u.user_id = uv.user_id WHERE email = ?;`

	rows, err := r.db.Raw(query, email).Rows()
	if err != nil {
		return models.User{}, fmt.Errorf("error getting user: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		verification := &models.UserVerification{}
		err := rows.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role, &user.Provider, &user.PhotoURL, &user.CreatedOn, &user.ModifiedOn, &verification.ID, &verification.UserID, &verification.VerificationToken, &verification.ExpiresOn, &verification.IsVerified, &verification.VerifiedOn, &verification.CreatedOn, &verification.ModifiedOn)
		if err != nil {
			return models.User{}, fmt.Errorf("error scanning user: %v", err)
		}
		user.Verification = verification
	}

	fmt.Println("WHAT THE FUCK?: ", user)

	return user, nil
}

// func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {

// 	user := &models.User{}

// 	query := `
// 	SELECT * FROM users LEFT JOIN user_verifications ON users.user_id = user_verifications.user_id WHERE email = ?;
// 	`

// 	rows, err := r.db.Raw(query, email).Rows()
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting user: %v", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		verification := &models.UserVerification{}
// 		err := rows.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role, &user.Provider, &user.PhotoURL, &user.CreatedOn, &user.ModifiedOn, &user.LastAccessAt, &verification.ID, &verification.UserID, &verification.VerificationToken, &verification.IsVerified, &verification.ExpiresOn, &verification.VerifiedOn, &verification.CreatedOn, &verification.ModifiedOn)
// 		if err != nil {
// 			return nil, fmt.Errorf("error scanning user: %v", err)
// 		}
// 	}

// 	fmt.Println("WHAT THE FUCK?: ", user)

// 	return user, nil
// }

func (r *UserRepository) GetUser(ctx context.Context, email string) (models.User, error) {

	var user models.User

	query := `
	SELECT users.user_id, users.first_name, users.last_name, users.email, users.password, users.role, users.provider, users.photo_url, users.created_on, users.modified_on, user_verifications.id, user_verifications.user_id, user_verifications.verification_token, user_verifications.is_verified, user_verifications.expires_on, user_verifications.verified_on, user_verifications.created_on, user_verifications.modified_on
	FROM users 
	LEFT JOIN user_verifications ON users.user_id = user_verifications.user_id WHERE email = ?;
	`

	rows, err := r.db.Raw(query, email).Rows()
	if err != nil {
		return models.User{}, fmt.Errorf("error while getting user: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var verification = &models.UserVerification{}
		err := rows.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Role, &user.Provider, &user.PhotoURL, &user.CreatedOn, &user.ModifiedOn, &verification.ID, &verification.UserID, &verification.VerificationToken, &verification.IsVerified, &verification.ExpiresOn, &verification.VerifiedOn, &verification.CreatedOn, &verification.ModifiedOn)
		if err != nil {
			return models.User{}, fmt.Errorf("error while scanning user: %v", err)
		}
		user.Verification = verification
	}

	fmt.Println("WHAT THE FUCK?: ", user)

	return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error) {

	var user models.User

	query := `
	SELECT * FROM users WHERE user_id = ?
	`

	if err := r.db.Raw(query, userID).Scan(&user).Error; err != nil {
		return models.User{}, fmt.Errorf("error while getting user by id: %v", err)
	}

	if user == (models.User{}) {
		return models.User{}, fmt.Errorf("user not found")
	}

	return user, nil
}

func (r *UserRepository) GetUsers(ctx context.Context) ([]models.User, error) {

	var users []models.User

	query := `SELECT * FROM users`

	if err := r.db.Raw(query).Scan(&users).Error; err != nil {
		return []models.User{}, fmt.Errorf("error while getting user by id: %v", err)
	}

	if users == nil {
		return []models.User{}, fmt.Errorf("users not found")
	}

	return users, nil
}

func (r *UserRepository) GetUserByVerificationToken(ctx context.Context, vCode string) (models.User, error) {

	var user models.User

	query := `
	SELECT users.user_id, first_name, last_name, email, password, photo_url, users.created_on, users.modified_on, user_verifications.user_id, verification_token, is_verified, expires_on, user_verifications.created_on, user_verifications.modified_on
	FROM users LEFT JOIN user_verifications ON users.user_id = user_verifications.user_id WHERE verification_token = ?;
	`

	rows, err := r.db.Raw(query, vCode).Rows()
	if err != nil {
		return models.User{}, fmt.Errorf("error while getting user by verification code: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var verification = &models.UserVerification{}
		err := rows.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.PhotoURL, &user.CreatedOn, &user.ModifiedOn, &verification.UserID, &verification.VerificationToken, &verification.IsVerified, &verification.ExpiresOn, &verification.CreatedOn, &verification.ModifiedOn)
		if err != nil {
			return models.User{}, fmt.Errorf("error while scanning user in get with verification code")
		}

		user.Verification = verification
	}

	// if user == (models.User{}) {
	// 	return models.User{}, fmt.Errorf("user not found")
	// }

	return user, nil
}

func (r *UserRepository) GetUserByPasswordResetToken(ctx context.Context, resetToken string) (models.User, error) {

	var user models.User

	query := `
	SELECT u.user_id, u.first_name, u.last_name, u.email, up.id, up.user_id, up.pass_reset_token, up.expires_on, up.pass_reset_on, up.created_on, up.modified_on
	FROM users u LEFT JOIN user_password_resets up ON u.user_id = up.user_id WHERE up.pass_reset_token = ?;
	`

	rows, err := r.db.Raw(query, resetToken).Rows()
	if err != nil {
		return models.User{}, fmt.Errorf("error while getting user by verification code: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var userPasswordReset = &models.UserPasswordReset{}
		err := rows.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &userPasswordReset.ID, &userPasswordReset.UserID, &userPasswordReset.PassResetToken, &userPasswordReset.ExpiresOn, &userPasswordReset.PassResetOn, &userPasswordReset.CreatedOn, &userPasswordReset.ModifiedOn)
		if err != nil {
			return models.User{}, fmt.Errorf("error while scanning user in get with verification code")
		}

		user.PasswordReset = userPasswordReset
	}

	if user == (models.User{}) {
		return models.User{}, fmt.Errorf("invalid password reset token or token has expired. please request a new one")
	}

	return user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, data *models.User) (*models.User, error) {

	var user *models.User

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query := `
	WITH updated_user AS (
		UPDATE users SET first_name = ?, last_name = ?, email = ?, password = ?, photo_url = ?, modified_on = Now() WHERE user_id = ?
		RETURNING *
	),
	updated_user_verification AS (
		UPDATE user_verifications SET verification_token = ?, is_verified = ?, expires_on = ? WHERE user_id = (SELECT user_id FROM updated_user)
		RETURNING *
	)
	SELECT * FROM updated_user LEFT JOIN updated_user_verification ON updated_user.user_id = updated_user_verification.user_id;
	`

	rows, err := tx.Raw(query, data.FirstName, data.LastName, data.Email, data.Password, data.PhotoURL, data.UserID, data.Verification.VerificationToken, data.Verification.IsVerified, data.Verification.ExpiresOn).Rows()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error while updating user: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = r.db.ScanRows(rows, &data)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error while scanning updated user: %v", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error while committing transaction: %v", err)
	}

	return user, nil
}

func (r *UserRepository) PatchResetPassword(ctx context.Context, data *models.User, fields map[string]interface{}) (*models.User, error) {

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query := `UPDATE users SET password = ? WHERE user_id = ? RETURNING *`
	rows, err := tx.Raw(query, data.Password, data.UserID).Rows()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error while patching user password: %v", err)
	}
	defer rows.Close()

	user := &models.User{}
	for rows.Next() {
		err = r.db.ScanRows(rows, &user)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error while scanning user in patch password: %v", err)
		}
	}

	resetQuery := `UPDATE user_password_resets SET`
	resetParams := []interface{}{}
	validResetKeys := map[string]bool{
		"password_reset":   true,
		"pass_reset_token": true,
		"expires_on":       true,
		"pass_reset_on":    true,
	}

	for k := range fields {
		if _, ok := validResetKeys[k]; !ok {
			fmt.Printf("invalid key: %s", k)
			continue
		}
		if k == "password_reset" {
			for _, v := range fields[k].([]interface{}) {
				fmt.Println("v: ", v)
				passresetMap := v.(map[string]interface{})
				fmt.Println("passresetMap: ", passresetMap)
				numResetFields := len(passresetMap)
				resetI := 1
				for k, v := range passresetMap {
					if _, ok := validResetKeys[k]; !ok {
						fmt.Printf("invalid key: %s", k)
						continue
					}
					resetQuery += fmt.Sprintf(" %s = ?", k)
					resetParams = append(resetParams, v)
					if resetI < numResetFields {
						resetQuery += `,`
					}
					resetI++
					fmt.Println("myresetquery: ", resetQuery)
					fmt.Println(k, v)
				}
			}
		} else {
			resetQuery += fmt.Sprintf(" %s = ?,", k)
			resetParams = append(resetParams, fields[k])
		}
	}
	resetQuery = strings.TrimSuffix(resetQuery, ", ")
	resetQuery += ` modified_on = ? WHERE user_id = ? RETURNING *`
	resetParams = append(resetParams, time.Now(), data.UserID)

	myrows, err := tx.Raw(resetQuery, resetParams...).Rows()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error while patching user password reset: %v", err)
	}
	defer myrows.Close()

	userPasswordReset := &models.UserPasswordReset{}
	if myrows.Next() {
		err = tx.ScanRows(myrows, &userPasswordReset)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error while scanning user in patch password reset: %v", err)
		}
	}

	user.PasswordReset = userPasswordReset

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error while committing transaction: %v", err)
	}

	return user, nil
}

func (r *UserRepository) SetPasswordResetToken(ctx context.Context, userID uuid.UUID, token string, expiry *time.Time) error {

	var user models.User

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	fmt.Println("expiry", expiry)
	fmt.Println("token", token)
	fmt.Println("UserID", userID)

	query := `
	INSERT INTO user_password_resets (user_id, pass_reset_token, expires_on, modified_on) 
	VALUES (?, ?, ?, ?)
	`

	rows, err := tx.Raw(query, userID, token, expiry, time.Now()).Rows()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error while setting password reset token: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		err = r.db.ScanRows(rows, &user)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error while scanning user in set password reset token: %v", err)
		}
	}

	fmt.Println("setPasswordResetToken-User", user)

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error while committing transaction: %v", err)
	}

	return nil
}

func (r *UserRepository) PatchUserVerification(ctx context.Context, userID uuid.UUID, fields map[string]interface{}) (*models.User, error) {

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query := `UPDATE user_verifications SET`
	params := []interface{}{}
	validKeys := map[string]bool{
		"user_id":            true,
		"verification_token": true,
		"expires_on":         true,
		"is_verified":        true,
		"verified_on":        true,
	}
	for k, v := range fields {
		if !validKeys[k] {
			fmt.Printf("invalid field: %s", k)
			continue
		}
		query += fmt.Sprintf(" %s = ?,", k)
		params = append(params, v)
	}
	query = strings.TrimSuffix(query, ",")
	query += ", modified_on = ? WHERE user_id = ? RETURNING *"

	params = append(params, time.Now())
	params = append(params, userID)

	rows, err := tx.Raw(query, params...).Rows()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error while updating user: %v", err)
	}
	defer rows.Close()

	verification := &models.UserVerification{}
	for rows.Next() {
		err = tx.ScanRows(rows, &verification)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error while scanning updated user: %v", err)
		}
	}

	user := &models.User{}
	user.Verification = verification

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error while committing transaction: %v", err)
	}

	return user, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {

	query := `
	DELETE FROM users WHERE user_id = ?;
	`

	if err := r.db.Raw(query, userID).Error; err != nil {
		return fmt.Errorf("error while deleting user: %v", err)
	}

	return nil
}
