package dictionary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsContained(t *testing.T) {
	d := Dictionary_CreateDictionary()

	assert.True(t, d.Dictionary_IsContained([]byte{224}), "Fail 1")
	assert.True(t, d.Dictionary_IsContained([]byte{7}), "Fail 2")
	assert.False(t, d.Dictionary_IsContained([]byte{56, 78}), "Fail 1")

}

func TestAddKey(t *testing.T) {
	d := Dictionary_CreateDictionary()

	d.Dictionary_AddKey([]byte{12, 12})
	d.Dictionary_AddKey([]byte{43, 65, 124})

	assert.True(t, d.Dictionary_IsContained([]byte{12, 12}), "Fail 1")
	assert.True(t, d.Dictionary_IsContained([]byte{43, 65, 124}), "Fail 2")

}

func TestGetVal(t *testing.T) {
	d := Dictionary_CreateDictionary()

	d.Dictionary_AddKey([]byte{12, 12})
	d.Dictionary_AddKey([]byte{43, 65, 124})

	assert.Equal(t, 256, d.Dictionary_GetVal([]byte{12, 12}), "Fail 1")
	assert.Equal(t, 257, d.Dictionary_GetVal([]byte{43, 65, 124}), "Fail 1")

}
