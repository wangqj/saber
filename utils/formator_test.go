package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAddZeroForInt(t *testing.T) {
	assert.Equal(t, AddZeroForInt(10, 4), "0010")
}

func TestAddZeroForStr(t *testing.T) {
	assert.Equal(t, AddZeroForStr("10", 4), "0010")
}
