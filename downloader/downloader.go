// Package downloader provides a WASM downloader that can download the WASM
// file from a given URL. The downloader supports both HTTPS URLs and magnet links.
package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

//go:generate mockgen -package=downloader -destination=mocks_test.go . WASMDownloader,torrentClient,torrentInfo
//go:generate mockgen -package=downloader -destination=torrent_reader_mock_test.go github.com/anacrolix/torrent Reader

// WASMDownloader is an interface that defines the methods that a WASM downloader
type WASMDownloader interface {
	DownloadWASM(context.Context, io.Writer) error
	Close() error
}

type downloader struct {
	urls             []string
	httpClient       *http.Client
	httpDownloader   WASMDownloader
	magnetDownloader WASMDownloader
}

// NewWaterWASMDownloader creates a new WASMDownloader instance.
func NewWASMDownloader(urls []string, httpClient *http.Client) (WASMDownloader, error) {
	if len(urls) == 0 {
		return nil, fmt.Errorf("WASM downloader requires URLs to download but received empty list")
	}
	return &downloader{
		urls:       urls,
		httpClient: httpClient,
	}, nil
}

func (d *downloader) Close() error {
	if d.magnetDownloader != nil {
		return d.magnetDownloader.Close()
	}
	return nil
}

// DownloadWASM downloads the WASM file from the given URLs, verifies the hash
// sum and writes the file to the given writer.
func (d *downloader) DownloadWASM(ctx context.Context, w io.Writer) error {
	joinedErrs := errors.New("failed to download WASM from all URLs")
	for _, url := range d.urls {
		err := d.downloadWASM(ctx, w, url)
		if err != nil {
			joinedErrs = errors.Join(joinedErrs, err)
			continue
		}

		return nil
	}
	return joinedErrs
}

// downloadWASM checks what kind of URL was given and downloads the WASM file
// from the URL. It can be a HTTPS URL or a magnet link.
func (d *downloader) downloadWASM(ctx context.Context, w io.Writer, url string) error {
	switch {
	case strings.HasPrefix(url, "http://"), strings.HasPrefix(url, "https://"):
		if d.httpDownloader == nil {
			d.httpDownloader = newHTTPSDownloader(d.httpClient, url)
		}
		return d.httpDownloader.DownloadWASM(ctx, w)
	case strings.HasPrefix(url, "magnet:?"):
		if d.magnetDownloader == nil {
			var err error
			downloader, err := newMagnetDownloader(ctx, url)
			if err != nil {
				return err
			}
			d.magnetDownloader = downloader
		}
		return d.magnetDownloader.DownloadWASM(ctx, w)
	default:
		return fmt.Errorf("unsupported protocol: %s", url)
	}
}
