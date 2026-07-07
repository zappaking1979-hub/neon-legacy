package handlers

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/neonlegacy/server/internal/application"
	"github.com/neonlegacy/server/internal/domain/player"
	"github.com/neonlegacy/server/internal/middleware"
)

type AuthHandler struct {
	authService *application.AuthService
	tmpl        *template.Template
}

type registerForm struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
}

type loginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type apiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func NewAuthHandler(authService *application.AuthService, tmpl *template.Template) *AuthHandler {
	return &AuthHandler{authService: authService, tmpl: tmpl}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if middleware.GetPlayer(r) != nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	h.tmpl.ExecuteTemplate(w, "login.html", nil)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") == "application/json" {
		h.loginAPI(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	form := loginForm{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	session, _, err := h.authService.Login(r.Context(), form.Email, form.Password)
	if err != nil {
		h.tmpl.ExecuteTemplate(w, "login.html", map[string]string{"error": err.Error()})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    session.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(time.Until(session.ExpiresAt).Seconds()),
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *AuthHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if middleware.GetPlayer(r) != nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	h.tmpl.ExecuteTemplate(w, "register.html", nil)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") == "application/json" {
		h.registerAPI(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	form := registerForm{
		Email:    r.FormValue("email"),
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Gender:   r.FormValue("gender"),
	}

	p, err := h.authService.Register(r.Context(), form.Email, form.Username, form.Password, player.Gender(form.Gender))
	if err != nil {
		h.tmpl.ExecuteTemplate(w, "register.html", map[string]string{"error": err.Error()})
		return
	}

	session, _, err := h.authService.Login(r.Context(), form.Email, form.Password)
	if err == nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    session.Token,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   int(time.Until(session.ExpiresAt).Seconds()),
		})
	}

	http.Redirect(w, r, "/dashboard?welcome="+p.Username, http.StatusSeeOther)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie("session"); err == nil {
		h.authService.Logout(r.Context(), c.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *AuthHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	p := middleware.GetPlayer(r)
	if p == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	h.tmpl.ExecuteTemplate(w, "dashboard.html", map[string]interface{}{
		"Player": p,
		"welcome": r.URL.Query().Get("welcome"),
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	p := middleware.GetPlayer(r)
	if p == nil {
		writeJSON(w, http.StatusUnauthorized, apiResponse{Error: "not authenticated"})
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    p,
	})
}

func (h *AuthHandler) loginAPI(w http.ResponseWriter, r *http.Request) {
	var form loginForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Error: "invalid json"})
		return
	}

	session, _, err := h.authService.Login(r.Context(), form.Email, form.Password)
	if err != nil {
		code := http.StatusUnauthorized
		if errors.Is(err, application.ErrInvalidCreds) {
			code = http.StatusUnauthorized
		}
		writeJSON(w, code, apiResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Data:    map[string]string{"token": session.Token},
	})
}

func (h *AuthHandler) registerAPI(w http.ResponseWriter, r *http.Request) {
	var form registerForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Error: "invalid json"})
		return
	}

	p, err := h.authService.Register(r.Context(), form.Email, form.Username, form.Password, player.Gender(form.Gender))
	if err != nil {
		code := http.StatusBadRequest
		switch {
		case errors.Is(err, application.ErrEmailTaken),
			errors.Is(err, application.ErrUsernameTaken):
			code = http.StatusConflict
		}
		writeJSON(w, code, apiResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, apiResponse{
		Success: true,
		Data:    p,
	})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
