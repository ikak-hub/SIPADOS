package helper

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"sipasdos/app/model"
)

// Role dipakai oleh route dan middleware. Nilainya berasal dari model
// supaya satu-satunya sumber kebenaran tetap di layer entitas.
const (
	RoleAdmin    = model.RoleAdmin
	RoleLecturer = model.RoleLecturer
	RoleStudent  = model.RoleStudent
)

// Kunci c.Locals tempat middleware JWT menyimpan identitas pengguna.
const (
	LocalUserID = "auth_user_id"
	LocalRole   = "auth_role"
)

const tokenIssuer = "sipasdos"

var ErrInvalidToken = errors.New("token tidak valid atau kedaluwarsa")

// Claims adalah isi token JWT.
type Claims struct {
	UserID int    `json:"uid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// HashPassword menghasilkan hash bcrypt. Password asli tidak pernah disimpan.
func HashPassword(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword membandingkan password dengan hash bcrypt.
func CheckPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

// GenerateToken membuat token JWT (HS256) beserta waktu kedaluwarsanya.
func GenerateToken(secret string, ttl time.Duration, userID int, role string) (string, time.Time, error) {
	now := time.Now()
	expires := now.Add(ttl)

	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    tokenIssuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expires, nil
}

// ParseToken memeriksa tanda tangan, algoritma, penerbit, dan masa berlaku token.
// Hanya HS256 yang diterima, sehingga token dengan algoritma lain ditolak.
func ParseToken(secret, tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (any, error) { return []byte(secret), nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(tokenIssuer),
	)
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

// CurrentUser membaca identitas yang dipasang middleware JWT.
// ok bernilai false bila route tidak dilindungi middleware tersebut.
func CurrentUser(c *fiber.Ctx) (int, string, bool) {
	id, okID := c.Locals(LocalUserID).(int)
	role, okRole := c.Locals(LocalRole).(string)
	return id, role, okID && okRole
}
