package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

func TestCreateSendMessageIinput(t *testing.T) {

	t.Run("Create messages from test files", func(t *testing.T) {
		filenames := []string{
			"../../test-data/sales-messages/1808710-body.json",
			"../../test-data/sales-messages/1808711-body.json",
			"../../test-data/sales-messages/1808712-body.json",
			"../../test-data/sales-messages/1808713-body.json",
			"../../test-data/sales-messages/1808714-body.json",
		}

		for seqNo, filename := range filenames {
			jsonFile, _ := ioutil.ReadFile(filename)

			var data map[string]interface{}

			jsonErr := json.Unmarshal([]byte(jsonFile), &data)
			if jsonErr != nil {
				log.WithError(jsonErr).Errorf("failed to unmarshal file to map: %s", filename)
				return
			}
			smi, err := createSendMessageInput(data, "888", "999", "aaa", seqNo)
			assert.Nil(t, err)
			assert.NotNil(t, smi)
		}
	})
}
