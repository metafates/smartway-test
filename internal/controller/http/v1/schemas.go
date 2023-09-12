package v1

import (
	"context"
	"encoding/json"
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

	accountsRouter := router.PathPrefix("/schemas/").Subrouter()

	withID := accountsRouter.PathPrefix("/{id}").Subrouter()

	withID.NewRoute().Methods(http.MethodPost).HandlerFunc(r.PostID)
	withID.NewRoute().Methods(http.MethodPut).HandlerFunc(r.PutID)
}

func (s *schemasRoutes) extractID(r *http.Request) (entity.SchemaID, error) {
	vars := mux.Vars(r)

	var id entity.SchemaID
	if err := json.Unmarshal([]byte(vars["id"]), &id); err != nil {
		return 0, err
	}

	return id, nil
}

type postSchemasIDRequest struct {
	Name string `json:"name"`
}

func (s *schemasRoutes) PostID(w http.ResponseWriter, r *http.Request) {
	id, err := s.extractID(r)
	if err != nil {
		writeError(w, err)
		return
	}

	var request postSchemasIDRequest
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

func (s *schemasRoutes) PutID(w http.ResponseWriter, r *http.Request) {
	id, err := s.extractID(r)
	if err != nil {
		writeError(w, err)
		return
	}

	var request entity.Schema
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
