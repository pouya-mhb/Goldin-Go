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
func WithRoutes(registerUser inbound.RegisterUser, loginUser inbound.LoginUser) func(chi.Router) {
	registrationHandler := NewRegistrationHandler(registerUser)
	loginHandler := NewLoginHandler(loginUser)

	return func(r chi.Router) {
		r.Route("/identity", func(r chi.Router) {
			r.Post("/register", registrationHandler.Register)
			r.Post("/login", loginHandler.Login)
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

// LoginHandler handles Identity login HTTP requests.
type LoginHandler struct {
	loginUser inbound.LoginUser
}

// NewLoginHandler constructs a LoginHandler.
func NewLoginHandler(loginUser inbound.LoginUser) *LoginHandler {
	return &LoginHandler{loginUser: loginUser}
}

// Login handles POST /identity/login.
func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request loginUserRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "request body must be valid JSON")
		return
	}

	loggedInUser, err := h.loginUser.LoginUser(r.Context(), command.LoginUser{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		writeLoginError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, loginUserResponseFromResult(loggedInUser))
}

type loginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginUserResponse struct {
	UserID                string `json:"user_id"`
	Email                 string `json:"email"`
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	TokenType             string `json:"token_type"`
	AccessTokenExpiresIn  int64  `json:"access_token_expires_in"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
}

func loginUserResponseFromResult(loggedInUser result.LoginUser) loginUserResponse {
	return loginUserResponse{
		UserID:                loggedInUser.UserID,
		Email:                 loggedInUser.Email,
		AccessToken:           loggedInUser.AccessToken,
		RefreshToken:          loggedInUser.RefreshToken,
		TokenType:             loggedInUser.TokenType,
		AccessTokenExpiresIn:  loggedInUser.AccessTokenExpiresIn,
		RefreshTokenExpiresIn: loggedInUser.RefreshTokenExpiresIn,
	}
}

func writeLoginError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		writeError(w, http.StatusUnauthorized, "invalid_credentials", "email or password is invalid")
	default:
		writeError(w, http.StatusInternalServerError, "internal_error", "login failed")
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
