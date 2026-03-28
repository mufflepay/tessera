package controllers

import (
	"tessera/config"
	"tessera/models"
	"tessera/models/dtos"
	"tessera/repos"
	"tessera/util"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

type IAuthController interface {
	RegisterUser(c *fiber.Ctx) error
	VerifyEmail(c *fiber.Ctx) error
	ResendVerificationEmail(c *fiber.Ctx) error
	LoginUser(c *fiber.Ctx) error
	LogoutUser(c *fiber.Ctx) error
	ForgotPassword(c *fiber.Ctx) error
	ValidatePasswordResetToken(c *fiber.Ctx) error
	ResetPassword(c *fiber.Ctx) error
	RevokeSession(c *fiber.Ctx) error
}

type authController struct {
	userRepo repos.IUserRepository
}

func NewAuthController(userRepo repos.IUserRepository) IAuthController {
	return &authController{userRepo}
}

// TODO: let new users login without email verification but restrict access to certain routes until they verify their email. Show a "Send Verification Link" banner message on the frontend to remind them to verify their email to access all features.
func (a *authController) RegisterUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var payload dtos.RegisterDTO
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err)
	}

	if payload.IsRoleInvalid() {
		return c.Status(fiber.StatusBadRequest).JSON(util.ErrInvalidRole)
	}

	valid, err := util.ValidateStruct(payload)
	if !valid {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": govalidator.ErrorsByField(err)})
	}

	payload.Email = util.FilterEmail(payload.Email)

	fmt.Println("testtyreg: ", payload.Email)

	userExists, err := a.userRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		fmt.Println("erggy: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	if userExists.Email == payload.Email {
		return c.Status(fiber.StatusConflict).JSON(util.ErrEmailAlreadyExists)
	}

	if payload.Password != payload.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(util.ErrPasswordNotMatch)
	}

	hashedPass, err := util.HashPass(payload.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	token, err := util.RandomString(32)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	verificationToken := token

	expires := time.Now().Add(time.Second * 30)
	newUser := models.User{
		FirstName: util.FilterName(payload.FirstName),
		LastName:  util.FilterName(payload.LastName),
		Email:     payload.Email,
		Password:  hashedPass,
		PhotoURL:  payload.PhotoURL,
		Verification: &models.UserVerification{
			VerificationToken: verificationToken,
			ExpiresOn:         &expires,
		},
	}

	user, err := a.userRepo.RegisterUser(ctx, &newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	firstName := user.FirstName
	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	loadConfig, err := config.LoadConfig("./")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}

	// ? Send Email
	emailData := util.EmailData{
		URL:       loadConfig.ClientOrigin + "/api/v1/auth/verify-email?verificationToken=" + token,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	// ? Send Email in a goroutine
	emailSent := make(chan bool, 1)
	go func() {
		defer close(emailSent)
		err := util.SendVerificationEmail(&newUser, &emailData)
		if err != nil {
			log.Println("Failed to send verification email:", err)
			emailSent <- false
			// Rollback registration process
			//_ = a.userRepo.DeleteUser(ctx, user.UserID) // ? change to delete verification token for user so they resend verification email
		} else {
			emailSent <- true
		}
	}()

	select {
	case sent := <-emailSent:
		if sent {
			message := "We sent an email with a verification code to " + newUser.Email
			return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": message})
		} else {
			message := "Failed to send verification email to " + newUser.Email + ". Please try again."
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": message})
		}
	case <-time.After(10 * time.Second):
		//_ = a.userRepo.DeleteUser(ctx, user.UserID) // ? change to delete verification token for user so they resend verification email
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": util.ErrEmailNotSent})
	}
}

func (a *authController) VerifyEmail(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	verificationToken := c.Params("verificationToken")
	if verificationToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": util.ErrVerificationTokenRequired,
		})
	}

	fmt.Println("verificationToken: ", verificationToken)

	var user models.User

	user, err := a.userRepo.GetUserByVerificationToken(ctx, verificationToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"errors": "Invalid verification code: " + err.Error(),
		})
	}

	if user.Verification == nil {
		return c.Status(fiber.StatusBadRequest).JSON(util.ErrUserNotFound)
	}

	if user.Verification.ExpiresOn.Before(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"errors": util.ErrVerificationTokenExpired})
	}

	if user.Verification.IsVerified {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"errors": util.ErrEmailAlreadyVerified})
	}

	now := time.Now()
	expired := time.Now().Add(-time.Hour * 24)
	user.Verification.VerificationToken = ""
	user.Verification.ExpiresOn = &expired
	user.Verification.IsVerified = true
	user.Verification.VerifiedOn = &now

	patchUser := user.Verification

	jsonPatch, err := json.Marshal(patchUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"errors": "Failed to marshal user verification data!: " + err.Error()})
	}

	var verificationMap map[string]interface{}
	if err := json.Unmarshal(jsonPatch, &verificationMap); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"errors": "Failed to unmarshal user verification data!: " + err.Error()})
	}

	myuser, err := a.userRepo.PatchUserVerification(ctx, user.UserID, verificationMap)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully verified email!",
		"data":    myuser,
	})
}

