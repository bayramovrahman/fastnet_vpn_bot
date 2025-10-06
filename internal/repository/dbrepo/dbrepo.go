package dbrepo

import (
	"database/sql"

	"github.com/bayramovrahman/fastnet_vpn_bot/internal/config"
	"github.com/bayramovrahman/fastnet_vpn_bot/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
