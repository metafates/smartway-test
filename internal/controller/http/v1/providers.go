package v1

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/metafates/smartway-test/internal/entity"
	"github.com/metafates/smartway-test/internal/usecase"
	"github.com/metafates/smartway-test/pkg/logger"
)

type providersRoutes struct {
	l logger.Interface
	p usecase.Provider
}

func registerProvidersRoutes(router *mux.Router, p usecase.Provider, l logger.Interface) {
	r := &providersRoutes{
		l: l,
		p: p,
	}

	providersRouter := router.PathPrefix("/providers/").Subrouter()
	{
		withID := providersRouter.PathPrefix("/{id}").Subrouter()
		{
			withID.NewRoute().Methods(http.MethodPost).HandlerFunc(r.PostID)
			withID.NewRoute().Methods(http.MethodDelete).HandlerFunc(r.DeleteID)
			withID.NewRoute().Methods(http.MethodGet).Path("/airlines").HandlerFunc(r.GetIDAirlines)
		}
	}
}

func (p *providersRoutes) extractID(r *http.Request) (entity.ProviderID, error) {
	vars := mux.Vars(r)

	var id entity.ProviderID

	err := id.UnmarshalText([]byte(vars["id"]))
	if err != nil {
		return "", err
	}

	return id, nil
}

func (p *providersRoutes) PostID(w http.ResponseWriter, r *http.Request) {
	id, err := p.extractID(r)
	if err != nil {
		writeError(w, err)
		return
	}

	var request struct {
		Name string `json:"name,omitempty"`
	}
	if err := bindJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}

	err = p.p.Add(context.Background(), entity.Provider{
		ID:   id,
		Name: request.Name,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (p *providersRoutes) DeleteID(w http.ResponseWriter, r *http.Request) {
	id, err := p.extractID(r)
	if err != nil {
		writeError(w, err)
		return
	}

	if err := p.p.Delete(context.Background(), id); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *providersRoutes) GetIDAirlines(w http.ResponseWriter, r *http.Request) {
	id, err := p.extractID(r)
	if err != nil {
		writeError(w, err)
		return
	}

	airlines, err := p.p.GetAirlines(context.Background(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, airlines, http.StatusOK)
}
