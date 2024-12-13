package listener

import (
	"bytes"
	"context"
	"embed"
	"io"
	"net"
	"testing"

	"github.com/getlantern/golog"
	"github.com/refraction-networking/water"
	_ "github.com/refraction-networking/water/transport/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/*
var testData embed.FS

func TestWATERListener(t *testing.T) {
	f, err := testData.Open("testdata/reverse_v1.wasm")
	require.Nil(t, err)

	wasm, err := io.ReadAll(f)
	require.Nil(t, err)

	ctx := context.Background()

	// configuration used at the dialer
	cfg := &water.Config{
		TransportModuleBin: wasm,
	}

	listenerParameters := ListenerParams{
		Logger:    golog.LoggerFor("water"),
		Transport: "reverse_v1",
		Address:   "127.0.0.1:3000",
		WASM:      wasm,
	}

	ll, err := NewWATERListener(ctx, listenerParameters)
	require.Nil(t, err)

	messageRequest := "hello"
	expectedResponse := "world"
	// running listener
	go func() {
		for {
			var conn net.Conn
			conn, err = ll.Accept()
			if err != nil {
				t.Error(err)
				return
			}

			go func() {
				if conn == nil {
					t.Error("nil connection")
					return
				}
				buf := make([]byte, 2*len(messageRequest))
				n, connErr := conn.Read(buf)
				if connErr != nil {
					t.Errorf("error reading: %v", err)
					return
				}

				buf = buf[:n]
				if !bytes.Equal(buf, []byte(messageRequest)) {
					t.Errorf("unexpected request %v %v", buf, messageRequest)
					return
				}
				conn.Write([]byte(expectedResponse))
			}()
		}
	}()

	dialer, err := water.NewDialerWithContext(ctx, cfg)
	require.Nil(t, err)

	conn, err := dialer.DialContext(ctx, "tcp", ll.Addr().String())
	require.Nil(t, err)
	defer conn.Close()

	n, err := conn.Write([]byte(messageRequest))
	assert.Nil(t, err)
	assert.Equal(t, len(messageRequest), n)

	buf := make([]byte, 1024)
	n, err = conn.Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, len(expectedResponse), n)
	assert.Equal(t, expectedResponse, string(buf[:n]))
}
