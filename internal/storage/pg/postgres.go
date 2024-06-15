package pg

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/wlcmtunknwndth/AtomicHackBackend/internal/config"
	"github.com/wlcmtunknwndth/AtomicHackBackend/internal/storage"
	"github.com/wlcmtunknwndth/AtomicHackBackend/lib/slogResp"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
)

type Storage struct {
	db *gorm.DB
}

func New(config *config.Database) (*Storage, error) {
	const op = "internal.storage.postges.New"
	connStr := fmt.Sprintf("postgres://%s:%s@postgres:%s/%s?sslmode=%s",
		config.DbUser, config.DbPass, config.Port,
		config.DbName, config.SslMode,
	)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		slog.Error("couldn't connect to database", slogResp.Error(op, err))
		return nil, err
	}

	if err = db.AutoMigrate(&storage.User{},
		&storage.Request{}, &storage.Response{},
		&storage.Solved{}); err != nil {
		slog.Error("couldn't auto migrate", slogResp.Error(op, err))
		return nil, err
	}

	return &Storage{db: db}, nil
}
