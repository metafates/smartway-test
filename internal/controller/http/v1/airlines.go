package v1

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/metafates/smartway-test/internal/entity"
	"github.com/metafates/smartway-test/internal/pkg/hashset"
	"github.com/metafates/smartway-test/internal/usecase"
	"github.com/metafates/smartway-test/pkg/logger"
)

type airlinesRoutes struct {
	l logger.Interface
	a usecase.Airline
}

func registerAirlinesRoutes(router *mux.Router, a usecase.Airline, l logger.Interface) {
	r := &airlinesRoutes{
		l: l,
		a: a,
	}

	airlinesRouter := router.PathPrefix("/airlines/").Subrouter()
	{
		withCode := airlinesRouter.PathPrefix("/{code}").Subrouter()
		{
			withCode.NewRoute().Methods(http.MethodPost).HandlerFunc(r.PostID)
			withCode.NewRoute().Methods(http.MethodDelete).HandlerFunc(r.DeleteID)
			withCode.NewRoute().Methods(http.MethodPut).Path("/providers").HandlerFunc(r.PutIDProviders)
		}
	}
}

func (a *airlinesRoutes) extractCode(r *http.Request) (entity.AirlineCode, error) {
	vars := mux.Vars(r)

	var code entity.AirlineCode

	err := code.UnmarshalText([]byte(vars["code"]))
	if err != nil {
		return "", err
	}

	return code, nil
}

func (a *airlinesRoutes) PostID(w http.ResponseWriter, r *http.Request) {
	code, err := a.extractCode(r)
	if err != nil {
		writeError(w, err)
		return
	}

	var request struct {
		Name string `json:"name"`
	}
	if err := bindJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}

	err = a.a.Add(context.Background(), entity.Airline{
		Code: code,
		Name: request.Name,
	})

	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *airlinesRoutes) DeleteID(w http.ResponseWriter, r *http.Request) {
	code, err := a.extractCode(r)
	if err != nil {
		writeError(w, err)
		return
	}

	if err := a.a.Delete(context.Background(), code); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *airlinesRoutes) PutIDProviders(w http.ResponseWriter, r *http.Request) {
	code, err := a.extractCode(r)
	if err != nil {
		writeError(w, err)
		return
	}

	var request struct {
		Providers *hashset.Set[entity.ProviderID] `json:"providers"`
	}

	if err := bindJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}

	err = a.a.SetProviders(context.Background(), code, request.Providers.Values())
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
