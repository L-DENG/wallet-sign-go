package main

import (
	"context"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/urfave/cli/v2"

	"github.com/L-DENG/wallet-sign-go/common/cliapp"
	"github.com/L-DENG/wallet-sign-go/config"
	"github.com/L-DENG/wallet-sign-go/flags"
	"github.com/L-DENG/wallet-sign-go/leveldb"
	"github.com/L-DENG/wallet-sign-go/services/rpc"
)

func runRpc(ctx *cli.Context, shutDown context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	log.Info("Runing grpc services...")
	cfg := config.NewConfig(ctx)
	db, err := leveldb.NewKeys(cfg.LevelDBPath)
	if err != nil {
		log.Error("Failed to open leveldb", "err", err)
		return nil, err
	}
	log.Info("port:", ctx.Int(flags.RpcPortFlag.Name), "Host:", ctx.String(flags.RpcHostFlag.Name))
	log.Info("port:", cfg.RpcServer.Port, "Host:", cfg.RpcServer.Host)

	grpcServer := &rpc.RpcServerConfig{
		GrpcHostName: cfg.RpcServer.Host,
		GrpcPort:     cfg.RpcServer.Port,
	}
	return rpc.NewRpcServer(db, grpcServer)
}

func NewCli(GitCommit string, GitData string) *cli.App {
	flags := flags.Flags
	return &cli.App{
		Version:              params.VersionWithCommit(GitCommit, GitData),
		Description:          "A generate key pair servers(ssm)",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:        "rpc",
				Flags:       flags,
				Description: "Run rpc service",
				Action:      cliapp.LifecycleCmd(runRpc),
			},
			{
				Name:        "version",
				Description: "Show project version",
				Action: func(ctx *cli.Context) error {
					cli.ShowVersion(ctx)
					return nil
				},
			},
		},
	}
}
