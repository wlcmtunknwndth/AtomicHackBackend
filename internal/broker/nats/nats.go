package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/wlcmtunknwndth/AtomicHackBackend/internal/config"
	"github.com/wlcmtunknwndth/AtomicHackBackend/internal/storage"
	"time"
)

type Storage interface {
	SaveRequest(*storage.Request) error
	GetRequest(string) (*storage.Request, error)
	SaveResponse(response *storage.Response) error
	GetResponse(string) (*storage.Response, error)
}

type Broker struct {
	conn *nats.Conn
	db   Storage
}

func New(cfg *config.Broker, db Storage) (*Broker, error) {
	const op = "internal.broker.nats.New"
	natsSrv, err := nats.Connect(cfg.Address,
		nats.RetryOnFailedConnect(cfg.Retry),
		nats.MaxReconnects(cfg.MaxReconnects),
		nats.ReconnectWait(cfg.ReconnectWait),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = natsSrv.Flush(); err != nil {
		return nil, fmt.Errorf("%s: flush: %w", op, err)
	}
	if err = natsSrv.FlushTimeout(time.Second); err != nil {
		return nil, fmt.Errorf("%s: flush timeout: %w", op, err)
	}
	return &Broker{
		conn: natsSrv,
		db:   db,
	}, nil
}

func (b *Broker) Close() {
	b.conn.Close()
}
