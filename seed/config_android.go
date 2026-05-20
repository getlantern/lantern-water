//go:build android

package seed

import "github.com/anacrolix/torrent"

func applyPlatformConfig(cfg *torrent.ClientConfig) {
	// net.Interfaces() via netlinkrib is blocked by Android SELinux for
	// untrusted_app domains. UPnP is also useless on carrier NAT.
	cfg.NoDefaultPortForwarding = true
}
