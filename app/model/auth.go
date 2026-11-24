package model

import "time"

// RegisterRequest dipakai untuk pendaftaran akun. Role tidak dikirim klien:
// pendaftaran publik selalu menghasilkan role mahasiswa.
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse adalah isi data pada respons login.
type AuthResponse struct {
	Token     string    `json:"token"`
	TokenType string    `json:"token_type"`
	ExpiresAt time.Time `json:"expires_at"`
	User      User      `json:"user"`
}
