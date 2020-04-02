package ulid

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUlidPackage(t *testing.T) {

	t.Run("Valid ULID", func(t *testing.T) {
		response := New()

		require.NotNil(t, response)
		require.True(t, len(response) > 0)
		assert.IsType(t, reflect.TypeOf("").String(), response)
	})

	t.Run("Valid Monotonic time", func(t *testing.T) {
		response := NewTime()

		require.NotNil(t, response)
		require.True(t, len(response) > 0)
		assert.Equal(t, reflect.String, reflect.TypeOf(response).Kind())
		assert.IsType(t, reflect.TypeOf("").String(), response)
	})
}
