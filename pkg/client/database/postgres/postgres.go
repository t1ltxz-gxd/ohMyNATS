package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"ohMyNATS/pkg/client"
)

type Postgres struct {
	client.Builder
	Host     string
	Port     int
	User     string
	Dbname   string
	Password string
}

func (pg *Postgres) New() (*gorm.DB, error) {
	args := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		pg.Host, pg.Port, pg.User, pg.Dbname, pg.Password)
	db, err := gorm.Open(
		"postgres",
		args)
	return db, err
}
