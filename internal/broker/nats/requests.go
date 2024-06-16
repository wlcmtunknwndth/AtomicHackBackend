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
	MustSaveRequest = "save_req"
	AskFindRequest  = "req:"
	MustFindRequest = AskFindRequest + "*"
)

//func convertStrToUint(str string) (uint64, error) {
//	return strconv.ParseUint(str, 10, 64)
//}
//
//func convertUintToStr(id uint64) string {
//	return strconv.FormatUint(id, 10)
//}

func (b *Broker) RequestFinder() (*nats.Subscription, error) {
	const op = "internal.broker.nats.RequestFinder"
	sub, err := b.conn.Subscribe(MustFindRequest, func(msg *nats.Msg) {
		id := msg.Subject[4:]
		if len([]byte(id)) != 16 {
			slog.Error("couldn't get uuid from wildcard", slogResp.Info(op, id))
			return
		}

		request, err := b.db.GetRequest(id)
		if err != nil {
			slog.Error("couldn't get request", slogResp.Error(op, err))
			return
		}

		data, err := json.Marshal(request)
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

func (b *Broker) RequestSaver() (*nats.Subscription, error) {
	const op = "broker.nats.event.RequestSaver"
	sub, err := b.conn.Subscribe(MustSaveRequest, func(msg *nats.Msg) {
		var request storage.Request
		if err := json.Unmarshal(msg.Data, &request); err != nil {
			slog.Error("couldn't unmarshal request", slogResp.Error(op, err))
			return
		}

		if err := b.db.SaveRequest(&request); err != nil {
			slog.Error("couldn't save request", slogResp.Error(op, err))
			return
		}
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return sub, nil
}

func (b *Broker) AskSaveRequest(request *storage.Request) error {
	const op = "internal.broker.nats.AskSaveRequest"
	data, err := json.Marshal(*request)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = b.conn.Publish(MustSaveRequest, data); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (b *Broker) AskRequest(id string) ([]byte, error) {
	const op = "internal.broker.nats.AskRequest"
	msg, err := b.conn.Request(AskFindRequest+id, nil, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return msg.Data, nil
}
