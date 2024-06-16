package pg

import (
	"fmt"
	"github.com/wlcmtunknwndth/AtomicHackBackend/internal/storage"
)

func (s *Storage) SaveResponse(response *storage.Response) error {
	const op = "internal.storage.pg.SaveRequest"
	if err := s.db.Create(response).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	//s.db.Create(response)
	return nil
}

func (s *Storage) GetResponse(id string) (*storage.Response, error) {
	const op = "internal.storage.pg.SaveRequest"
	var response storage.Response
	if err := s.db.First(&response, id).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &response, nil
}
