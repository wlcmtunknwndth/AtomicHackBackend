package chat

import (
	"github.com/gorilla/websocket"
	"github.com/wlcmtunknwndth/AtomicHackBackend/internal/storage"
)

type Storage interface {
	AskRequest(id string) ([]byte, error)
	AskSaveRequest(request *storage.Request) error
	AskSaveResponse(response *storage.Response) error
	AskResponse(id string) ([]byte, error)
}

type Handler struct {
	ws              websocket.Upgrader
	storage         Storage
	receiverAddress string
}

func New(ReadBufSize, WriteBufSize int, storage Storage, addr string) *Handler {
	return &Handler{
		ws: websocket.Upgrader{
			ReadBufferSize:  ReadBufSize,
			WriteBufferSize: WriteBufSize,
		},
		storage:         storage,
		receiverAddress: addr,
	}
}
