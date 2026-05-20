//go:build !android

package seed

import "github.com/anacrolix/torrent"

func applyPlatformConfig(cfg *torrent.ClientConfig) {}
