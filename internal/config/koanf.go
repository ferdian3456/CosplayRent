package config

import (
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
)

func NewKoanf() *koanf.Koanf {
	k := koanf.New(".")
	err := k.Load(file.Provider("../.env"), dotenv.Parser())
	if err != nil {
		log.Fatal().Msg("Failed to load .env files")
	}
	return k

}
