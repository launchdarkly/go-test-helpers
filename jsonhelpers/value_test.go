package jsonhelpers

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJValueOf(t *testing.T) {
	s := `{"a":true}`
	m := map[string]any{"a": true}

	v1 := JValueOf([]byte(s))
	assert.Nil(t, v1.Error())
	assert.Equal(t, m, v1.parsed)
	assert.Equal(t, s, v1.String())

	v2 := JValueOf(json.RawMessage(s))
	assert.Nil(t, v2.Error())
	assert.Equal(t, m, v2.parsed)
	assert.Equal(t, s, v2.String())
	assert.Equal(t, v1, v2)

	v3 := JValueOf(s)
	assert.Nil(t, v3.Error())
	assert.Equal(t, m, v3.parsed)
	assert.Equal(t, s, v3.String())
	assert.Equal(t, v1, v3)

	v4 := JValueOf(m)
	assert.Nil(t, v4.Error())
	assert.Equal(t, m, v4.parsed)
	assert.Equal(t, s, v4.String())
	assert.Equal(t, v1, v4)

	v5 := JValueOf(v4)
	assert.Equal(t, v4, v5)

	v6 := JValueOf("{no")
	assert.NotNil(t, v6.Error())
	assert.Equal(t, "{no", v6.String())
}
