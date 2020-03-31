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

func TestParseCommaSeparatedFiles(t *testing.T) {
	filenames := `/Users/sg/Dropbox (Personal)/API/s2c/aws2sap-dlq/18082210000.json,` +
		`/Users/sg/Dropbox (Personal)/API/s2c/aws2sap-dlq/15827625345.json`

	t.Run("test for correct split", func(t *testing.T) {
		filenames := ParseCommaSeparatedFiles(filenames)
		assert.NotNil(t, filenames)
		assert.NotEmpty(t, filenames)
		assert.Len(t, filenames, 2)
	})
}

func TestParseCommaSeparatedStrings(t *testing.T) {
	t.Run("test separation", func(t *testing.T) {
		str := "1, 2, 3,   4.1, 5"
		strs := ParseCommaSeparatedStrings(str)
		assert.NotNil(t, strs)
		assert.NotEmpty(t, strs)
		assert.Len(t, strs, 5)
	})

	t.Run("test empty element", func(t *testing.T) {
		str := "1, 2, 3,, 5"
		strs := ParseCommaSeparatedStrings(str)
		assert.NotNil(t, strs)
		assert.NotEmpty(t, strs)
		assert.Len(t, strs, 4)
	})
}

func TestSelectRandomString(t *testing.T) {
	t.Run("select from slice", func(t *testing.T) {
		strs := make([]string, 3)
		strs[0] = "one"
		strs[1] = "two"
		strs[2] = "tre"

		str := SelectRandomString(strs)
		assert.NotNil(t, str)
		assert.NotEmpty(t, str)
		assert.Contains(t, strs, str)
	})

	t.Run("select from empty slice", func(t *testing.T) {
		strs := []string{}
		str := SelectRandomString(strs)
		assert.Empty(t, str)
		assert.Equal(t, "", str)
	})
}
