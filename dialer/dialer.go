// Package dialer holds the dialer implementation for the water transport.
package dialer

import (
	"context"
	"log/slog"

	"github.com/getlantern/golog"
	"github.com/getlantern/lantern-water/logger"
	"github.com/refraction-networking/water"
	_ "github.com/refraction-networking/water/transport/v1"
)

// DialerParameters are used when creating a new dialer.
type DialerParameters struct {
	// An optional golog.Logger used for keeping compatibility with http-proxy
	// and flashlight logger. If not defined the dialer will use the default
	// water logger.
	Logger    golog.Logger
	Transport string // Specifies transport being used.
	WASM      []byte // The WASM module to use.
}

// NewDialer creates a new water dialer with the given parameters.
func NewDialer(ctx context.Context, params DialerParameters) (water.Dialer, error) {
	cfg := &water.Config{
		TransportModuleBin: params.WASM,
	}

	if params.Logger != nil {
		cfg.OverrideLogger = slog.New(logger.NewLogHandler(params.Logger, params.Transport))
	}

	dialer, err := water.NewDialerWithContext(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return dialer, nil
}
