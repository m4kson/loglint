package lowercase

import "log/slog"

func badMessages() {
	slog.Info("Server started")      // want `log message must start with a lowercase letter`
	slog.Error("Failed to connect")  // want `log message must start with a lowercase letter`
	slog.Warn("Retry attempt")       // want `log message must start with a lowercase letter`
	slog.Debug("Processing request") // want `log message must start with a lowercase letter`
	slog.Info("ERROR disk full")     // want `log message must start with a lowercase letter`
}

func goodMessages() {
	slog.Info("server started")
	slog.Error("failed to connect")
	slog.Warn("retry attempt")
	slog.Debug("processing request")
	slog.Info("123 items processed")
	slog.Info("")
}
