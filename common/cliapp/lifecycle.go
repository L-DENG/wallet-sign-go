package cliapp

import (
	"context"
	"errors"
	"fmt"
	"github.com/L-DENG/wallet-sign-go/common/opio"
	"github.com/urfave/cli/v2"
	"os"
)

type Lifecycle interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Stopped() bool
}

type LifecycleAction func(ctx *cli.Context, close context.CancelCauseFunc) (Lifecycle, error)

func LifecycleCmd(fn LifecycleAction) cli.ActionFunc {
	return lifecycleCmd(fn, opio.BlockOnInterruptsContext)
}

var interruptErr = errors.New("interrupt signal ...")

type WaitSignalFn func(ctx context.Context, signal ...os.Signal)

func lifecycleCmd(fn LifecycleAction, blockFnOnInterrupt WaitSignalFn) cli.ActionFunc {

	return func(ctx *cli.Context) error {
		hostCtx := ctx.Context
		appCtx, appCancel := context.WithCancelCause(hostCtx)
		ctx.Context = appCtx

		go func() {
			blockFnOnInterrupt(appCtx)
			appCancel(interruptErr)
		}()
		appLifecycle, err := fn(ctx, appCancel)
		if err != nil {
			return errors.Join(
				fmt.Errorf("failed to setup: %W", err),
				context.Cause(appCtx),
			)
		}
		if err := appLifecycle.Start(appCtx); err != nil {
			return errors.Join(
				fmt.Errorf("failed to start: %W", err),
				context.Cause(appCtx),
			)
		}

		<-appCtx.Done()

		stopCtx, StopCause := context.WithCancelCause(hostCtx)
		go func() {
			blockFnOnInterrupt(stopCtx)
			StopCause(interruptErr)
		}()

		stopErr := appLifecycle.Stop(stopCtx)
		if err != nil {
			return errors.Join(
				fmt.Errorf("failed to stop: %W", stopErr),
				context.Cause(stopCtx),
			)
		}
		return nil
	}
}
