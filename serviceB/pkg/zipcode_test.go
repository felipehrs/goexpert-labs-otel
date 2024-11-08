package pkg_test

import (
	"testing"

	. "github.com/felipehrs/goexpert-labs-otel-serciceB/pkg"
	"github.com/stretchr/testify/assert"
)

func TestIsValidZipCode(t *testing.T) {
	assert.True(t, IsValidZipCode("01001-000"))
	assert.True(t, IsValidZipCode("01001000"))
	assert.False(t, IsValidZipCode("01001-0000"))
	assert.False(t, IsValidZipCode("0100-100"))
	assert.False(t, IsValidZipCode("abcde-fgh"))
}
