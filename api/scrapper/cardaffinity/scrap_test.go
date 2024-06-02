package cardaffinity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Scrap(t *testing.T) {
	s := NewScrapper()
	result, err := s.Scrap("Abrade")
	require.NoError(t, err)
	require.True(t, len(result) > 0)
}
