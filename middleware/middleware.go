package middleware

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"sipasdos/helper"
)

// Register memasang seluruh middleware yang berlaku untuk semua route.
// URUTAN PENTING: middleware dieksekusi sesuai urutan pemasangan.
func Register(app *fiber.App, logger *slog.Logger) {
	app.Use(requestid.New())       // 1. beri setiap request satu ID unik
	app.Use(recover.New())         // 2. tangkap panic agar server tidak mati
	app.Use(helmet.New())          // 3. pasang header keamanan dasar
	app.Use(cors.New())            // 4. atur Cross-Origin Resource Sharing
	app.Use(RequestLogger(logger)) // 5. catat setiap request
}

// RequestLogger mencatat setiap request ke log terstruktur (JSON).
func RequestLogger(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next() // serahkan ke middleware/handler berikutnya

		requestID, _ := c.Locals("requestid").(string)
		logger.Info("http_request",
			slog.String("request_id", requestID),
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", c.Response().StatusCode()),
			slog.Duration("duration", time.Since(start)),
			slog.String("ip", c.IP()),
		)
		return err
	}
}

var methodsWithBody = map[string]bool{
	fiber.MethodPost:  true,
	fiber.MethodPut:   true,
	fiber.MethodPatch: true,
}

// RequireJSON menolak request berisi body yang Content-Type-nya bukan JSON.
// Dipasang per grup route, bukan global.
func RequireJSON(c *fiber.Ctx) error {
	if methodsWithBody[c.Method()] {
		ct := c.Get("Content-Type")
		if !strings.HasPrefix(ct, fiber.MIMEApplicationJSON) {
			return helper.Fail(c, fiber.StatusUnsupportedMediaType,
				"Content-Type harus application/json")
		}
	}
	return c.Next()
}

// JWT memeriksa header "Authorization: Bearer <token>". Bila valid, identitas
// pengguna disimpan di c.Locals (dibaca lewat helper.CurrentUser); bila tidak,
// permintaan dihentikan dengan 401. Dipasang per grup route, bukan global.
func JWT(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		scheme, token, found := strings.Cut(c.Get(fiber.HeaderAuthorization), " ")
		token = strings.TrimSpace(token)
		if !found || !strings.EqualFold(scheme, "Bearer") || token == "" {
			c.Set(fiber.HeaderWWWAuthenticate, "Bearer")
			return helper.Fail(c, fiber.StatusUnauthorized, "token autentikasi diperlukan")
		}

		claims, err := helper.ParseToken(secret, token)
		if err != nil {
			c.Set(fiber.HeaderWWWAuthenticate, "Bearer")
			return helper.Fail(c, fiber.StatusUnauthorized, "token tidak valid atau kedaluwarsa")
		}

		c.Locals(helper.LocalUserID, claims.UserID)
		c.Locals(helper.LocalRole, claims.Role)
		return c.Next()
	}
}

// RequireRole hanya meloloskan pengguna dengan salah satu role yang diizinkan.
// Harus dipasang SETELAH middleware JWT; bila tidak, hasilnya 401.
func RequireRole(roles ...string) fiber.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, role := range roles {
		allowed[role] = true
	}

	return func(c *fiber.Ctx) error {
		_, role, ok := helper.CurrentUser(c)
		if !ok {
			return helper.Fail(c, fiber.StatusUnauthorized, "token autentikasi diperlukan")
		}
		if !allowed[role] {
			return helper.Fail(c, fiber.StatusForbidden, "akses ditolak untuk role ini")
		}
		return c.Next()
	}
}
