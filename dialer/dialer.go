package dialer

import (
	"context"
	"log/slog"

	"github.com/getlantern/golog"
	"github.com/getlantern/lantern-water/logger"
	"github.com/refraction-networking/water"
)

type DialerParameters struct {
	Logger    golog.Logger
	Transport string
	WASM      []byte
}

func NewDialer(ctx context.Context, params DialerParameters) (water.Dialer, error) {
	cfg := &water.Config{
		TransportModuleBin: params.WASM,
		OverrideLogger:     slog.New(logger.NewLogHandler(params.Logger, params.Transport)),
	}

	dialer, err := water.NewDialerWithContext(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return dialer, nil
}
