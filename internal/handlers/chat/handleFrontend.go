package chat

import (
	"encoding/json"
	"github.com/wlcmtunknwndth/AtomicHackBackend/internal/storage"
	"github.com/wlcmtunknwndth/AtomicHackBackend/lib/slogResp"
	"log/slog"
	"net/http"
)

func (h *Handler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	const op = "internal.handlers.chat.handleFrontend.HandlerConnections"
	ws, err := h.ws.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("couldn't establish websocket", slogResp.Error(op, err))
		return
	}
	defer ws.Close()

	for {
		var msg []byte
		_, msg, err = ws.ReadMessage()
		if err != nil {
			slog.Error("couldn't read msg", slogResp.Error(op, err))
			break
		}
		slog.Info("received message", slog.Any("message", msg))

		var request storage.Request
		if err = json.Unmarshal(msg, &request); err != nil {
			slog.Error("couldn't unmarshal request", slogResp.Error(op, err))
			break
		}

		if err = h.storage.AskSaveRequest(&request); err != nil {
			slog.Error("couldn't unmarshal request", slogResp.Error(op, err))
			break
		}

	}
}
