package handlers

import (
	"log"
	"net/http"

	"github.com/bayramovrahman/fastnet_vpn_bot/internal/config"
	"github.com/bayramovrahman/fastnet_vpn_bot/internal/driver"
	"github.com/bayramovrahman/fastnet_vpn_bot/internal/forms"
	"github.com/bayramovrahman/fastnet_vpn_bot/internal/models"
	"github.com/bayramovrahman/fastnet_vpn_bot/internal/render"
	"github.com/bayramovrahman/fastnet_vpn_bot/internal/repository"
	"github.com/bayramovrahman/fastnet_vpn_bot/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
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

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)

		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// Logout logs a user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
