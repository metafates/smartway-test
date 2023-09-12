package v1

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/metafates/smartway-test/internal/entity"
	"github.com/metafates/smartway-test/internal/usecase"
	"github.com/metafates/smartway-test/pkg/logger"
)

type schemasRoutes struct {
	l logger.Interface
	s usecase.Schema
}

func registerSchemasRoutes(router *mux.Router, s usecase.Schema, l logger.Interface) {
	r := schemasRoutes{
		l: l,
		s: s,
	}

	schemasRouter := router.PathPrefix("/schemas/").Subrouter()

	schemasRouter.NewRoute().Methods(http.MethodGet).Path("/find").Queries("name", "{name}").HandlerFunc(r.GetSearch)

	withID := schemasRouter.PathPrefix("/{id}").Subrouter()

	withID.NewRoute().Methods(http.MethodPost).HandlerFunc(r.PostID)
	withID.NewRoute().Methods(http.MethodDelete).HandlerFunc(r.DeleteID)
	withID.NewRoute().Methods(http.MethodPut).HandlerFunc(r.PutID)
}

func (s *schemasRoutes) extractID(r *http.Request) (entity.SchemaID, error) {
	vars := mux.Vars(r)

	var id entity.SchemaID

	err := id.UnmarshalText([]byte(vars["id"]))
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *schemasRoutes) PostID(w http.ResponseWriter, r *http.Request) {
	id, err := s.extractID(r)
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

	err = s.s.Add(context.Background(), entity.Schema{
		Name: request.Name,
		ID:   id,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *schemasRoutes) DeleteID(w http.ResponseWriter, r *http.Request) {
	id, err := s.extractID(r)
	if err != nil {
		writeError(w, err)
		return
	}

	err = s.s.Delete(context.Background(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *schemasRoutes) PutID(w http.ResponseWriter, r *http.Request) {
	id, err := s.extractID(r)
	if err != nil {
		writeError(w, err)
		return
	}

	var request entity.SchemaChanges
	if err := bindJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}

	if err := s.s.Update(context.Background(), id, request); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *schemasRoutes) GetSearch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	schema, found, err := s.s.Find(context.Background(), name)
	if err != nil {
		writeError(w, err)
		return
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	writeJSON(w, schema, http.StatusOK)
}
