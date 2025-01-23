package test

import (
	"io"

	"github.com/nut-game/nano/logger/interfaces"
	lwrapper "github.com/nut-game/nano/logger/logrus"
	tests "github.com/sirupsen/logrus/hooks/test"
)

// NewNullLogger creates a discarding logger and installs the test hook.
func NewNullLogger() (interfaces.Logger, *tests.Hook) {
	logger, hook := tests.NewNullLogger()
	logger.Out = io.Discard
	return lwrapper.NewWithFieldLogger(logger), hook
}
