package werk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	pool := NewPool(10).Start()

	assert.Equal(t, pool.Available(), 10)

	pool.Do("hello", func(v Work) {
		assert.Equal(t, "hello", v.(string))
		assert.Equal(t, 9, pool.Available())
	})

	assert.Equal(t, pool.Available(), 10)
}
