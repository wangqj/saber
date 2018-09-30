package utils

import (
	"testing"
	"github.com/sirupsen/logrus"
)

func TestHashCode(t *testing.T) {
	logrus.Println(HashCode("aaaaa"))
	logrus.Println(HashCode("dDSDSF"))
	logrus.Println(HashCode("890@323"))
}
