package seed

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSeeder(t *testing.T) {
	seed, err := New("testdata/shadowsocks_client.wasm", [][]string{
		{"udp://tracker.opentrackr.org:1337/announce"},
	})
	require.NoError(t, err)
	defer seed.Close()
	t.Logf("Magnet URI: %s", seed.MagnetURI())
}
