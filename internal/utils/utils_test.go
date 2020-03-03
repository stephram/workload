package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestGetLogger(t *testing.T) {
	t.Run("check that logger is initialised", func(t *testing.T) {
		_log := GetLogger()
		assert.NotNil(t, _log)
	})
}
