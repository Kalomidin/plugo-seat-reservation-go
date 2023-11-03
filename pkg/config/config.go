package config

import (
	"seat-reservation/jwt"
)

type Config struct {
	Auth jwt.Config `yaml:"auth"`
	DB   Postgres   `yaml:"db"`
}
