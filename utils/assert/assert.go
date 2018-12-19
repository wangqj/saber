// Copyright 2016 CodisLabs. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

package assert

import "github.com/sirupsen/logrus"

func Must(b bool) {
	if b {
		return
	}
	logrus.Panic("assertion failed")
}

func MustNoError(err error) {
	if err == nil {
		return
	}
	logrus.Panic(err, "error happens, assertion failed")
}
