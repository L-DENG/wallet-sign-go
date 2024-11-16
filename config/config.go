package config

import (
	"github.com/L-DENG/wallet-sign-go/flags"
	"github.com/urfave/cli/v2"
)

type Config struct {
	LevelDBPath string
	RpcServer   ServerConfig
}

type ServerConfig struct {
	Host string
	Port int
}

func NewConfig(ctx *cli.Context) Config {
	return Config{
		LevelDBPath: ctx.String(flags.LevelDbPathFlag.Name),
		RpcServer: ServerConfig{
			Host: ctx.String(flags.RpcHostFlag.Name),
			Port: ctx.Int(flags.RpcPortFlag.Name),
		},
	}
}
