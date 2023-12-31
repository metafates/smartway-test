package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/metafates/smartway-test/internal/usecase"
	"github.com/metafates/smartway-test/pkg/logger"
)

func RegisterRoutes(router *mux.Router, useCases usecase.UseCases, l logger.Interface) {
	v1 := router.PathPrefix("/v1/").Subrouter()
	{
		v1.
			NewRoute().
			Methods(http.MethodGet).
			Path("/health").
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

		registerAccountsRoutes(v1, useCases.Account, l)
		registerProvidersRoutes(v1, useCases.Provider, l)
		registerSchemasRoutes(v1, useCases.Schema, l)
		registerAirlinesRoutes(v1, useCases.Airline, l)
	}
}
