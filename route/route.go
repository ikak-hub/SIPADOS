package route

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"sipasdos/app/service"
	"sipasdos/helper"
	"sipasdos/middleware"
)

// Services mengumpulkan seluruh service aplikasi. Setiap anggota kelompok
// menambahkan satu field untuk modul miliknya, lalu mendaftarkan route-nya
// pada Register di bawah.
type Services struct {
	Auth   *service.AuthService
	Period *service.PeriodService
	// Vacancy     *service.VacancyService
	// Application *service.ApplicationService
}

// Register memetakan URL ke method pada service.
//
// Perhatikan isi file ini: tidak ada logika bisnis, tidak ada query,
// tidak ada validasi. Hanya daftar alamat dan siapa yang melayaninya.
//
// authRequired adalah middleware.JWT yang sudah berisi secret; route tidak
// perlu mengetahui secret itu sendiri.
func Register(app *fiber.App, pool *pgxpool.Pool, s Services, authRequired fiber.Handler) {
	api := app.Group("/api/v1")
	adminOnly := middleware.RequireRole(helper.RoleAdmin)

	// Publik
	api.Get("/health", healthCheck(pool))

	auth := api.Group("/auth", middleware.RequireJSON)
	auth.Post("/register", s.Auth.Register)
	auth.Post("/login", s.Auth.Login)

	// Perlu token (semua role)
	api.Get("/auth/me", authRequired, s.Auth.Me)

	// Periode: baca untuk semua role, tulis hanya admin
	periods := api.Group("/periods", authRequired, middleware.RequireJSON)
	periods.Get("/", s.Period.List)
	periods.Get("/:id", s.Period.Get)
	periods.Post("/", adminOnly, s.Period.Create)
	periods.Put("/:id", adminOnly, s.Period.Replace)
	periods.Patch("/:id", adminOnly, s.Period.Patch)
	periods.Delete("/:id", adminOnly, s.Period.Delete)
}

// healthCheck melaporkan kondisi layanan beserta databasenya.
func healthCheck(pool *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			return helper.Fail(c, fiber.StatusServiceUnavailable,
				"database tidak dapat dihubungi")
		}
		return helper.Success(c, fiber.StatusOK, "server dan database berjalan", nil)
	}
}
