// Package listener contains the WATER listener creation functions
package listener

import (
	"context"
	"log/slog"
	"net"

	"github.com/refraction-networking/water"
	_ "github.com/refraction-networking/water/transport/v1"
)

// ListenerParams contain arguments/parameters used for creating a new WATER listener
type ListenerParams struct {
	// BaseListener is a listener that should be wrapped by the WATER listener, it's optional and can be nil
	BaseListener net.Listener
	// An optional *slog.Logger. If not defined the listener will use the default
	// water logger.
	Logger *slog.Logger
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
	}

	if params.Logger != nil {
		logger := params.Logger
		if params.Transport != "" {
			logger = logger.With("transport", params.Transport)
		}
		cfg.OverrideLogger = logger
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
