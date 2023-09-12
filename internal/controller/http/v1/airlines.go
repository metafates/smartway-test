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

type airlinesRoutes struct {
	l logger.Interface
	a usecase.Airline
}

func registerAirlinesRoutes(router *mux.Route, a usecase.Airline, l logger.Interface) {
	r := &airlinesRoutes{
		l: l,
		a: a,
	}

	accountsRouter := router.PathPrefix("/airlines/").Subrouter()

	withCode := accountsRouter.PathPrefix("/{code}").Subrouter()

	withCode.NewRoute().Methods(http.MethodPost).HandlerFunc(r.PostID)
}

func (a *airlinesRoutes) extractCode(r *http.Request) (entity.AirlineCode, error) {
	vars := mux.Vars(r)

	var code entity.AirlineCode
	if err := json.Unmarshal([]byte(vars["code"]), &code); err != nil {
		return "", err
	}

	return code, nil
}

type postAirlinesIDRequest struct {
	Name string `json:"name,omitempty"`
}

func (a *airlinesRoutes) PostID(w http.ResponseWriter, r *http.Request) {
	code, err := a.extractCode(r)
	if err != nil {
		writeError(w, err)
		return
	}

	var request postAirlinesIDRequest
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
