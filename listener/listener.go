package listener

import (
	"context"
	"log/slog"
	"net"

	"github.com/getlantern/golog"
	"github.com/getlantern/lantern-water/logger"
	"github.com/refraction-networking/water"
)

// ListenerParams contain arguments/parameters used for creating a new WATER listener
type ListenerParams struct {
	// baseListener is a listener that should be wrapped by the WATER listener, it's optional and can be nil
	BaseListener net.Listener
	// Logger is a golog logger that should be used internally by the WASM
	Logger golog.Logger
	// Transport represents the protocol, version or whatever detail that will
	// be used at local logs to help understanding which WASM file is being used
	Transport string

	// Address represents the address used by the listener
	Address string
	// WASM must contain the WASM data used by the WATER listener
	WASM []byte
}

// NewWATERListener creates a WATER listener
// Currently water doesn't support customized TCP connections and we need to listen and receive requests directly from the WATER listener
func NewWATERListener(ctx context.Context, params ListenerParams) (net.Listener, error) {
	cfg := &water.Config{
		TransportModuleBin: params.WASM,
		OverrideLogger:     slog.New(logger.NewLogHandler(params.Logger, params.Transport)),
	}

	if params.BaseListener != nil {
		cfg.NetworkListener = params.BaseListener
	}

	waterListener, err := cfg.ListenContext(ctx, "tcp", params.Address)
	if err != nil {
		return nil, err
	}

	return waterListener, nil
}
