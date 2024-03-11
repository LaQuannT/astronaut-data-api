package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/LaQuannT/astronaut-data-api/internal/model"
	"github.com/gorilla/mux"
)

type astronautHandler struct {
	service model.AstronautUsecase
	log     *slog.Logger
}

func RegisterAstronautHandlers(s model.AstronautUsecase, r *mux.Router, l *slog.Logger) {
	handler := &astronautHandler{
		service: s,
		log:     l,
	}

	sr := r.PathPrefix("/astronauts").Subrouter()
	// TODO - add api key middleware

	sr.HandleFunc("", handler.CreateAstronaut).Methods("POST")
	sr.HandleFunc("", handler.ListAstronauts).Methods("GET")
	sr.HandleFunc("/{astronautID}", handler.GetAstronaut).Methods("GET")
	sr.HandleFunc("/{astronautID}", handler.UpdateAstronaut).Methods("PUT")
	sr.HandleFunc("/{astronautID}", handler.DeleteAstronaut).Methods("DELETE")
}

func (h *astronautHandler) CreateAstronaut(w http.ResponseWriter, r *http.Request) {
	a := new(model.Astronaut)

	key := r.Header.Get(string(ApiKeyHeader))
	ctx := context.WithValue(r.Context(), ApiKeyHeader, key)

	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		writeJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Invalid request body"})
		h.log.Warn("error decoding request body to astronaut", slog.Any("error", err))
		return
	}

	a, err := h.service.Create(ctx, a)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Bad Request"})
		h.log.Warn("error creating new astronaut", slog.Any("error", err))
		return
	}

	writeJSON(w, http.StatusCreated, model.JSONResponse{Astronaut: a})
}

func (h *astronautHandler) ListAstronauts(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get(string(ApiKeyHeader))
	ctx := context.WithValue(r.Context(), ApiKeyHeader, key)

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "invalid request query"})
		h.log.Warn("error parsing url request query", slog.Any("error", err))
		return
	}

	l := params.Get("limit")
	o := params.Get("offset")

	var limit, offset int

	limit, err = strconv.Atoi(l)
	if err != nil {
		limit = 30
	}

	offset, err = strconv.Atoi(o)
	if err != nil {
		offset = 0
	}

	astronauts, err := h.service.List(ctx, limit, offset)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Bad Request"})
		h.log.Warn("error listing astronauts", slog.Any("error", err))
		return
	}

	writeJSON(w, http.StatusOK, model.JSONResponse{Astronauts: astronauts})
}

func (h *astronautHandler) GetAstronaut(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get(string(ApiKeyHeader))
	ctx := context.WithValue(r.Context(), ApiKeyHeader, key)

	astronautID := mux.Vars(r)["astronautID"]
	id, err := strconv.Atoi(astronautID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Invalid Astronaut ID"})
		return
	}

	a, err := h.service.Get(ctx, id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Bad Request"})
		h.log.Warn("error fetching a astronaut", slog.Any("error", err))
		return
	}

	writeJSON(w, http.StatusOK, model.JSONResponse{Astronaut: a})
}

func (h *astronautHandler) UpdateAstronaut(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get(string(ApiKeyHeader))
	ctx := context.WithValue(r.Context(), ApiKeyHeader, key)
	a := new(model.Astronaut)

	astronautID := mux.Vars(r)["astronautID"]
	id, err := strconv.Atoi(astronautID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Invalid Astronaut ID"})
		return
	}

	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		writeJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Invalid request body"})
		h.log.Warn("error decoding request body to astronaut", slog.Any("error", err))
		return
	}

	a.ID = id

	a, err = h.service.Update(ctx, a)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Bad Request"})
		h.log.Warn("error updating a astronaut", slog.Any("error", err))
		return
	}

	writeJSON(w, http.StatusOK, model.JSONResponse{Astronaut: a})
}

func (h *astronautHandler) DeleteAstronaut(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get(string(ApiKeyHeader))
	ctx := context.WithValue(r.Context(), ApiKeyHeader, key)

	astronautID := mux.Vars(r)["astronautID"]
	id, err := strconv.Atoi(astronautID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Invalid Astronaut ID"})
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		writeJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Bad Request"})
		h.log.Warn("error deleting an astronaut", slog.Any("error", err))
		return
	}

	writeJSON(w, http.StatusOK, model.JSONResponse{Message: "Astronaut Deleted"})
}
