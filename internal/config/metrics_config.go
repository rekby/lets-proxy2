// need for prevent import loop

package config

import (
	"github.com/rekby/lets-proxy2/internal/secrethandler"
	"github.com/rekby/lets-proxy2/internal/tlslistener"
)

type listenConfig tlslistener.Config
type secretHandlerConfig secrethandler.Config

type Config struct {
	Enable bool

	listenConfig
	secretHandlerConfig
}

func (c Config) GetListenConfig() tlslistener.Config {
	return tlslistener.Config(c.listenConfig)
}

func (c Config) GetSecretHandlerConfig() secrethandler.Config {
	return secrethandler.Config(c.secretHandlerConfig)
}
