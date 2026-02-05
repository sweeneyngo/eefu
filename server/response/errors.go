package response

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

type HTTPHandlerError struct {
	HTTPStatus         int
	HTTPMessage        string
	InternalLogMessage string
	Error              error
}

func BadRequest(message string) *HTTPHandlerError {
	return &HTTPHandlerError{
		HTTPStatus:         http.StatusBadRequest,
		HTTPMessage:        message,
		InternalLogMessage: message,
	}
}

func NotFound(entity string) *HTTPHandlerError {
	return &HTTPHandlerError{
		HTTPStatus:         http.StatusNotFound,
		HTTPMessage:        entity + " not found",
		InternalLogMessage: entity + " missing in DB",
	}
}

func InternalServerError(err error, message string) *HTTPHandlerError {
	return &HTTPHandlerError{
		HTTPStatus:         http.StatusInternalServerError,
		HTTPMessage:        message,
		InternalLogMessage: message,
		Error:              err,
	}
}

func (handlerError *HTTPHandlerError) Respond(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())
	evt := logger.Error().Int("status", handlerError.HTTPStatus)
	if handlerError.Error != nil {
		evt = evt.Err(handlerError.Error)
	}
	evt.Msg(handlerError.InternalLogMessage)
	http.Error(w, handlerError.HTTPMessage, handlerError.HTTPStatus)
}
