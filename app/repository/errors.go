package repository

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// Sentinel error milik repository. Service menerjemahkannya menjadi status HTTP.
var (
	ErrNotFound  = errors.New("data tidak ditemukan")
	ErrDuplicate = errors.New("data sudah ada")
	ErrInUse     = errors.New("data masih dipakai data lain")
)

const (
	pgUniqueViolation     = "23505"
	pgForeignKeyViolation = "23503"
)

// pgErrorCode mengambil kode error PostgreSQL, atau string kosong bila
// error tersebut bukan error dari PostgreSQL.
func pgErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

// mapWriteError menerjemahkan pelanggaran UNIQUE menjadi ErrDuplicate.
func mapWriteError(err error) error {
	if pgErrorCode(err) == pgUniqueViolation {
		return ErrDuplicate
	}
	return err
}
