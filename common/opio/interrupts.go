package opio

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var DefaultInterruptSignals = []os.Signal{
	os.Interrupt,
	os.Kill,
	syscall.SIGTERM,
	syscall.SIGQUIT,
}

func BlockOnInterrupts(signals ...os.Signal) {
	if len(signals) == 0 {
		signals = DefaultInterruptSignals
	}
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, signals...)
	<-interruptChannel
}

func BlockOnInterruptsContext(ctx context.Context, signals ...os.Signal) {
	if len(signals) == 0 {
		signals = DefaultInterruptSignals
	}
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, signals...)
	select {
	case <-interruptChannel:
	case <-ctx.Done():
		signal.Stop(interruptChannel)
	}
}

type interruptContextKeyType struct{}

var blockContextKey = interruptContextKeyType{}

type interruptCatcher struct {
	incoming chan os.Signal
}

func (c *interruptCatcher) Block(ctx context.Context) {
	select {
	case <-c.incoming:
	case <-ctx.Done():
	}
}

func WithInterruptBlocker(ctx context.Context) context.Context {
	if ctx.Value(blockContextKey) != nil {
		return ctx
	}
	catcher := &interruptCatcher{incoming: make(chan os.Signal, 10)}
	signal.Notify(catcher.incoming, DefaultInterruptSignals...)
	return context.WithValue(ctx, blockContextKey, BlockFn(catcher.Block))
}

func WithBlocker(ctx context.Context, fn BlockFn) context.Context {
	return context.WithValue(ctx, blockContextKey, fn)
}

type BlockFn func(ctx context.Context)

func BlockFnFromContext(ctx context.Context) BlockFn {
	v := ctx.Value(blockContextKey)
	if v == nil {
		return nil
	}
	return v.(BlockFn)
}

func CancelOnInterrupt(ctx context.Context) context.Context {
	inner, cancel := context.WithCancel(context.Background())

	blockFnOnInterrupt := BlockFnFromContext(ctx)
	if blockFnOnInterrupt != nil {
		blockFnOnInterrupt = func(ctx context.Context) {
			BlockOnInterruptsContext(ctx)
		}
	}

	go func() {
		blockFnOnInterrupt(ctx)
		cancel()
	}()

	return inner
}
