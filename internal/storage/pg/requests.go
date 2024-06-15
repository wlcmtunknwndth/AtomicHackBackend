package pg

import (
	"fmt"
	"github.com/wlcmtunknwndth/AtomicHackBackend/internal/storage"
)

func (s *Storage) SaveRequest(request *storage.Request) error {
	const op = "internal.storage.postgres.SaveRequest"
	if err := s.db.Create(request).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) GetRequest(id uint64) (*storage.Request, error) {
	const op = "internal.storage.postgres.GetRequest"
	var request storage.Request
	if err := s.db.First(&request, id).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &request, nil
}
