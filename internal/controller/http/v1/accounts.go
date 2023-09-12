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

type accountsRoutes struct {
	a usecase.Account
	l logger.Interface
}

func registerAccountsRoutes(router *mux.Router, a usecase.Account, l logger.Interface) {
	r := &accountsRoutes{
		a: a,
		l: l,
	}

	accountsRouter := router.PathPrefix("/accounts/").Subrouter()

	withID := accountsRouter.PathPrefix("/{id:[1-9][0-9]*}").Subrouter()

	withID.NewRoute().Methods(http.MethodPost).HandlerFunc(r.PostID)
	withID.NewRoute().Methods(http.MethodDelete).HandlerFunc(r.DeleteID)
	withID.NewRoute().Methods(http.MethodGet).Path("/airlines").HandlerFunc(r.GetIDAirlines)
	withID.NewRoute().Methods(http.MethodPut).Path("/schema").HandlerFunc(r.PutIDSchema)
}

func (a *accountsRoutes) extractID(r *http.Request) int {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		a.l.Fatal("unexpected error while parsing id: %w", err)
	}

	return int(id)
}

func (a *accountsRoutes) PostID(w http.ResponseWriter, r *http.Request) {
	id := a.extractID(r)

	err := a.a.Add(context.Background(), entity.Account{
		ID: id,
	})

	if err != nil {
		writeError(w, err)
		return
	}

	writeOK(w)
}

func (a *accountsRoutes) DeleteID(w http.ResponseWriter, r *http.Request) {
	id := a.extractID(r)

	err := a.a.Delete(context.Background(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeOK(w)
}

func (a *accountsRoutes) GetIDAirlines(w http.ResponseWriter, r *http.Request) {
	id := a.extractID(r)

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

func (a *accountsRoutes) PutIDSchema(w http.ResponseWriter, r *http.Request) {
	id := a.extractID(r)

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

	writeOK(w)
}
