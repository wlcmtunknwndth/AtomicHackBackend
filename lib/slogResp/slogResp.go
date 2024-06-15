package slogResp

import "log/slog"

func Error(op string, err error) slog.Attr {
	return slog.Any(op, err)
}

func Info(msg string, data any) slog.Attr {
	return slog.Any(msg, data)
}
