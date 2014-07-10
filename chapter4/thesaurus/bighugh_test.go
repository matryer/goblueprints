package thesaurus

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterface(t *testing.T) {
	var _ Thesaurus = &BigHugh{}
}

func GetAPIKey(t *testing.T) string {
	apiKey := os.Getenv("BHT_APIKEY")
	assert.NotEmpty(t, apiKey, "Tests need an API key to run")
	return apiKey
}

func TestSynonyms(t *testing.T) {
	var thesaurus Thesaurus = &BigHugh{APIKey: GetAPIKey(t)}
	if syns, err := thesaurus.Synonyms("love"); assert.NoError(t, err) {
		assert.True(t, len(syns) > 0, "Should have at least one synonym")
	}
}
