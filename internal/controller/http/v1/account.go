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

type accountRoutes struct {
	a usecase.Account
	l logger.Interface
}

func registerAccountRoutes(router *mux.Router, a usecase.Account, l logger.Interface) {
	r := &accountRoutes{a, l}

	accountsRouter := router.PathPrefix("/accounts/").Subrouter()

	accountRouter := accountsRouter.PathPrefix("/{id:[1-9][0-9]*}").Subrouter()

	accountRouter.NewRoute().Methods(http.MethodPost).HandlerFunc(r.add)
	accountRouter.NewRoute().Methods(http.MethodDelete).HandlerFunc(r.delete)
	accountRouter.NewRoute().Methods(http.MethodGet).Path("/airlines").HandlerFunc(r.getAirlines)
	accountRouter.NewRoute().Methods(http.MethodPut).Path("/schema").HandlerFunc(r.setSchema)
}

func (a *accountRoutes) getID(r *http.Request) int {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		a.l.Fatal("unexpected error while parsing id: %w", err)
	}

	return int(id)
}

func (a *accountRoutes) add(w http.ResponseWriter, r *http.Request) {
	id := a.getID(r)

	err := a.a.Add(context.Background(), entity.Account{
		ID: id,
	})

	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *accountRoutes) delete(w http.ResponseWriter, r *http.Request) {
	id := a.getID(r)

	err := a.a.Delete(context.Background(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *accountRoutes) getAirlines(w http.ResponseWriter, r *http.Request) {
	id := a.getID(r)

	airlines, err := a.a.GetAirlines(context.Background(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, airlines, http.StatusOK)
}

type setSchemaRequest struct {
	ID int `json:"id"`
}

func (a *accountRoutes) setSchema(w http.ResponseWriter, r *http.Request) {
	id := a.getID(r)

	var request setSchemaRequest
	err := bindJSON(r, &request)
	if err != nil {
		writeError(w, err)
		return
	}

	err = a.a.SetSchema(context.Background(), id, request.ID)
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
