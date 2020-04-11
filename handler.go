package graphql_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/graph-gophers/graphql-go"
	"github.com/satori/go.uuid"
	"net/http"
)

type Logger interface {
	Error(err error, requestID uuid.UUID)
}

func NewHandler(schema *graphql.Schema, logger Logger) (*Handler, error) {

	if schema == nil {
		return nil, errors.New("missing schema")
	}

	if logger == nil {
		return nil, errors.New("missing logger")
	}

	return &Handler{
		schema: schema,
		logger: logger,
	}, nil

}

type Handler struct {
	schema *graphql.Schema
	logger Logger
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	requestID := uuid.NewV4()

	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := h.schema.Exec(r.Context(), params.Query, params.OperationName, params.Variables)

	if response.Errors != nil {
		for _, err := range response.Errors {
			h.logger.Error(err, requestID)
			err.Message = fmt.Sprintf("error in request: %s", requestID.String())
			err.ResolverError = nil
		}
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(responseJSON); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
