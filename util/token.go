package util

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	ErrTokenUnauthorized = errors.New("token is unauthorized")
	ErrExpired           = errors.New("token is expired")
	ErrNBFInvalid        = errors.New("token nbf validation failed")
	ErrIATInvalid        = errors.New("token iat validation failed")
	ErrNoTokenFound      = errors.New("no token found")
	ErrAlgoInvalid       = errors.New("algorithm mismatch")
)

// func ErrorReason(err error) error {
// 	switch {
// 	case errors.Is(err, jwt.ErrTokenExpired()), err == ErrExpired:
// 		return ErrExpired
// 	case errors.Is(err, jwt.ErrInvalidIssuedAt()), err == ErrIATInvalid:
// 		return ErrIATInvalid
// 	case errors.Is(err, jwt.ErrTokenNotYetValid()), err == ErrNBFInvalid:
// 		return ErrNBFInvalid
// 	default:
// 		return ErrUnauthorized
// 	}
// }

type TokenInfo struct {
	Token     *string
	TokenUUID string
	UserID    string
	ExpiresIn *int64
}

// func GenerateTokenAndCookie(c echo.Context, userID *uuid.UUID, loadConfig config.Config) (string, error) {
// 	now := time.Now().UTC()
// 	expiresIn := now.Add(loadConfig.JwtExpiresIn).Unix()
// 	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"sub": userID.String(),
// 		"exp": expiresIn,
// 		"iat": now.Unix(),
// 		"nbf": now.Unix(),
// 	})
// 	tokenString, err := jwtToken.SignedString([]byte(loadConfig.JwtSecret))
// 	if err != nil {
// 		return "", c.JSON(http.StatusInternalServerError, "failed while signing jwt token: "+err.Error())
// 	}
// 	jwtCookie(c, tokenString, loadConfig)
// 	return tokenString, nil
// }

func GenerateToken(userID uuid.UUID, ttl time.Duration, privateKey string) (*TokenInfo, error) {
	now := time.Now().UTC()
	ti := &TokenInfo{
		ExpiresIn: new(int64),
		Token:     new(string),
	}

	*ti.ExpiresIn = now.Add(ttl).Unix()
	ti.TokenUUID = uuid.New().String()
	ti.UserID = userID.String()

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key token: %w", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key token: %w", err)
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":        userID.String(),
		"token_uuid": ti.TokenUUID,
		"exp":        ti.ExpiresIn,
		"iat":        now.Unix(),
		"nbf":        now.Unix(),
	})

	signedJwtToken, err := jwtToken.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("failed while signing token: %w", err)
	}

	*ti.Token = signedJwtToken

	return ti, nil
}

func ParseAndValidateToken(token string, publicKey string) (*TokenInfo, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, err := validateToken(parsedToken)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %s", err.Error())
	}

	return &TokenInfo{
		TokenUUID: fmt.Sprint(claims["token_uuid"]),
		UserID:    fmt.Sprint(claims["sub"]),
	}, nil
}

func validateToken(parsedtoken *jwt.Token) (jwt.MapClaims, error) {
	if !parsedtoken.Valid {
		return nil, fmt.Errorf("parsed token is invalid: %s", ErrUnauthorized)
	}

	claims, ok := parsedtoken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("bad jwt claims: %s", ErrUnauthorized)
	}

	return claims, nil
}

// func RefreshToken(c echo.Context, loadConfig config.Config) (string, error) {
// 	token, err := ParseAndValidateToken(c, loadConfig)
// 	if err != nil {
// 		return "", c.JSON(http.StatusUnauthorized, "failed to parse token: "+err.Error())
// 	}

// 	claims := token.Claims.(jwt.MapClaims)

// 	var userID uuid.UUID
// 	userID, err = uuid.Parse(claims["sub"].(string))
// 	if err != nil {
// 		return "", fmt.Errorf("failed to parse uuid: %s", err.Error())
// 	}

// 	return GenerateTokenAndCookie(c, &userID, loadConfig)
// }

func DeleteJwtCookie(c echo.Context, tokenName string, httpOnly bool) {
	cookie := new(http.Cookie)
	cookie.Name = tokenName
	cookie.Value = ""
	cookie.Path = "/"
	cookie.MaxAge = -1
	cookie.Expires = time.Now().Add(-24 * time.Hour)
	cookie.Secure = false
	cookie.HttpOnly = httpOnly
	cookie.Domain = "localhost"
	cookie.SameSite = http.SameSiteStrictMode

	c.SetCookie(cookie)
}

// func jwtCookie(c echo.Context, tokenString string, loadConfig config.Config) {
// 	cookie := new(http.Cookie)
// 	cookie.Name = "jwtToken"
// 	cookie.Value = tokenString
// 	cookie.Path = "/" // cookie is valid for the entire site
// 	cookie.MaxAge = loadConfig.JwtMaxAge * 60
// 	//cookie.Expires = time.Now().Add(time.Minute * 60)
// 	cookie.Secure = false  // secure false for development(http) and true for production(https)
// 	cookie.HttpOnly = true // cookie is not accessible from javascript
// 	cookie.Domain = loadConfig.JwtDomain
// 	cookie.SameSite = http.SameSiteStrictMode

// 	c.SetCookie(cookie)
// }

func MakeCookie(c echo.Context, tokenName string, token *string, maxAge int, httpOnly bool) {
	cookie := new(http.Cookie)
	cookie.Name = tokenName
	cookie.Value = *token
	cookie.Path = "/" // cookie is valid for the entire site
	cookie.MaxAge = maxAge * 60
	cookie.Secure = false      // secure false for development(http) and true for production(https)
	cookie.HttpOnly = httpOnly // cookie is not accessible from javascript
	cookie.Domain = "localhost"
	cookie.SameSite = http.SameSiteStrictMode

	c.SetCookie(cookie)
}
