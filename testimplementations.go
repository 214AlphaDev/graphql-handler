package graphql_handler

import "github.com/satori/go.uuid"

type testLogger struct {
	error func(err error, requestID uuid.UUID)
}

func (l *testLogger) Error(err error, requestID uuid.UUID) {
	l.error(err, requestID)
}
