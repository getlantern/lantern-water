package seed

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSeeder(t *testing.T) {
	seed, err := New("testdata/shadowsocks_client.wasm")
	require.NoError(t, err)
	defer seed.Close()
	t.Logf("Magnet URI: %s", seed.MagnetURI())
}
