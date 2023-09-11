package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/metafates/smartway-test/internal/usecase"
	"github.com/metafates/smartway-test/pkg/logger"
)

type UseCases struct {
	Account usecase.Account
}

func RegisterRoutes(router *mux.Router, useCases UseCases, l logger.Interface) {
	v1 := router.PathPrefix("/v1/").Subrouter()

	v1.
		NewRoute().
		Methods(http.MethodGet).
		Path("/health").
		Handler(healthHandler{})

	registerAccountRoutes(v1, useCases.Account, l)
}