func (a *authController) ResendVerificationEmail(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var payload dtos.ResendVerificationEmail
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err)
	}

	valid, err := util.ValidateStruct(payload)
	if !valid {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": govalidator.ErrorsByField(err)})
	}

	payload.Email = util.FilterEmail(payload.Email)

	fmt.Println("payload.Email: ", payload.Email)

	user, err := a.userRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	fmt.Println("UserVerification: ", user.Verification)

	if user.Verification == nil {
		return c.Status(fiber.StatusBadRequest).JSON(util.ErrUserNotFound)
	}

	if user.Verification.IsVerified {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
			"errors": util.ErrEmailAlreadyVerified})
	}

	// Check if the existing verification token has expired
	if user.Verification.ExpiresOn.Before(time.Now()) {
		// Generate a new verification token and update the user's record
		token, err := util.RandomString(32)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
		expires := time.Now().Add(time.Minute * 1)
		user.Verification.VerificationToken = token
		user.Verification.ExpiresOn = &expires

		patchUser := user.Verification

		jsonPatch, err := json.Marshal(patchUser)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"errors": "Failed to marshal user verification data!: " + err.Error()})
		}

		var verificationMap map[string]interface{}
		if err := json.Unmarshal(jsonPatch, &verificationMap); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "error",
				"errors": "Failed to unmarshal user verification data!: " + err.Error()})
		}

		_, err = a.userRepo.PatchUserVerification(ctx, user.UserID, verificationMap)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
	}

	verificationToken := user.Verification.VerificationToken

	firstName := user.FirstName
	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	loadConfig, err := config.LoadConfig("./")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}

	// Send Email
	emailData := util.EmailData{
		URL:       loadConfig.ClientOrigin + "/api/v1/auth/verify-email?verificationToken=" + verificationToken,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	// Send Email in a goroutine
	emailSent := make(chan bool, 1)
	go func() {
		err := util.SendVerificationEmail(&user, &emailData)
		if err != nil {
			log.Println("Failed to send verification email:", err)
			emailSent <- false
		} else {
			emailSent <- true
		}
	}()

	select {
	case sent := <-emailSent:
		if sent {
			message := "We sent an email with a verification code to " + user.Email
			return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": message})
		} else {
			message := "Failed to send verification email to " + user.Email + ". Please try again."
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": message})
		}
	case <-time.After(10 * time.Second):
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": util.ErrEmailNotSent})
	}
}

func (a *authController) LoginUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var payload dtos.LoginDTO
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err)
	}

	valid, err := util.ValidateStruct(payload)
	if !valid {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": govalidator.ErrorsByField(err)})
	}

	payload.Email = util.FilterEmail(payload.Email)

	user, err := a.userRepo.GetUser(ctx, payload.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	if user.Verification == nil || user == (models.User{}) {
		return c.Status(fiber.StatusBadRequest).JSON(util.ErrUserNotFound)
	}

	if user.Email != payload.Email {
		return c.Status(fiber.StatusBadRequest).JSON(util.ErrInvalidCredentials)
	}

	if !user.Verification.IsVerified {
		return c.Status(fiber.StatusBadRequest).JSON(util.ErrEmailNotVerified)
	}

	if err := util.ComparePass(user.Password, payload.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(util.ErrInvalidCredentials)
	}

	userID := user.UserID.String()

	store := util.GetSessionStore()

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get session in loginUser!: " + err.Error(),
		})
	}

	userInfo := UserInfo{
		UserID:    userID,
		Role:      UserType(user.Role),
		UserAgent: c.Context().UserAgent(),
		IPAddress: c.IP(),
		ExpiresOn: time.Now().Add(time.Hour * 12),
	}

	userInfoJSON, err := json.Marshal(userInfo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to marshal user in loginUser!: " + err.Error(),
		})
	}

	sess.Set("user_info", userInfoJSON)

	// c.Locals("user_id", userID)
	// c.Locals("user_agent", c.Context().UserAgent())
	// c.Locals("ip_address", c.IP())

	if err := sess.Save(); err != nil {
		sess.Destroy()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to save session in loginUser!: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		//"request_id":
		"message": "Login successful",
		// "data":    user,
	})
}

func (a *authController) LogoutUser(c *fiber.Ctx) error {

	store := util.GetSessionStore()

	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "fail",
			"errors": "Failed to get session in logoutUser!",
		})
	}

	sess.Destroy()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Logged out!"})
}

