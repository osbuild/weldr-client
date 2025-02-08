package cloud

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDistroList(t *testing.T) {
	json := `{
  "distro-1": {
    "arch-1": {
	  "image-1-1-1": [{"name": "fedora"}, {"name": "updates"}],
	  "image-1-1-2": [{"name": "fedora"}, {"name": "updates"}]
	},
    "arch-2": {
	  "image-1-2-1": [{"name": "fedora"}, {"name": "updates"}],
	  "image-1-2-2": [{"name": "fedora"}, {"name": "updates"}]
	}
  },
  "distro-2": {
    "arch-1": {
	  "image-2-1-1": [{"name": "fedora"}, {"name": "updates"}],
	  "image-2-1-2": [{"name": "fedora"}, {"name": "updates"}]
	},
    "arch-2": {
	  "image-2-2-1": [{"name": "fedora"}, {"name": "updates"}],
	  "image-2-2-2": [{"name": "fedora"}, {"name": "updates"}]
	}
  }
}`

	mc := MockClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(json))),
			}, nil
		},
	}
	tc := NewClient(context.Background(), &mc, "")

	distros, err := tc.ListDistros()
	require.Nil(t, err)
	require.Greater(t, len(distros), 0)
	assert.Equal(t, distros, []string{"distro-1", "distro-2"})
}
