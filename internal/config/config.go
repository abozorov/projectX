package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HttpHost        string `env:"HOST"`
	Storage         string `env:"STORAGE"`
	AuditLogStorage string `env:"AUDITLOGSTORAGE"`
}

func NewConfig(path string) (Config, error) {
	var cnf Config

	err := cleanenv.ReadConfig(path, &cnf)
	if err != nil {
		return Config{}, fmt.Errorf("cleanenv.ReadConfig: %w", err)
	}

	return cnf, nil
}
