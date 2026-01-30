package downloader

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestNewWASMDownloader(t *testing.T) {
	hashsum := "hashsum"
	client := http.DefaultClient
	urls := []string{"http://example.com"}
	var tests = []struct {
		name            string
		givenHashSum    string
		givenURLs       []string
		givenHTTPClient *http.Client
		assert          func(*testing.T, WASMDownloader, error)
	}{
		{
			name: "it should return an error when providing an empty hash sum",
			assert: func(t *testing.T, d WASMDownloader, err error) {
				assert.Error(t, err)
				assert.Nil(t, d)
			},
		},
		{
			name:         "it should return an error when providing an empty list of URLs",
			givenHashSum: hashsum,
			assert: func(t *testing.T, d WASMDownloader, err error) {
				assert.Error(t, err)
				assert.Nil(t, d)
			},
		},
		{
			name:            "it should successfully return a wasm downloader",
			givenHashSum:    hashsum,
			givenURLs:       urls,
			givenHTTPClient: client,
			assert: func(t *testing.T, wDownloader WASMDownloader, err error) {
				assert.NoError(t, err)
				d := wDownloader.(*downloader)
				assert.Equal(t, hashsum, d.expectedHashSum)
				assert.Equal(t, urls, d.urls)
				assert.Equal(t, client, d.httpClient)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := NewWASMDownloader(tt.givenHashSum, tt.givenURLs, tt.givenHTTPClient)
			tt.assert(t, d, err)
		})
	}
}

func TestDownloadWASM(t *testing.T) {
	ctx := context.Background()

	contentMessage := "content"
	hashsum := fmt.Sprintf("%x", sha256.Sum256([]byte(contentMessage)))
	var tests = []struct {
		name                string
		givenHashSum        string
		givenHTTPClient     *http.Client
		givenURLs           []string
		givenWriter         io.Writer
		setupHTTPDownloader func(ctrl *gomock.Controller) WASMDownloader
		assert              func(*testing.T, io.Reader, error)
	}{
		{
			name:         "udp urls are unsupported",
			givenHashSum: hashsum,
			givenURLs: []string{
				"udp://example.com",
			},
			assert: func(t *testing.T, r io.Reader, err error) {
				b, berr := io.ReadAll(r)
				require.NoError(t, berr)
				assert.Empty(t, b)
				assert.Error(t, err)
				assert.ErrorContains(t, err, "unsupported protocol")
			},
		},
		{
			name:         "http download error",
			givenHashSum: hashsum,
			givenURLs: []string{
				"http://example.com",
			},
			setupHTTPDownloader: func(ctrl *gomock.Controller) WASMDownloader {
				httpDownloader := NewMockWASMDownloader(ctrl)
				httpDownloader.EXPECT().DownloadWASM(ctx, gomock.Any()).Return(assert.AnError)
				return httpDownloader
			},
			assert: func(t *testing.T, r io.Reader, err error) {
				b, berr := io.ReadAll(r)
				require.NoError(t, berr)
				assert.Empty(t, b)
				assert.Error(t, err)
				assert.ErrorContains(t, err, assert.AnError.Error())
				assert.ErrorContains(t, err, "failed to download WASM from all URLs")
			},
		},
		{
			name:         "success",
			givenHashSum: hashsum,
			givenURLs: []string{
				"http://example.com",
			},
			setupHTTPDownloader: func(ctrl *gomock.Controller) WASMDownloader {
				httpDownloader := NewMockWASMDownloader(ctrl)
				httpDownloader.EXPECT().DownloadWASM(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, w io.Writer) error {
						_, err := w.Write([]byte(contentMessage))
						return err
					})
				return httpDownloader
			},
			assert: func(t *testing.T, r io.Reader, err error) {
				b, berr := io.ReadAll(r)
				require.NoError(t, berr)
				assert.NoError(t, err)
				assert.Equal(t, contentMessage, string(b))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var httpDownloader WASMDownloader
			if tt.setupHTTPDownloader != nil {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				httpDownloader = tt.setupHTTPDownloader(ctrl)
			}

			b := &bytes.Buffer{}
			d, err := NewWASMDownloader(tt.givenHashSum, tt.givenURLs, tt.givenHTTPClient)
			require.NoError(t, err)

			if httpDownloader != nil {
				d.(*downloader).httpDownloader = httpDownloader
			}
			err = d.DownloadWASM(ctx, b)
			tt.assert(t, b, err)
		})
	}
}
