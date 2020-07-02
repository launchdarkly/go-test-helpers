package helpers

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadWithTimeout(t *testing.T) {
	var buf1 bytes.Buffer
	data := ReadWithTimeout(&buf1, 5, time.Millisecond*50)
	assert.Len(t, data, 0)

	var buf2 bytes.Buffer
	buf2.WriteString("he")
	data = ReadWithTimeout(&buf2, 5, time.Millisecond*50)
	assert.Equal(t, "he", string(data))

	var buf3 bytes.Buffer
	buf3.WriteString("hello")
	data = ReadWithTimeout(&buf3, 5, time.Millisecond*50)
	assert.Equal(t, "hello", string(data))

	r1, w1 := io.Pipe()
	go func() {
		w1.Write([]byte("good"))
		time.Sleep(10 * time.Millisecond)
		w1.Write([]byte("bye"))
	}()
	data = ReadWithTimeout(r1, 7, time.Millisecond*100)
	assert.Equal(t, "goodbye", string(data))

	r2, w2 := io.Pipe()
	go func() {
		w2.Write([]byte("good"))
	}()
	data = ReadWithTimeout(r2, 7, time.Millisecond*100)
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, "good", string(data))
}
