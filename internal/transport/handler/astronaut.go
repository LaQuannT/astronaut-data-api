package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/LaQuannT/astronaut-data-api/internal/model"
	"github.com/gorilla/mux"
)

type astronautHandler struct {
	service model.AstronautUsecase
}

func RegisterAstronautHandlers(s model.AstronautUsecase, r *mux.Router) {
	handler := &astronautHandler{
		service: s,
	}

	r.HandleFunc("/", handler.CreateAstronaut).Methods("POST")
	r.HandleFunc("/?", handler.ListAstronauts).Methods("GET")
	r.HandleFunc("/{astronautID}", handler.GetAstronaut).Methods("GET")
	r.HandleFunc("/{astronautID}", handler.UpdateAstronaut).Methods("PUT")
	r.HandleFunc("/{astronautID}", handler.DeleteAstronaut).Methods("DELETE")
}

func (h *astronautHandler) CreateAstronaut(w http.ResponseWriter, r *http.Request) {
	a := new(model.Astronaut)

	key := r.Header.Get(string(auth))
	ctx := context.WithValue(r.Context(), auth, key)

	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		http.Error(w, "500 Internal Server error", http.StatusInternalServerError)
		log.Printf("error decoding request body to astronaut: %v", err)
		return
	}

	a, err := h.service.Create(ctx, a)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println(err)
		return
	}

	jsonBytes, err := json.Marshal(a)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Printf("error marshalling astronaut to json: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

func (h *astronautHandler) ListAstronauts(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get(string(auth))
	ctx := context.WithValue(r.Context(), auth, key)

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Printf("error parsing url request query: %v", err)
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
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println(err)
	}

	jsonBytes, err := json.Marshal(astronauts)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Printf("error marshalling list of astronauts to json: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

func (h *astronautHandler) GetAstronaut(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get(string(auth))
	ctx := context.WithValue(r.Context(), auth, key)

	astronautID := mux.Vars(r)["astronautID"]
	id, err := strconv.Atoi(astronautID)
	if err != nil {
		http.Error(w, "Bad Request - Invalid Astronaut ID", http.StatusBadRequest)
		return
	}

	a, err := h.service.Get(ctx, id)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println(err)
		return
	}

	jsonBytes, err := json.Marshal(a)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Printf("error marshalling astronaut to json: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

func (h *astronautHandler) UpdateAstronaut(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get(string(auth))
	ctx := context.WithValue(r.Context(), auth, key)
	a := new(model.Astronaut)

	astronautID := mux.Vars(r)["astronautID"]
	id, err := strconv.Atoi(astronautID)
	if err != nil {
		http.Error(w, "Bad Request - Invalid Astronaut ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		http.Error(w, "500 Internal Server error", http.StatusInternalServerError)
		log.Printf("error decoding request body to astronaut: %v", err)
		return
	}

	a.ID = id

	a, err = h.service.Update(ctx, a)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println(err)
		return
	}

	jsonBytes, err := json.Marshal(a)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Printf("error marshalling astronaut to json: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

func (h *astronautHandler) DeleteAstronaut(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get(string(auth))
	ctx := context.WithValue(r.Context(), auth, key)

	astronautID := mux.Vars(r)["astronautID"]
	id, err := strconv.Atoi(astronautID)
	if err != nil {
		http.Error(w, "Bad Request - Invalid Astronaut ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println(err)
		return
	}

	jsonBytes, err := json.Marshal(map[string]string{"message": "astronaut deleted"})
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Printf("error marshalling astronaut deleted message: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}
