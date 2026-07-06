package response

import "net/http"

type Exception struct {
	StatusCode int
	Message    string
}

func (e *Exception) Error() string {
	return e.Message
}

func NewBadRequestException(message ...string) *Exception {
	statusCode := http.StatusBadRequest
	messageDefault := http.StatusText(statusCode)
	if len(message) > 0 && message[0] != ""{
		messageDefault = message[0]
	}
	return &Exception{
		StatusCode: statusCode,
		Message: messageDefault,
	}
}

func NewInternalServerErrorException(message ...string) *Exception {
	statusCode := http.StatusInternalServerError
	messageDefault := http.StatusText(statusCode)
	
	if len(message) > 0 && message[0] != ""{
		messageDefault = message[0]
	}

	return &Exception{
		StatusCode: statusCode,
		Message: messageDefault,
	}
}

func NewUnauthorizedException(message ...string) *Exception {
	statusCode := http.StatusUnauthorized
	messageDefault := http.StatusText(statusCode)
	if len(message) > 0 && message[0] != ""{
		messageDefault = message[0]
	}
	return &Exception{
		StatusCode: statusCode,
		Message: messageDefault,
	}
}

func NewForbiddenException(message ...string) *Exception {
	statusCode := http.StatusForbidden
	messageDefault := http.StatusText(statusCode)
	if len(message) > 0 && message[0] != ""{
		messageDefault = message[0]
	}
	return &Exception{
		StatusCode: statusCode,
		Message: messageDefault,
	}
}
