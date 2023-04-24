package graceful

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func HandleSignals(ctx context.Context, stopFunc ...func(context.Context)) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	<-signals
	for _, fn := range stopFunc {
		fn(ctx)
	}
}
