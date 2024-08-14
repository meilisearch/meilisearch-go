package meilisearch

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRawType_UnmarshalJSON(t *testing.T) {
	var r RawType

	data := []byte(`"example"`)
	err := r.UnmarshalJSON(data)
	assert.NoError(t, err)
	assert.Equal(t, RawType(`"example"`), r)

	data = []byte(`""`)
	err = r.UnmarshalJSON(data)
	assert.NoError(t, err)
	assert.Equal(t, RawType(`""`), r)

	data = []byte(`{invalid}`)
	err = r.UnmarshalJSON(data)
	assert.NoError(t, err)
	assert.Equal(t, RawType(`{invalid}`), r)
}

func TestRawType_MarshalJSON(t *testing.T) {
	r := RawType(`"example"`)
	data, err := r.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"example"`), data)

	r = RawType(`""`)
	data, err = r.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte(`""`), data)

	r = RawType(`{random}`)
	data, err = r.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{random}`), data)
}
