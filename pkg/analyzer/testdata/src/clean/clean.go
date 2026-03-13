package clean

import (
	"log"
	"log/slog"
)

func allGood() {
	slog.Debug("cache miss for key")
	slog.Info("server started on port 8080")
	slog.Warn("retry attempt 3 of 5")
	slog.Error("connection refused by remote host")

	log.Print("application initialized")
	log.Printf("listening on %s", "8080")

	slog.Info("user", "user_id", 42, "method", "oauth2")
	slog.Error("query failed", "table", "users", "duration_ms", 150)

	slog.Info("123 events processed")
	slog.Info("")
}
