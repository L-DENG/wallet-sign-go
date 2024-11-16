package flags

import (
	"github.com/urfave/cli/v2"
)

const evnVarPrefix = "SIGNATURE_"

func prefixEnvVars(name string) []string {
	return []string{evnVarPrefix + name}
}

var (
	LevelDbPathFlag = &cli.StringFlag{
		Name:    "leveldb-path",
		Usage:   "the path of leveldb",
		EnvVars: prefixEnvVars("LEVEL_DB_PATH"),
		Value:   "./leveldb",
	}

	RpcHostFlag = &cli.StringFlag{
		Name:     "rpc-host",
		Usage:    "The host of the rpc server",
		EnvVars:  prefixEnvVars("RPC_HOST"),
		Value:    "localhost",
		Required: true,
	}

	RpcPortFlag = &cli.IntFlag{
		Name:     "rpc-port",
		Usage:    "The port of the rpc server",
		EnvVars:  prefixEnvVars("RPC_PORT"),
		Value:    8989,
		Required: true,
	}
)

var requiredFlags = []cli.Flag{
	RpcPortFlag,
	RpcHostFlag,
	LevelDbPathFlag,
}

var optionalFlags = []cli.Flag{}

func init() {
	Flags = append(requiredFlags, optionalFlags...)
}

var Flags []cli.Flag
