package jsonhelpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanonicalizeJSON(t *testing.T) {
	assert.Equal(t, `null`, string(CanonicalizeJSON([]byte(`null`))))
	assert.Equal(t, `true`, string(CanonicalizeJSON([]byte(`true`))))
	assert.Equal(t, `123`, string(CanonicalizeJSON([]byte(`123`))))
	assert.Equal(t, `"x"`, string(CanonicalizeJSON([]byte(`"x"`))))

	assert.Equal(t, `{"a":1,"b":2,"c":3}`,
		string(CanonicalizeJSON([]byte(`{"a":1,"b":2,"c":3}`))))

	assert.Equal(t, `{"a":1,"b":2,"c":3}`,
		string(CanonicalizeJSON([]byte(`{"c":3,"b":2,"a":1}`))))

	assert.Equal(t, `{"a":1,"b":{"A":10,"B":20,"C":30},"c":3}`,
		string(CanonicalizeJSON([]byte(`{"c":3,"b":{"B":20,"A":10,"C":30},"a":1}`))))

	assert.Equal(t, `[{"a":1,"b":2,"c":3},{"A":10,"B":20,"C":30}]`,
		string(CanonicalizeJSON([]byte(`[{"c":3,"b":2,"a":1},{"B":20,"A":10,"C":30}]`))))
}
