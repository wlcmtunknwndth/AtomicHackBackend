package config

import (
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/wlcmtunknwndth/AtomicHackBackend/lib/slogResp"
	"log/slog"
	"os"
	"time"
)

const (
	configPathKey = "config_path"
)

var errNoConfigPath = errors.New("invalid config path")

type FileServer struct {
	StorageFolder string `yaml:"storage_folder" env-required:"true"`
	UrlPath       string `json:"url_path" env-default:"static"`
}

type Database struct {
	DbUser  string `yaml:"db_user" env-required:"true"`
	DbPass  string `yaml:"db_pass" env-required:"true"`
	DbName  string `yaml:"db_name" env-required:"true"`
	SslMode string `yaml:"ssl_mode" env-default:"disable"`
	Port    string `yaml:"db_port" env-default:"5432"`
}

type Server struct {
	Timeout     time.Duration `yaml:"timeout" env-default:"15s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"30s"`
	Address     string        `yaml:"address" env-required:"true"`
}

type Broker struct {
	MaxReconnects int           `yaml:"max_reconnects"`
	ReconnectWait time.Duration `yaml:"reconnect_wait"`
	Address       string        `yaml:"address"`
	Retry         bool          `yaml:"retry"`
}

type Config struct {
	Db         Database   `yaml:"db"`
	Server     Server     `yaml:"server"`
	Broker     Broker     `yaml:"broker"`
	FileServer FileServer `yaml:"file_server"`
}

func MustLoad() *Config {
	const op = "internal.config.MustLoad"

	path, ok := os.LookupEnv(configPathKey)
	if !ok || path == "" {
		slog.Error("couldn't find config path env", slogResp.Error(op, errNoConfigPath))
		panic(errNoConfigPath.Error())
	}

	if _, err := os.Stat(path); err != nil {
		slog.Error("couldn't find config path", slogResp.Error(op, err))
		panic("couldn't find config path")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		slog.Error("couldn't read config", slogResp.Error(op, err))
		panic("couldn't read config")
	}
	return &cfg
}
