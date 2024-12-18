package config

import (
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
)

func NewKoanf() *koanf.Koanf {
	k := koanf.New(".")
	err := k.Load(file.Provider("config.json"), json.Parser())
	if err != nil {
		log.Fatal().Msg("Failed to load config.json")
	}
	return k
}
