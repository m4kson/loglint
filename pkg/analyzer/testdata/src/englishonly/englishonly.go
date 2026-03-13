package englishonly

import "log/slog"

func badMessages() {
	slog.Info("сервер запущен")    // want `log message must contain only English`
	slog.Error("соединение упало") // want `log message must contain only English`
	slog.Warn("повтор попытки")    // want `log message must contain only English`
	slog.Info("café au lait")      // want `log message must contain only English`
	slog.Info("naïve approach")    // want `log message must contain only English`
}

func goodMessages() {
	slog.Info("server started")
	slog.Error("connection failed")
	slog.Warn("retry attempt")
	slog.Info("user id is 42 status is active")
	slog.Info("processing complete")
	slog.Info("")
}
