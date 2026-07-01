package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/adapters/persistence"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/command"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/ports/inbound"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/result"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/service"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
)

const maxRegistrationRequestBytes int64 = 1 << 20

// WithRoutes registers Identity HTTP routes.
func WithRoutes(registerUser inbound.RegisterUser) func(chi.Router) {
	handler := NewRegistrationHandler(registerUser)

	return func(r chi.Router) {
		r.Route("/identity", func(r chi.Router) {
			r.Post("/register", handler.Register)
		})
	}
}

// RegistrationHandler handles Identity registration HTTP requests.
type RegistrationHandler struct {
	registerUser inbound.RegisterUser
}

// NewRegistrationHandler constructs a RegistrationHandler.
func NewRegistrationHandler(registerUser inbound.RegisterUser) *RegistrationHandler {
	return &RegistrationHandler{registerUser: registerUser}
}

// Register handles POST /identity/register.
func (h *RegistrationHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request registerUserRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "request body must be valid JSON")
		return
	}

	registeredUser, err := h.registerUser.RegisterUser(r.Context(), command.RegisterUser{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		writeRegistrationError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, registerUserResponseFromResult(registeredUser))
}

type registerUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerUserResponse struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	RegisteredAt string `json:"registered_at"`
}

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func registerUserResponseFromResult(registeredUser result.RegisterUser) registerUserResponse {
	return registerUserResponse{
		UserID:       registeredUser.UserID,
		Email:        registeredUser.Email,
		RegisteredAt: registeredUser.RegisteredAt.Format("2006-01-02T15:04:05.999999Z07:00"),
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxRegistrationRequestBytes)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	if decoder.Decode(&struct{}{}) == nil {
		return errors.New("request body must contain a single JSON object")
	}

	return nil
}

func writeRegistrationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, valueobject.ErrEmptyEmail), errors.Is(err, valueobject.ErrInvalidEmail), errors.Is(err, service.ErrPasswordTooShort):
		writeError(w, http.StatusBadRequest, "invalid_registration", "email or password is invalid")
	case errors.Is(err, service.ErrEmailAlreadyRegistered), errors.Is(err, persistence.ErrDuplicateUserEmail):
		writeError(w, http.StatusConflict, "email_already_registered", "email is already registered")
	default:
		writeError(w, http.StatusInternalServerError, "internal_error", "registration failed")
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}

func writeError(w http.ResponseWriter, statusCode int, code string, message string) {
	writeJSON(w, statusCode, errorResponse{
		Error:   code,
		Message: message,
	})
}
