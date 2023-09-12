package v1

import (
	"context"
	"net/http"
	"strconv"

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

	withID := providersRouter.PathPrefix("/{id:[1-9][0-9]*}").Subrouter()

	withID.NewRoute().Methods(http.MethodPost).HandlerFunc(r.PostID)
	withID.NewRoute().Methods(http.MethodDelete).HandlerFunc(r.DeleteID)
	withID.NewRoute().Methods(http.MethodGet).Path("/airlines").HandlerFunc(r.GetIDAirlines)
}

func (p *providersRoutes) extractID(r *http.Request) int {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		p.l.Fatal("unexpected error while parsing id: %w", err)
	}

	return int(id)
}

type postProvidersIDRequest struct {
	Name string `json:"name,omitempty"`
}

func (p *providersRoutes) PostID(w http.ResponseWriter, r *http.Request) {
	id := p.extractID(r)

	var request postProvidersIDRequest
	if err := bindJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}

	err := p.p.Add(context.Background(), entity.Provider{
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
	id := p.extractID(r)

	if err := p.p.Delete(context.Background(), id); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *providersRoutes) GetIDAirlines(w http.ResponseWriter, r *http.Request) {
	id := p.extractID(r)

	airlines, err := p.p.GetAirlines(context.Background(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, airlines, http.StatusOK)
}
