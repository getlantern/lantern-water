// Package seed provides a BitTorrent seeder for WASM files.
// It builds a metainfo and magnet URI from a local file and seeds it to peers.
package seed

import (
	"errors"
	"fmt"
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
func New(filePath string) (*Seeder, error) {
	mi, err := buildMetainfo(filePath)
	if err != nil {
		return nil, fmt.Errorf("building metainfo: %w", err)
	}

	magnet, err := mi.MagnetV2()
	if err != nil {
		return nil, fmt.Errorf("building magnet URI: %w", err)
	}
	magnetURI := magnet.String()

	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = filepath.Dir(filePath)
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
		client.Close()
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
func buildMetainfo(filePath string) (*metainfo.MetaInfo, error) {
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
		InfoBytes: infoBytes,
		AnnounceList: metainfo.AnnounceList{
			{"udp://tracker.opentrackr.org:1337/announce"},
			{"udp://open.demonii.com:1337/announce"},
			{"udp://open.stealth.si:80/announce"},
			{"udp://exodus.desync.com:6969/announce"},
			{"https://torrent.tracker.durukanbal.com:443/announce"},
			{"udp://tracker.torrent.eu.org:451/announce"},
			{"udp://tracker.theoks.net:6969/announce"},
			{"udp://tracker.srv00.com:6969/announce"},
			{"udp://tracker.filemail.com:6969/announce"},
			{"udp://tracker.dler.org:6969/announce"},
			{"udp://tracker.corpscorp.online:80/announce"},
			{"udp://tracker.alaskantf.com:6969/announce"},
			{"udp://tracker-udp.gbitt.info:80/announce"},
			{"udp://t.overflow.biz:6969/announce"},
			{"udp://open.dstud.io:6969/announce"},
			{"udp://leet-tracker.moe:1337/announce"},
			{"udp://explodie.org:6969/announce"},
			{"udp://bittorrent-tracker.e-n-c-r-y-p-t.net:1337/announce"},
			{"udp://6ahddutb1ucc3cp.ru:6969/announce"},
		},
	}
	mi.SetDefaults()

	return mi, nil
}
