package sensitive

import "log/slog"

func badMessages() {
	slog.Info("user password reset")  // want `log message may contain sensitive data`
	slog.Error("token expired")       // want `log message may contain sensitive data`
	slog.Info("invalid jwt provided") // want `log message may contain sensitive data`
	slog.Info("bearer token missing") // want `log message may contain sensitive data`
	slog.Warn("credential rejected")  // want `log message may contain sensitive data`
	slog.Info("wrong pwd supplied")   // want `log message may contain sensitive data`
}

func goodMessages() {
	slog.Info("server started")
	slog.Error("connection failed")
	slog.Info("user not found")
	slog.Info("request processed successfully")
	slog.Info("")
}
