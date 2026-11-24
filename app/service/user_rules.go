package service

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"sipasdos/app/model"
)

// File ini berisi business rules MURNI untuk akun dan autentikasi:
// tidak menyentuh fiber.Ctx, database, maupun JWT.

const (
	maxNameLength     = 100
	maxEmailLength    = 150
	minPasswordLength = 8
	maxPasswordBytes  = 72 // batas masukan bcrypt
)

// NormalizeEmail menyeragamkan email agar pencarian dan UNIQUE konsisten.
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// isValidEmail adalah pemeriksaan sederhana, bukan validasi RFC.
func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	at := strings.Index(email, "@")
	dot := strings.LastIndex(email, ".")
	return at > 0 && dot > at+1 && dot < len(email)-1 && !strings.ContainsAny(email, " \t\n")
}

// checkPassword mengembalikan pesan kesalahan, atau string kosong bila lolos.
func checkPassword(password string) string {
	if utf8.RuneCountInString(password) < minPasswordLength {
		return "minimal 8 karakter"
	}
	if len(password) > maxPasswordBytes {
		return "maksimal 72 byte"
	}
	hasLetter, hasDigit := false, false
	for _, r := range password {
		if unicode.IsLetter(r) {
			hasLetter = true
		}
		if unicode.IsDigit(r) {
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return "harus mengandung huruf dan angka"
	}
	return ""
}

// ValidateRegister memeriksa permintaan pendaftaran dan membentuk entitas
// akun (tanpa password_hash; hash dibuat oleh service). Pendaftaran publik
// selalu berrole mahasiswa, sehingga klien tidak bisa menjadikan dirinya admin.
func ValidateRegister(req model.RegisterRequest) (model.User, map[string]string) {
	errs := map[string]string{}

	name := strings.TrimSpace(req.Name)
	switch {
	case name == "":
		errs["name"] = "wajib diisi"
	case utf8.RuneCountInString(name) > maxNameLength:
		errs["name"] = "maksimal 100 karakter"
	}

	email := NormalizeEmail(req.Email)
	switch {
	case !isValidEmail(email):
		errs["email"] = "format email tidak valid"
	case len(email) > maxEmailLength:
		errs["email"] = "maksimal 150 karakter"
	}

	if msg := checkPassword(req.Password); msg != "" {
		errs["password"] = msg
	}

	if len(errs) > 0 {
		return model.User{}, errs
	}
	return model.User{Name: name, Email: email, Role: model.RoleStudent}, errs
}

// ValidateLogin hanya memeriksa bahwa field terisi. Aturan password tidak
// diperiksa di sini agar tidak membocorkan kebijakan lewat pesan login.
func ValidateLogin(req model.LoginRequest) map[string]string {
	errs := map[string]string{}
	if strings.TrimSpace(req.Email) == "" {
		errs["email"] = "wajib diisi"
	}
	if req.Password == "" {
		errs["password"] = "wajib diisi"
	}
	return errs
}
