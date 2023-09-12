package v1

import (
	"context"
	"net/http"

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

	withID := accountsRouter.PathPrefix("/{id}").Subrouter()

	withID.NewRoute().Methods(http.MethodPost).HandlerFunc(r.PostID)
	withID.NewRoute().Methods(http.MethodDelete).HandlerFunc(r.DeleteID)
	withID.NewRoute().Methods(http.MethodGet).Path("/airlines").HandlerFunc(r.GetIDAirlines)
	withID.NewRoute().Methods(http.MethodPut).Path("/schema").HandlerFunc(r.PutIDSchema)
}

func (a *accountsRoutes) extractID(r *http.Request) (entity.AccountID, error) {
	vars := mux.Vars(r)

	var id entity.AccountID

	err := id.UnmarshalText([]byte(vars["id"]))
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (a *accountsRoutes) PostID(w http.ResponseWriter, r *http.Request) {
	id, err := a.extractID(r)
	if err != nil {
		writeError(w, err)
		return
	}

	err = a.a.Add(context.Background(), entity.Account{
		ID: id,
	})

	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *accountsRoutes) DeleteID(w http.ResponseWriter, r *http.Request) {
	id, err := a.extractID(r)
	if err != nil {
		writeError(w, err)
		return
	}

	err = a.a.Delete(context.Background(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *accountsRoutes) GetIDAirlines(w http.ResponseWriter, r *http.Request) {
	id, err := a.extractID(r)
	if err != nil {
		writeError(w, err)
		return
	}

	airlines, err := a.a.GetAirlines(context.Background(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, airlines, http.StatusOK)
}

func (a *accountsRoutes) PutIDSchema(w http.ResponseWriter, r *http.Request) {
	id, err := a.extractID(r)
	if err != nil {
		writeError(w, err)
		return
	}

	var request struct {
		ID entity.SchemaID `json:"id"`
	}

	err = bindJSON(r, &request)
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
