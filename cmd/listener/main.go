package main

import (
	"bytes"
	"context"
	"flag"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	waterDownloader "github.com/getlantern/lantern-water/downloader"
	"github.com/getlantern/lantern-water/listener"
)

func main() {
	ctx := context.Background()
	log := slog.New(slog.NewTextHandler(os.Stdin, nil))

	var listenAddr, wasmAvailableAt, transportName string
	flag.StringVar(&listenAddr, "proxyURL", "localhost:8080", "URL of the proxy")
	flag.StringVar(&wasmAvailableAt, "wasmAvailableAt", "https://github.com/getlantern/wateringhole/raw/716a062ffa977fb4004d17827d46bc401265e2ac/protocols/plain/v1.0.0/plain.wasm", "URL where the WASM is available")
	flag.StringVar(&transportName, "transport", "plain", "Transport to use")
	flag.Parse()

	// Client for downloading WASM file
	cli := &http.Client{
		Timeout: 600 * time.Second,
	}
	downloader, err := waterDownloader.NewWASMDownloader(strings.Split(wasmAvailableAt, ","), cli)
	if err != nil {
		log.Error("failed to create wasm downloader", slog.Any("err", err))
		return
	}

	buffer := new(bytes.Buffer)
	err = downloader.DownloadWASM(ctx, buffer)
	if err != nil {
		log.Error("failed to retrieve WASM", slog.Any("err", err))
		return
	}
	wasm, err := io.ReadAll(buffer)
	if err != nil {
		log.Error("failed to read WASM", slog.Any("err", err))
	}

	l, err := listener.NewWATERListener(ctx, listener.ListenerParams{
		Transport: transportName,
		WASM:      wasm,
		Address:   listenAddr,
	})
	if err != nil {
		log.Error("failed to create listener", slog.Any("err", err))
		return
	}
	defer l.Close()

	go func() {
		conn, err := l.Accept()
		if err != nil {
			log.Error("failed to accept", slog.Any("err", err))
		}

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Error("failed to read", slog.Any("err", err))
		}

		slog.Info("received", slog.Any("data", string(buf[:n])))

		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 0\r\n\r\n"))
		if err != nil {
			log.Error("failed to write", slog.Any("err", err))
		}
	}()

	<-ctx.Done()
}
