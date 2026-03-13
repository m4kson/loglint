package specialchars

import "log/slog"

func controlCharMessages() {
	slog.Info("connection failed!!!")             // want `log message must not contain special characters`
	slog.Warn("warning: something went wrong...") // want `log message must not contain special characters`
	slog.Error("this is #3")                      // want `log message must not contain special characters`
}

func goodMessages() {
	slog.Info("server started")
	slog.Error("connection failed")
	slog.Info("user is 42 status is ok")
	slog.Info("")
}
