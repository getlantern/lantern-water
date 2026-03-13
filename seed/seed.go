// Package seed provides a BitTorrent seeder for WASM files.
// It builds a metainfo and magnet URI from a local file and seeds it to peers.
package seed

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
)

const defaultPieceLength = 256 * 1024 // 256 KiB

// Seeder seeds a WASM file via BitTorrent.
type Seeder struct {
	client    *torrent.Client
	magnetURI string
}

// New creates a Seeder for the file at filePath, begins seeding it, and
// returns the Seeder alongside the generated magnet URI.
// It uses the default options from anacrolix/torrent config and make sure
// it enables seed options based on the announce list and enable DHT
func New(filePath string, announceList [][]string) (*Seeder, error) {
	mi, err := buildMetainfo(filePath, announceList)
	if err != nil {
		return nil, fmt.Errorf("building metainfo: %w", err)
	}

	magnet, err := mi.MagnetV2()
	if err != nil {
		return nil, fmt.Errorf("building magnet URI: %w", err)
	}
	magnetURI := magnet.String()

	cfg := torrent.NewDefaultClientConfig()
	path, err := os.MkdirTemp(filepath.Dir(filePath), "lantern-water-torrent")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	cfg.DataDir = path
	cfg.Seed = true
	cfg.NoDHT = false
	cfg.NoUpload = false
	cfg.AcceptPeerConnections = true

	client, err := torrent.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating torrent client: %w", err)
	}

	t, err := client.AddTorrent(mi)
	if err != nil {
		if errs := client.Close(); errs != nil {
			closeErr := errors.Join(errs...)
			err = errors.Join(err, closeErr)
		}
		return nil, fmt.Errorf("adding torrent: %w", err)
	}

	<-t.GotInfo()
	t.DisallowDataDownload()

	return &Seeder{client: client, magnetURI: magnetURI}, nil
}

// MagnetURI returns the magnet URI for the seeded file.
func (s *Seeder) MagnetURI() string {
	return s.magnetURI
}

// Close stops seeding and shuts down the torrent client.
func (s *Seeder) Close() error {
	if errs := s.client.Close(); len(errs) > 0 {
		return fmt.Errorf("closing torrent client: %w", errors.Join(errs...))
	}
	return nil
}

// buildMetainfo creates a MetaInfo for the file at filePath.
func buildMetainfo(filePath string, announceList [][]string) (*metainfo.MetaInfo, error) {
	info := metainfo.Info{
		PieceLength: defaultPieceLength,
	}
	if err := info.BuildFromFilePath(filePath); err != nil {
		return nil, fmt.Errorf("building info from file path: %w", err)
	}

	infoBytes, err := bencode.Marshal(info)
	if err != nil {
		return nil, fmt.Errorf("marshaling info: %w", err)
	}

	mi := &metainfo.MetaInfo{
		InfoBytes:    infoBytes,
		AnnounceList: metainfo.AnnounceList(announceList),
	}
	mi.SetDefaults()

	return mi, nil
}