func (a *authController) ForgotPassword(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var payload dtos.ForgotPasswordDTO

	fmt.Println("Forgot password payload: ", payload)

	if err := c.BodyParser(&payload); err != nil {
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(err)
	}

	valid, err := util.ValidateStruct(payload)
	if !valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": govalidator.ErrorsByField(err)})
	}

	payload.Email = util.FilterEmail(payload.Email)

	user, err := a.userRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	if user == (models.User{}) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "fail",
			"errors": "User not found!"})
	}

	token, err := util.RandomString(20)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}
	resetToken := util.Encode(token)
	expiry := time.Now().Add(time.Minute * 30)

	err = a.userRepo.SetPasswordResetToken(ctx, user.UserID, resetToken, &expiry)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	loadConfig, err := config.LoadConfig("./")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}

	// Send email with password reset link
	firstName := user.FirstName
	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	emailData := util.EmailData{
		URL:       loadConfig.ClientOrigin + "/api/v1/auth/reset-password/" + token,
		FirstName: firstName,
		Subject:   "Password reset code",
	}

	emailSent := make(chan bool, 1)
	go func() {
		err := util.SendVerificationEmail(&user, &emailData)
		if err != nil {
			log.Println("Failed to send password reset email, Please try again:", err)
			emailSent <- false
		} else {
			emailSent <- true
		}
	}()

	select {
	case success := <-emailSent:
		if !success {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail",
				"errors": "Failed to send email, Please try again!"})
		} else {

			message := "We sent an email with password reset instructions to " + user.Email
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": message})
		}
	case <-time.After(10 * time.Second):
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "fail",
			"errors": "Timeout reached, failed to send email, Please try again!"})
	}
}

func (a *authController) ValidatePasswordResetToken(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token := c.Params("passResetToken")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"errors": "Pas is required!",
		})
	}

	resetToken := util.Encode(token)

	user, err := a.userRepo.GetUserByPasswordResetToken(ctx, resetToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"errors":  err.Error(),
			"message": "Failed to validate user for password reset",
		})
	}

	if user == (models.User{}) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "fail",
			"errors": "User not found!"})
	}

	if user.PasswordReset.PassResetToken != resetToken {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"errors": "Invalid password reset token!"})
	}

	if user.PasswordReset.ExpiresOn.Before(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"errors": "Password reset token has expired, Please try generating a new code!"})
	}

	// TODO: If not valid, redirect to reset password page with error message to resend reset password link

	c.Cookie(&fiber.Cookie{
		Name:     "passResetToken",
		Value:    resetToken,
		Path:     "/",
		MaxAge:   15 * 60,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	// TODO: Redirect to reset password page

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Password reset token is valid!",
		"data":    user,
	})
}

func (a *authController) ResetPassword(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var payload dtos.ResetPasswordDTO
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err)
	}

	if payload.Password != payload.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"errors": "Password and confirm password do not match!"})
	}

	// TODO: Get user id or password reset token from cookie
	cookie := c.Cookies("passResetToken")
	if cookie == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"errors": "Password reset token is required!",
		})
	}

	fmt.Println("Cookie: ", cookie)

	passResetToken := cookie

	user, err := a.userRepo.GetUserByPasswordResetToken(ctx, passResetToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	fmt.Println("User: ", user)

	if user == (models.User{}) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "fail",
			"errors": "User not found!"})
	}

	if user.PasswordReset.PassResetToken != passResetToken {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"errors": "Invalid password reset token!"})
	}

	if user.PasswordReset.ExpiresOn.Before(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail",
			"errors": "Password reset token has expired, Please try generating a new code!"})
	}

	hashPass, err := util.HashPass(payload.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	expired := time.Now().Add(-time.Hour * 24 * 7)
	resetOn := time.Now()
	user.Password = hashPass
	user.PasswordReset.PassResetToken = ""
	user.PasswordReset.ExpiresOn = &expired
	user.PasswordReset.PassResetOn = &resetOn

	upUser := user.PasswordReset

	jsonData, err := json.Marshal(upUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"errors": "Failed to marshal user password reset data!: " + err.Error()})
	}

	var userPasswordResetMap map[string]interface{}
	err = json.Unmarshal(jsonData, &userPasswordResetMap)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"errors": "Failed to unmarshal user password reset data!: " + err.Error()})
	}

	myuser, err := a.userRepo.PatchResetPassword(ctx, &user, userPasswordResetMap)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	// TODO: redirect to login

	// TODO: clear cookie - user id or password reset token
	c.ClearCookie("passResetToken")

	// TODO: send email to user with password reset success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Password reset successful!",
		"data":    myuser,
	})
}

// TODO: Use session ID to revoke session after using user ID to get sessions from Redis and you will have to delete sessions for that user
// ? Revoke session route will be POST /api/v1/auth/revoke-session and will be called from the frontend with payload of session ID to revoke
// ? In the frontend, a user will have a session manager for retrieved sessions from Redis and that we display in a list with a revoke button for each session
// ? When revoke button clicked it will send the session id as body to the revoke-session route
// ? In revoke-session controller i will retrieve the session id and pass it to Delete Session function to revoke it and logout the client using the session id

func (a *authController) RevokeSession(c *fiber.Ctx) error {
	sessionID := c.Params("sessionID")
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "session id required",
		})
	}

	util.DelSession(c, sessionID)

	return c.Status(fiber.StatusOK).Send([]byte("session revoked"))
}
