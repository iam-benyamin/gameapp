package httpmsg

import (
	"gameapp/pkg/errmsg"
	"gameapp/pkg/richerror"
	"net/http"
)

func Error(err error) (message string, code int) {
	// we should not expose unexpected error messages
	switch err.(type) {
	case richerror.RichError:
		re := err.(richerror.RichError)
		msg := re.Message()
		kind := mapKindToStatusCode(re.Kind())
		if kind >= 500 {
			msg = errmsg.ErrorSomeThingWentWrong
		}
		return msg, kind
	default:
		return err.Error(), http.StatusBadRequest
	}
}

func mapKindToStatusCode(kind richerror.Kind) int {
	switch kind {
	case richerror.KindInvalid:
		return http.StatusUnprocessableEntity
	case richerror.KindNotFound:
		return http.StatusNotFound
	case richerror.KindForbidden:
		return http.StatusForbidden
	case richerror.KindUnexpected:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
