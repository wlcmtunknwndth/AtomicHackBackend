package chat

import (
	"github.com/gorilla/websocket"
	"github.com/wlcmtunknwndth/AtomicHackBackend/internal/storage"
)

type Storage interface {
	AskRequest(id uint64) ([]byte, error)
	AskSaveRequest(request *storage.Request) error
	AskSaveResponse(response *storage.Response) error
	AskResponse(id uint64) ([]byte, error)
}

type Handler struct {
	ws      websocket.Upgrader
	storage Storage
}

func New(ReadBufSize, WriteBufSize int, storage Storage) *Handler {
	return &Handler{
		ws: websocket.Upgrader{
			ReadBufferSize:  ReadBufSize,
			WriteBufferSize: WriteBufSize,
		},
		storage: storage,
	}
}
