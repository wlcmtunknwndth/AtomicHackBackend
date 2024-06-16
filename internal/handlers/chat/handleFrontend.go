package chat

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/wlcmtunknwndth/AtomicHackBackend/internal/storage"
	"github.com/wlcmtunknwndth/AtomicHackBackend/lib/slogResp"
	"log/slog"
	"net/http"
	"time"
)

func (h *Handler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	const op = "internal.handlers.chat.handleFrontend.HandlerConnections"
	ws, err := h.ws.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("couldn't establish websocket", slogResp.Error(op, err))
		return
	}
	defer func(ws *websocket.Conn) {
		if err = ws.Close(); err != nil {
			slog.Error("couldn't close conn", slogResp.Error(op, err))
		}
		return
	}(ws)

	var received = make(chan *storage.Response)
	var send = make(chan *storage.Request)
	go h.receiveHandler(send, received)

	for {
		select {
		case <-time.After(time.Second):
			var msg []byte
			_, msg, err = ws.ReadMessage()
			if err != nil {
				slog.Error("couldn't read msg", slogResp.Error(op, err))
				break
			}
			var request storage.Request
			if err = json.Unmarshal(msg, &request); err != nil {
				slog.Error("couldn't unmarshal request", slogResp.Error(op, err))
				continue
			}
			if request.ID, err = uuid.NewUUID(); err != nil {
				slog.Error("couldn't generate uuid for request", slogResp.Error(op, err))
				return
			}

			if err = h.storage.AskSaveRequest(&request); err != nil {
				slog.Error("couldn't unmarshal request", slogResp.Error(op, err))
				continue
			}

			send <- &request

			if err = h.storage.AskSaveResponse(<-received); err != nil {
				slog.Error("couldn't save response", slogResp.Error(op, err))
				continue
			}
		}
	}
}

func (h *Handler) receiveHandler(send chan *storage.Request, received chan *storage.Response) {
	const op = "internal.handlers.chat.receiveHandler"

	conn, _, err := websocket.DefaultDialer.Dial(h.receiverAddress, nil)
	if err != nil {
		slog.Error("couldn't connect to request receiver", slogResp.Error(op, err))
		return
	}

	for {
		select {
		case <-time.After(time.Second):
			_, msg, err := conn.ReadMessage()
			if err != nil {
				slog.Error("couldn't read response", slogResp.Error(op, err))
				break
			}
			var resp storage.Response
			if err = json.Unmarshal(msg, &resp); err != nil {
				slog.Error("couldn't unmarshal response", slogResp.Error(op, err))
				continue
			}

			if err = h.storage.AskSaveResponse(&resp); err != nil {
				slog.Error("couldn't save response", slogResp.Error(op, err))
			}
			received <- &resp
		}
	}
}
