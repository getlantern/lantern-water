// Package downloader provides a WASM downloader that can download the WASM
// file from a given URL. The downloader supports both HTTPS URLs and magnet links.
package downloader

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

//go:generate mockgen -package=downloader -destination=mocks.go . WASMDownloader,torrentClient,torrentInfo
//go:generate mockgen -package=downloader -destination=torrent_reader_mock_test.go github.com/anacrolix/torrent Reader

// WASMDownloader is an interface that defines the methods that a WASM downloader
type WASMDownloader interface {
	DownloadWASM(context.Context, io.Writer) error
	Close() error
}

type downloader struct {
	expectedHashSum string
	urls            []string
	httpClient      *http.Client
}

// NewWASMDownloader creates a new WASMDownloader instance.
func NewWASMDownloader(hashsum string, urls []string, httpClient *http.Client) (WASMDownloader, error) {
	if hashsum == "" {
		return nil, fmt.Errorf("missing required hashsum")
	}
	if len(urls) == 0 {
		return nil, fmt.Errorf("WASM downloader requires URLs to download but received empty list")
	}
	return &downloader{
		urls:            urls,
		httpClient:      httpClient,
		expectedHashSum: hashsum,
	}, nil
}

func (d *downloader) Close() error {
	return nil
}

// DownloadWASM downloads the WASM file from the given URLs, verifies the hash
// sum and writes the file to the given writer.
func (d *downloader) DownloadWASM(ctx context.Context, w io.Writer) error {
	joinedErrs := errors.New("failed to download WASM from all URLs")
	for _, url := range d.urls {
		tempBuffer := &bytes.Buffer{}
		if err := d.downloadWASM(ctx, tempBuffer, url); err != nil {
			joinedErrs = errors.Join(joinedErrs, err)
			continue
		}

		if err := d.verifyHashSum(tempBuffer.Bytes()); err != nil {
			joinedErrs = errors.Join(joinedErrs, err)
			continue
		}

		if _, err := tempBuffer.WriteTo(w); err != nil {
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
		return newHTTPSDownloader(d.httpClient, url).DownloadWASM(ctx, w)
	case strings.HasPrefix(url, "magnet:?"):
		downloader, err := newMagnetDownloader(ctx, d.httpClient, url)
		if err != nil {
			return err
		}
		defer downloader.Close()
		return downloader.DownloadWASM(ctx, w)
	default:
		return fmt.Errorf("unsupported protocol: %s", url)
	}
}

func (d *downloader) verifyHashSum(data []byte) error {
	got := fmt.Sprintf("%x", sha256.Sum256(data))
	if d.expectedHashSum != got {
		return fmt.Errorf("hashsum verification failed, expected %s, but got %s", d.expectedHashSum, got)
	}
	return nil
}
