package scheme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LiveStreams(t *testing.T) {
	c := make(chan *LiveStream, 1)
	cErr := make(chan error)
	c <- &LiveStream{ID: "123"}

	list := LiveStreamIterator(c, cErr)
	ls, ok := list.Next()
	if !ok {
		t.Fatal("should true")
	}
	assert.Equal(t, "123", ls.ID)

	close(cErr)
	_, ok = list.Next()
	if ok {
		t.Fatal("should false")
	}
	assert.Equal(t, nil, list.Err())

}
