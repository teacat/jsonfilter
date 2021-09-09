package parser

import (
	"testing"

	"github.com/alecthomas/repr"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	a := assert.New(t)
	r, err := Parse("a(b(c/*/e))")
	a.NoError(err)
	repr.Println(r, repr.Indent("  "), repr.OmitEmpty(true))
	panic("hello")
}
