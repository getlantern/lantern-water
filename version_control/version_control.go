// Package version_control provides a version control system for the WASM files
package version_control

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/getlantern/lantern-water/downloader"
)

type waterVersionControl struct {
	dir    string
	logger *slog.Logger
}

type wasmInfo struct {
	lastTimeLoaded time.Time
	path           string
}

// NewWaterVersionControl creates a new instance of the version control system.
// It requires a directory where the WASM files will be stored and a logger.
func NewWaterVersionControl(dir string, logger *slog.Logger) *waterVersionControl {
	return &waterVersionControl{
		dir:    dir,
		logger: logger,
	}
}

// GetWASM returns the WASM file for the given transport.
// Please remember to Close the io.ReadCloser after using it.
// This function implements the following steps:
// 1. Check if the WASM file exists
// 2. If it does not exist, download it
// 3. If it exists, check if it was loaded correctly by checking if the last-loaded file exists
// 4. If it was not loaded correctly or the last-loaded file doesn't exist, download it again
// 5. If it was loaded correctly, return the file and mark the file as loaded
// 6. It deletes the WASM files that were not used for more than 7 days after successful loading
func (vc *waterVersionControl) GetWASM(ctx context.Context, transport string, downloader downloader.WASMDownloader) (io.ReadCloser, error) {
	path := filepath.Join(vc.dir, transport+".wasm")
	f, err := os.Open(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}

	if errors.Is(err, fs.ErrNotExist) || f == nil {
		if f != nil {
			f.Close()
		}
		response, err := vc.downloadWASM(ctx, transport, downloader)
		if err != nil {
			return nil, fmt.Errorf("failed to download WASM file: %w", err)
		}

		return response, nil
	}

	lastLoaded, err := os.Open(filepath.Join(vc.dir, transport+".last-loaded"))
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to open file %s: %w", transport+".last-loaded", err)
	}

	if errors.Is(err, fs.ErrNotExist) {
		// WASM file exists but it was never loaded correctly, downloading it again
		response, err := vc.downloadWASM(ctx, transport, downloader)
		if err != nil {
			return nil, fmt.Errorf("failed to download WASM file: %w", err)
		}
		return response, nil
	}
	defer lastLoaded.Close()

	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to seek file at the beginning: %w", err)
	}

	if err = vc.markUsed(transport); err != nil {
		return nil, fmt.Errorf("failed to update WASM history: %w", err)
	}
	return f, nil
}

// markUsed updates the last-loaded file for the given transport
func (vc *waterVersionControl) markUsed(transport string) error {
	f, err := os.Create(filepath.Join(vc.dir, transport+".last-loaded"))
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", transport+".last-loaded", err)
	}
	defer f.Close()

	if _, err = f.WriteString(strconv.FormatInt(time.Now().UTC().Unix(), 10)); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", transport+".last-loaded", err)
	}
	if err = vc.cleanOutdated(); err != nil {
		return fmt.Errorf("failed to clean outdated WASMs: %w", err)
	}
	return nil
}

// unusedWASMsDeletedAfter is the time after which the WASM files are considered outdated
const unusedWASMsDeletedAfter = 7 * 24 * time.Hour

func (vc *waterVersionControl) cleanOutdated() error {
	wg := new(sync.WaitGroup)
	filesToBeDeleted := make([]string, 0)
	// walk through dir, load last-loaded and delete if older than unusedWASMsDeletedAfter
	err := filepath.Walk(vc.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk through dir: %w", err)
		}
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".last-loaded" {
			return nil
		}

		lastLoaded, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		i, err := strconv.ParseInt(string(lastLoaded), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse int: %w", err)
		}
		lastLoadedTime := time.Unix(i, 0)
		if time.Since(lastLoadedTime) > unusedWASMsDeletedAfter {
			filesToBeDeleted = append(filesToBeDeleted, path)
		}
		return nil
	})
	for _, path := range filesToBeDeleted {
		wg.Add(1)
		go func() {
			defer wg.Done()
			transport := strings.TrimSuffix(filepath.Base(path), ".last-loaded")
			if err = os.Remove(filepath.Join(vc.dir, transport+".wasm")); err != nil {
				vc.logger.Error("failed to remove wasm file", slog.String("file", transport+".wasm"), slog.Any("err", err))
				return
			}
			if err = os.Remove(path); err != nil {
				vc.logger.Error("failed to remove last-loaded file", slog.String("path", path), slog.Any("err", err))
				return
			}
		}()
	}
	wg.Wait()
	return err
}

func (vc *waterVersionControl) downloadWASM(ctx context.Context, transport string, downloader downloader.WASMDownloader) (io.ReadCloser, error) {
	outputPath := filepath.Join(vc.dir, transport+".wasm")
	f, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %w", transport, err)
	}

	if err = downloader.DownloadWASM(ctx, f); err != nil {
		return nil, fmt.Errorf("failed to download wasm: %w", err)
	}

	if _, err = f.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek to the beginning of the file: %w", err)
	}

	if err = vc.markUsed(transport); err != nil {
		return nil, fmt.Errorf("failed to update WASM history: %w", err)
	}

	return f, nil
}
