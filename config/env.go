package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv memuat file .env bila ada. Bila tidak ada, aplikasi memakai
// environment sistem (misalnya saat dijalankan di dalam Docker).
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("file .env tidak ditemukan, memakai environment sistem")
	}
}

// GetEnv membaca variabel environment dan memberi nilai bawaan bila kosong.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
