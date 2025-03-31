package main

import (
	"context"
	"flag"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/getlantern/golog"
	"github.com/getlantern/lantern-water/dialer"
	waterDownloader "github.com/getlantern/lantern-water/downloader"
	"github.com/getlantern/lantern-water/logger"
	waterVC "github.com/getlantern/lantern-water/version_control"

	_ "github.com/refraction-networking/water/transport/v1"
)

func main() {
	ctx := context.Background()
	log := slog.New(slog.NewTextHandler(os.Stdin, nil))

	var listenerAddr, wasmAvailableAt, transportName string
	flag.StringVar(&listenerAddr, "proxyURL", "localhost:8080", "URL of the proxy")
	flag.StringVar(&wasmAvailableAt, "wasmAvailableAt", "https://github.com/getlantern/wateringhole/raw/716a062ffa977fb4004d17827d46bc401265e2ac/protocols/plain/v1.0.0/plain.wasm", "URL where the WASM is available")
	flag.StringVar(&transportName, "transport", "plain", "Transport to use")
	flag.Parse()

	storageDir, err := os.MkdirTemp("", "lantern-water-example")
	if err != nil {
		log.Error("failed to create storage dir", slog.Any("err", err))
		return
	}

	vc := waterVC.NewWaterVersionControl(storageDir, slog.New(logger.NewLogHandler(golog.LoggerFor("water-vc"), transportName)))

	// Client for downloading WASM file
	cli := &http.Client{
		Timeout: 600 * time.Second,
	}
	downloader, err := waterDownloader.NewWASMDownloader(strings.Split(wasmAvailableAt, ","), cli)
	if err != nil {
		log.Error("failed to create wasm downloader", slog.Any("err", err))
		return
	}
	wasmRC, err := vc.GetWASM(ctx, transportName, downloader)
	if err != nil {
		log.Error("failed to retrieve WASM", slog.Any("err", err))
		return
	}
	defer wasmRC.Close()

	wasm, err := io.ReadAll(wasmRC)
	if err != nil {
		log.Error("failed to read WASM", slog.Any("err", err))
		return
	}

	dialer, err := dialer.NewDialer(ctx, dialer.DialerParameters{Transport: transportName, WASM: wasm})
	if err != nil {
		log.Error("failed to create dialer", slog.Any("err", err))
		return
	}

	conn, err := dialer.DialContext(ctx, "tcp", listenerAddr)
	if err != nil {
		log.Error("failed to dial", slog.Any("err", err))
		return
	}
	defer conn.Close()
	_, err = conn.Write([]byte("Hello world!"))
	if err != nil {
		log.Error("failed to write", slog.Any("err", err))
		return
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Error("failed to read", slog.Any("err", err))
		return
	}
	log.Info("received", slog.Any("data", string(buf[:n])))
}
