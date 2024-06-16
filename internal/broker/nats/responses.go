package nats

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/wlcmtunknwndth/AtomicHackBackend/internal/storage"
	"github.com/wlcmtunknwndth/AtomicHackBackend/lib/slogResp"
	"log/slog"
	"time"
)

const (
	MustSaveResponse = "save_resp"
	AskFindResponse  = "resp:"
	MustFindResponse = AskFindResponse + "*"
)

func (b *Broker) ResponseFinder() (*nats.Subscription, error) {
	const op = "internal.broker.nats.ResponseFinder"
	sub, err := b.conn.Subscribe(MustFindRequest, func(msg *nats.Msg) {
		id := msg.Subject[5:]
		if len([]byte(id)) == 16 {
			slog.Error("couldn't convert wildcard var to uint64", slogResp.Info(op, id))
			return
		}

		response, err := b.db.GetResponse(id)
		if err != nil {
			slog.Error("couldn't get request", slogResp.Error(op, err))
			return
		}

		data, err := json.Marshal(response)
		if err != nil {
			slog.Error("couldn't marshall request", slogResp.Error(op, err))
			return
		}

		if err = msg.Respond(data); err != nil {
			slog.Error("couldn't send reply", slogResp.Error(op, err))
		}
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return sub, nil
}

func (b *Broker) ResponseSaver() (*nats.Subscription, error) {
	const op = "broker.nats.event.ResponseSaver"
	sub, err := b.conn.Subscribe(MustSaveRequest, func(msg *nats.Msg) {
		var response storage.Response
		if err := json.Unmarshal(msg.Data, &response); err != nil {
			slog.Error("couldn't unmarshal resp", slogResp.Error(op, err))
			return
		}

		if err := b.db.SaveResponse(&response); err != nil {
			slog.Error("couldn't save resp", slogResp.Error(op, err))
			return
		}
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return sub, nil
}

func (b *Broker) AskSaveResponse(response *storage.Response) error {
	const op = "internal.broker.nats.AskSaveRequest"
	data, err := json.Marshal(*response)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = b.conn.Publish(MustSaveResponse, data); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (b *Broker) AskResponse(id string) ([]byte, error) {
	const op = "internal.broker.nats.AskRequest"
	msg, err := b.conn.Request(AskFindResponse+id, nil, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return msg.Data, nil
}
