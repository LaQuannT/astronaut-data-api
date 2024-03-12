package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/LaQuannT/astronaut-data-api/internal/model"
	"github.com/LaQuannT/astronaut-data-api/internal/transport/middleware"
	"github.com/LaQuannT/astronaut-data-api/internal/transport/util"
	"github.com/gorilla/mux"
)

type userHandler struct {
	service model.UserUsecase
	log     *slog.Logger
}

func RegisterUserHandlers(s model.UserUsecase, r *mux.Router, l *slog.Logger) {
	handler := userHandler{
		service: s,
		log:     l,
	}

	sr := r.PathPrefix("/users").Subrouter()
	sr.Use(middleware.APIKeyValidation(s, l))

	r.HandleFunc("/users", handler.CreateUser).Methods("POST")
	sr.HandleFunc("", handler.ListUsers).Methods("GET")
	sr.HandleFunc("/{userID}", handler.GetUser).Methods("GET")
	sr.HandleFunc("/{userID}", handler.UpdateUser).Methods("PUT")
	sr.HandleFunc("/{userID}", handler.DeleteUser).Methods("DELETE")
	sr.HandleFunc("/password/{userID}", handler.PasswordReset).Methods("PATCH")
	sr.HandleFunc("/apikey/{userID}", handler.APIKeyReset).Methods("PATCH")
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	u := new(model.User)
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Invalid request body"})
		h.log.Warn("error decoding json request body", slog.Any("error", err))
		return
	}

	u, errs := h.service.Create(ctx, u)
	if errs != nil {
		util.WriteJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Bad Request"})
		h.log.Warn("error creating new user", slog.Any("error", errs))
		return
	}

	util.WriteJSON(w, http.StatusCreated, model.JSONResponse{User: u})
}

func (h *userHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Invalid request query"})
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

	users, err := h.service.List(ctx, limit, offset)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Bad Request"})
		h.log.Warn("error listing users", slog.Any("error", err))
		return
	}

	util.WriteJSON(w, http.StatusOK, model.JSONResponse{Users: users})
}

func (h *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := mux.Vars(r)["userID"]

	id, err := strconv.Atoi(userID)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Invalid User ID"})
		return
	}

	u, err := h.service.Get(ctx, id)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Bad Request"})
		h.log.Warn("error fetching a user", slog.Any("error", err))
		return
	}

	util.WriteJSON(w, http.StatusOK, model.JSONResponse{User: u})
}

func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := new(model.User)

	userID := mux.Vars(r)["userID"]

	id, err := strconv.Atoi(userID)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Invalid User ID"})
		return
	}

	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Invalid request body"})
		h.log.Warn("error decoding json request body", slog.Any("error", err))
		return
	}

	u.ID = id

	u, errs := h.service.Update(ctx, u)
	if errs != nil {
		util.WriteJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Bad Request"})
		h.log.Warn("error updating a user", slog.Any("errors", errs))
		return
	}

	util.WriteJSON(w, http.StatusOK, model.JSONResponse{User: u})
}

func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := mux.Vars(r)["userID"]

	id, err := strconv.Atoi(userID)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Invalid User ID"})
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Bad Request"})
		h.log.Warn("error deleting a user", slog.Any("error", err))
		return
	}

	util.WriteJSON(w, http.StatusOK, model.JSONResponse{Message: "User Deleted"})
}

func (h *userHandler) PasswordReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := new(model.User)

	userID := mux.Vars(r)["userID"]

	id, err := strconv.Atoi(userID)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Invalid User ID"})
		return
	}

	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Invalid request body"})
		h.log.Warn("error decoding json request body", slog.Any("error", err))
		return
	}

	u.ID = id

	if errs := h.service.ResetPassword(ctx, u); err != nil {
		util.WriteJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Bad Request"})
		h.log.Warn("error reseting user password", slog.Any("errors", errs))
		return
	}

	util.WriteJSON(w, http.StatusOK, model.JSONResponse{Message: "User password reset"})
}

func (h *userHandler) APIKeyReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := mux.Vars(r)["userID"]

	id, err := strconv.Atoi(userID)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Invalid User ID"})
		return
	}

	u, err := h.service.GenerateNewAPIKey(ctx, id)
	if err != nil {
		util.WriteJSON(w, http.StatusBadRequest, model.JSONResponse{Error: "Bad Request"})
		h.log.Warn("error resetting API key", slog.Any("error", err))
		return
	}

	util.WriteJSON(w, http.StatusOK, model.JSONResponse{User: u})
}
