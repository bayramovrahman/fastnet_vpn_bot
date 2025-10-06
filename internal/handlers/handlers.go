package handlers

import (
	"net/http"

	"github.com/bayramovrahman/fastnet_vpn_bot/internal/config"
	"github.com/bayramovrahman/fastnet_vpn_bot/internal/driver"
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
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Invoice(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "invoice.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Taxes(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "taxes.page.tmpl", &models.TemplateData{})
}

func (m *Repository) AdminLogin(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "login.page.tmpl", &models.TemplateData{})
}
