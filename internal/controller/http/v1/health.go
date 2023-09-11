package v1

import "net/http"

var _ http.Handler = (*healthHandler)(nil)

type healthHandler struct{}

func (h healthHandler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)
}
