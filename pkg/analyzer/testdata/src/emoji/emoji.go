package emoji

import "log/slog"

func emojiMessages() {
	slog.Info("😀 well done")      // want `log message must contain only English` `log message must not contain emoji`
	slog.Error("check passed 👻")  // want `log message must contain only English` `log message must not contain emoji`
	slog.Warn("warning 😰 issued") // want `log message must contain only English` `log message must not contain emoji`
	slog.Info("done cri🤔tical")   // want `log message must contain only English` `log message must not contain emoji`
}

func goodMessages() {
	slog.Info("server started")
	slog.Error("connection failed")
	slog.Info("")
}
