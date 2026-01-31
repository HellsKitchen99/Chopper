package main

import (
	"chopper/internal/build"

	"github.com/sirupsen/logrus"
)

func main() {
	if err := build.Run(); err != nil {
		logrus.Error(err)
		return
	}
}
