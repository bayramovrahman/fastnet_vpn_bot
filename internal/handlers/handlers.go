package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/bayramovrahman/fastnet_vpn_bot/internal/config"
	"github.com/bayramovrahman/fastnet_vpn_bot/internal/driver"
	"github.com/bayramovrahman/fastnet_vpn_bot/internal/email"
	"github.com/bayramovrahman/fastnet_vpn_bot/internal/forms"
	"github.com/bayramovrahman/fastnet_vpn_bot/internal/helpers"
	"github.com/bayramovrahman/fastnet_vpn_bot/internal/models"
	"github.com/bayramovrahman/fastnet_vpn_bot/internal/render"
	"github.com/bayramovrahman/fastnet_vpn_bot/internal/repository"
	"github.com/bayramovrahman/fastnet_vpn_bot/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App          *config.AppConfig
	DB           repository.DatabaseRepo
	EmailService *email.EmailService
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App:          a,
		DB:           dbrepo.NewPostgresRepo(db.SQL, a),
		EmailService: email.NewEmailService(),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(repo *Repository) {
	Repo = repo
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Invoice(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "invoice.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Taxes(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "taxes.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Profile(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "profile.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	if helpers.IsAuthenticated(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Unable to parse form")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	emailForm := r.Form.Get("email")
	password := r.Form.Get("password")
	rememberMe := r.Form.Get("remember_me")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(emailForm, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid e-mail or password")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := m.DB.GetUserById(id)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Unable to retrieve user information")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Generate verification code
	code, err := email.GenerateVerificationCode()
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Unable to generate verification code")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Send verification email
	err = m.EmailService.SendVerificationCode(user.Email, code)
	if err != nil {
		log.Println("Email send error:", err)
		m.App.Session.Put(r.Context(), "error", "Unable to send verification code. Please try again.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Store verification data in session
	m.App.Session.Put(r.Context(), "pending_user_id", id)
	m.App.Session.Put(r.Context(), "pending_user_username", user.Username)
	m.App.Session.Put(r.Context(), "pending_user_first_name", user.FirstName)
	m.App.Session.Put(r.Context(), "pending_user_last_name", user.LastName)
	m.App.Session.Put(r.Context(), "pending_user_email", user.Email)
	m.App.Session.Put(r.Context(), "verification_code", code)
	m.App.Session.Put(r.Context(), "code_expires", time.Now().Add(10*time.Minute).Unix())

	if rememberMe == "on" {
		m.App.Session.Put(r.Context(), "pending_remember_me", true)
	}

	m.App.Session.Put(r.Context(), "warning", "Verification code sent to your e-mail")
	http.Redirect(w, r, "/verify", http.StatusSeeOther)
}

func (m *Repository) Verify(w http.ResponseWriter, r *http.Request) {
	// Check if there's a pending verification
	if !m.App.Session.Exists(r.Context(), "pending_user_id") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	render.Template(w, r, "verify.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostVerify(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Unable to parse form")
		http.Redirect(w, r, "/verify", http.StatusSeeOther)
		return
	}

	code := r.Form.Get("code")

	form := forms.New(r.PostForm)
	form.Required("code")
	form.MinLength("code", 6)

	if !form.Valid() {
		render.Template(w, r, "verify.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	// Check if verification session exists
	if !m.App.Session.Exists(r.Context(), "verification_code") {
		m.App.Session.Put(r.Context(), "error", "Verification session expired. Please login again.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Check if code expired
	expiresAt := m.App.Session.GetInt64(r.Context(), "code_expires")
	if time.Now().Unix() > expiresAt {
		// Clean up session
		m.App.Session.Remove(r.Context(), "pending_user_id")
		m.App.Session.Remove(r.Context(), "pending_user_username")
		m.App.Session.Remove(r.Context(), "pending_user_first_name")
		m.App.Session.Remove(r.Context(), "pending_user_last_name")
		m.App.Session.Remove(r.Context(), "pending_user_email")
		m.App.Session.Remove(r.Context(), "verification_code")
		m.App.Session.Remove(r.Context(), "code_expires")
		m.App.Session.Remove(r.Context(), "pending_remember_me")

		m.App.Session.Put(r.Context(), "error", "Verification code expired. Please login again.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Verify code
	storedCode := m.App.Session.GetString(r.Context(), "verification_code")
	if code != storedCode {
		m.App.Session.Put(r.Context(), "error", "Invalid verification code")
		http.Redirect(w, r, "/verify", http.StatusSeeOther)
		return
	}

	// Code is correct - complete login
	userId := m.App.Session.GetInt(r.Context(), "pending_user_id")
	username := m.App.Session.GetString(r.Context(), "pending_user_username")
	firstName := m.App.Session.GetString(r.Context(), "pending_user_first_name")
	lastName := m.App.Session.GetString(r.Context(), "pending_user_last_name")
	userEmail := m.App.Session.GetString(r.Context(), "pending_user_email")
	rememberMe := m.App.Session.GetBool(r.Context(), "pending_remember_me")

	// Set authenticated session
	m.App.Session.Put(r.Context(), "user_id", userId)
	m.App.Session.Put(r.Context(), "user_username", username)
	m.App.Session.Put(r.Context(), "user_first_name", firstName)
	m.App.Session.Put(r.Context(), "user_last_name", lastName)
	m.App.Session.Put(r.Context(), "user_email", userEmail)

	if rememberMe {
		m.App.Session.Put(r.Context(), "remember_me", true)
	}

	// Clean up verification session
	m.App.Session.Remove(r.Context(), "pending_user_id")
	m.App.Session.Remove(r.Context(), "pending_user_username")
	m.App.Session.Remove(r.Context(), "pending_user_first_name")
	m.App.Session.Remove(r.Context(), "pending_user_last_name")
	m.App.Session.Remove(r.Context(), "pending_user_email")
	m.App.Session.Remove(r.Context(), "verification_code")
	m.App.Session.Remove(r.Context(), "code_expires")
	m.App.Session.Remove(r.Context(), "pending_remember_me")

	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func (m *Repository) ResendCode(w http.ResponseWriter, r *http.Request) {
	// Check if there's a pending verification
	if !m.App.Session.Exists(r.Context(), "pending_user_id") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userId := m.App.Session.GetInt(r.Context(), "pending_user_id")
	user, err := m.DB.GetUserById(userId)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Unable to retrieve user information")
		http.Redirect(w, r, "/verify", http.StatusSeeOther)
		return
	}

	// Generate new verification code
	code, err := email.GenerateVerificationCode()
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Unable to generate verification code")
		http.Redirect(w, r, "/verify", http.StatusSeeOther)
		return
	}

	// Send verification email
	err = m.EmailService.SendVerificationCode(user.Email, code)
	if err != nil {
		log.Println("E-mail send error:", err)
		m.App.Session.Put(r.Context(), "error", "Unable to send verification code. Please try again.")
		http.Redirect(w, r, "/verify", http.StatusSeeOther)
		return
	}

	// Update session with new code
	m.App.Session.Put(r.Context(), "verification_code", code)
	m.App.Session.Put(r.Context(), "code_expires", time.Now().Add(10*time.Minute).Unix())

	m.App.Session.Put(r.Context(), "warning", "New verification code sent to your email")
	http.Redirect(w, r, "/verify", http.StatusSeeOther)
}

// Logout logs a user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
