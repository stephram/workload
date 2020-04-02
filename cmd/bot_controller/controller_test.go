package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStoreNumbers(t *testing.T) {

	t.Run("Parse store numbers string", func(t *testing.T) {
		storeIDs, err := parseStoreNumbers("A001, A002, A003")
		assert.Nil(t, err)
		assert.True(t, len(storeIDs) > 1)
		assert.Equal(t, 3, len(storeIDs))
	})
}
