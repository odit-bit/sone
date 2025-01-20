package database

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_columns(t *testing.T) {
	c := columns{"a", "b", "c"}
	actual := c.query()
	expected := "a,b,c"

	assert.Equal(t, expected, actual)

	expected = "?,?,?"
	actual = c.placeholders()
	assert.Equal(t, expected, actual)

	query := fmt.Sprintf("SELECT a, b, c FROM tables WHERE %s = %s LIMIT 1", c.query(), c.placeholders())
	expected = "SELECT a, b, c FROM tables WHERE a,b,c = ?,?,? LIMIT 1"
	assert.Equal(t, expected, query)

}
