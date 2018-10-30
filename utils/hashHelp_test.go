package utils

import (
	"testing"
	"github.com/sirupsen/logrus"
)

func TestHashCode(t *testing.T) {
	logrus.Println(HashCode([]byte("aaaaa")))
	logrus.Println(HashCode([]byte("dDSDSF")))
	logrus.Println(HashCode([]byte("890@323")))
}
