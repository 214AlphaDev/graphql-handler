package graphql_handler

import (
	"fmt"
	"github.com/satori/go.uuid"
)

type ConsoleLogger struct{}

func (l ConsoleLogger) Error(err error, requestID uuid.UUID) {
	fmt.Println(fmt.Sprintf("error '%s' happened in request: %s", err.Error(), requestID.String()))
}
