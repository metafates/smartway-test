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

	withID := accountsRouter.PathPrefix("/{id:[1-9][0-9]*}").Subrouter()

	withID.NewRoute().Methods(http.MethodPost).HandlerFunc(r.PostID)
	withID.NewRoute().Methods(http.MethodPut).HandlerFunc(r.PutID)
}

func (s *schemasRoutes) extractID(r *http.Request) int {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		s.l.Fatal("unexpected error while parsing id: %w", err)
	}

	return int(id)
}

type postSchemasIDRequest struct {
	Name string `json:"name"`
}

func (s *schemasRoutes) PostID(w http.ResponseWriter, r *http.Request) {
	id := s.extractID(r)

	var request postSchemasIDRequest
	if err := bindJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}

	err := s.s.Add(context.Background(), entity.Schema{
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
	id := s.extractID(r)

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
