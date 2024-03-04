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

type apiKeyHeader string

const auth apiKeyHeader = "X-api-key"

type userHandler struct {
	service model.UserUsecase
}

func RegisterUserHandlers(s model.UserUsecase, r *mux.Router) {
	handler := userHandler{
		service: s,
	}

	r.HandleFunc("/", handler.CreateUser).Methods("POST")
	r.HandleFunc("/", handler.ListUsers).Methods("GET")
	r.HandleFunc("/{userID}", handler.GetUser).Methods("GET")
	r.HandleFunc("/{userID}", handler.UpdateUser).Methods("PUT")
	r.HandleFunc("/{userID}", handler.DeleteUser).Methods("DELETE")
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	u := new(model.User)
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Printf("error decoding json request body: %v", err)
		return
	}

	u, err := h.service.Create(ctx, u)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Print(err)
	}

	jsonBytes, err := json.Marshal(u)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Printf("error marshalling user to json: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

func (h *userHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
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

	users, err := h.service.List(ctx, limit, offset)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println(err)
	}

	jsonBytes, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("error marshalling batch of users: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get(string(auth))
	ctx := context.WithValue(r.Context(), auth, key)

	userID := mux.Vars(r)["userID"]

	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Bad Request - Invaild User ID", http.StatusBadRequest)
		return
	}

	u, err := h.service.Get(ctx, id)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println(err)
		return
	}

	jsonBytes, err := json.Marshal(u)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("error marshalling user: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get(string(auth))
	ctx := context.WithValue(r.Context(), auth, key)
	u := new(model.User)

	userID := mux.Vars(r)["userID"]

	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Bad Request - Invaild User ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Printf("error decoding json request body: %v", err)
		return
	}

	u.ID = id

	u, err = h.service.Update(ctx, u)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println(err)
		return
	}

	jsonBytes, err := json.Marshal(u)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("error marshalling user: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get(string(auth))
	ctx := context.WithValue(r.Context(), auth, key)

	userID := mux.Vars(r)["userID"]

	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Bad Request - Invaild User ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println(err)
	}

	jsonBytes, err := json.Marshal(map[string]string{"message": "User deleted"})
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("error marshalling user deleted message: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
