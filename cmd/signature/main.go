package main

import (
	"context"
	"github.com/L-DENG/wallet-sign-go/common/opio"
	"github.com/ethereum/go-ethereum/log"
	"os"
)

var (
	GitCommit = ""
	GitData   = ""
)

func main() {
	log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stderr, log.LevelInfo, true)))
	app := NewCli(GitCommit, GitData)
	ctx := opio.WithInterruptBlocker(context.Background())
	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Error("Application failed", "error", err.Error())
		os.Exit(1)
	}
}
