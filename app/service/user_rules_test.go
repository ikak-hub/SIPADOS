package service

import (
	"strings"
	"testing"

	"sipasdos/app/model"
)

func TestNormalizeEmail(t *testing.T) {
	if got := NormalizeEmail("  Ika@Unair.AC.ID "); got != "ika@unair.ac.id" {
		t.Errorf("harap ika@unair.ac.id, dapat %q", got)
	}
}

func TestValidateRegister(t *testing.T) {
	t.Run("permintaan valid", func(t *testing.T) {
		user, errs := ValidateRegister(model.RegisterRequest{
			Name: "  Ika ", Email: "Ika@Mail.com", Password: "rahasia123",
		})
		if len(errs) != 0 {
			t.Fatalf("tidak seharusnya ada error: %v", errs)
		}
		if user.Role != model.RoleStudent {
			t.Errorf("pendaftaran publik harus berrole mahasiswa, dapat %q", user.Role)
		}
		if user.Name != "Ika" || user.Email != "ika@mail.com" {
			t.Errorf("nama dan email harus dirapikan: %+v", user)
		}
	})

	cases := []struct {
		name      string
		req       model.RegisterRequest
		wantField string
	}{
		{"nama kosong", model.RegisterRequest{Name: " ", Email: "a@b.co", Password: "rahasia123"}, "name"},
		{"email tanpa @", model.RegisterRequest{Name: "A", Email: "ab.co", Password: "rahasia123"}, "email"},
		{"email berspasi", model.RegisterRequest{Name: "A", Email: "a b@c.co", Password: "rahasia123"}, "email"},
		{"password terlalu pendek", model.RegisterRequest{Name: "A", Email: "a@b.co", Password: "abc123"}, "password"},
		{"password tanpa angka", model.RegisterRequest{Name: "A", Email: "a@b.co", Password: "hurufsaja"}, "password"},
		{"password tanpa huruf", model.RegisterRequest{Name: "A", Email: "a@b.co", Password: "12345678"}, "password"},
		{"password lebih dari 72 byte", model.RegisterRequest{Name: "A", Email: "a@b.co", Password: strings.Repeat("a1", 40)}, "password"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, errs := ValidateRegister(tc.req)
			if _, ok := errs[tc.wantField]; !ok {
				t.Errorf("harap ada error pada field %q, dapat %v", tc.wantField, errs)
			}
		})
	}
}

func TestValidateLogin(t *testing.T) {
	if errs := ValidateLogin(model.LoginRequest{Email: "a@b.co", Password: "x"}); len(errs) != 0 {
		t.Errorf("tidak seharusnya ada error: %v", errs)
	}
	errs := ValidateLogin(model.LoginRequest{})
	if errs["email"] == "" || errs["password"] == "" {
		t.Errorf("email dan password kosong harus ditolak: %v", errs)
	}
}
