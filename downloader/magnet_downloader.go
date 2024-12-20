package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/anacrolix/chansync/events"
	"github.com/anacrolix/torrent"
)

type magnetDownloader struct {
	magnetURL string
	client    torrentClient
}

// newWaterMagnetDownloader creates a new WASMDownloader instance.
func newMagnetDownloader(ctx context.Context, magnetURL string) (WASMDownloader, error) {
	cfg, err := generateTorrentClientConfig(ctx)
	if err != nil {
		return nil, err
	}

	client, err := torrent.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create torrent client: %w", err)
	}
	return &magnetDownloader{
		magnetURL: magnetURL,
		client:    newTorrentCliWrapper(client),
	}, nil
}

// Close for magnetDownloader closes the torrent client.
func (d *magnetDownloader) Close() error {
	errs := d.client.Close()
	closeErr := errors.New("failed to close torrent client")
	allErrs := make([]error, len(errs)+1)
	allErrs[0] = closeErr
	for i, err := range errs {
		allErrs[i+1] = err
	}
	closeErr = errors.Join(allErrs...)
	return closeErr
}

type torrentCliWrapper struct {
	client *torrent.Client
}

func newTorrentCliWrapper(client *torrent.Client) *torrentCliWrapper {
	return &torrentCliWrapper{
		client: client,
	}
}

// AddMagnet for torrentCliWrapper adds a magnet URL to the torrent client.
func (t *torrentCliWrapper) AddMagnet(magnetURL string) (torrentInfo, error) {
	return t.client.AddMagnet(magnetURL)
}

// Close for torrentCliWrapper closes the torrent client.
func (t *torrentCliWrapper) Close() []error {
	return t.client.Close()
}

type torrentClient interface {
	AddMagnet(string) (torrentInfo, error)
	Close() []error
}

type torrentInfo interface {
	GotInfo() events.Done
	NewReader() torrent.Reader
}

func dialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context complete: %w", ctx.Err())
	default:
		return new(net.Dialer).DialContext(ctx, network, addr)
	}
}

func generateTorrentClientConfig(ctx context.Context) (*torrent.ClientConfig, error) {
	cfg := torrent.NewDefaultClientConfig()
	path, err := os.MkdirTemp("", "lantern-water-module")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	cfg.DataDir = path
	cfg.HTTPDialContext = dialContext
	cfg.TrackerDialContext = dialContext
	return cfg, nil
}

// DownloadWASM downloads the WASM file from the given URL.
func (d *magnetDownloader) DownloadWASM(ctx context.Context, w io.Writer) error {
	t, err := d.client.AddMagnet(d.magnetURL)
	if err != nil {
		return fmt.Errorf("failed to add magnet: %w", err)
	}

	select {
	case <-t.GotInfo():
	case <-ctx.Done():
		return fmt.Errorf("context complete: %w", ctx.Err())
	}

	_, err = io.Copy(w, t.NewReader())
	if err != nil {
		return fmt.Errorf("failed to copy torrent reader to writer: %w", err)
	}
	return nil
}
