package ulid

import (
	"fmt"
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
		fmt.Println(reflect.TypeOf(response))
		assert.IsType(t, reflect.TypeOf("").String(), response)
	})
}
