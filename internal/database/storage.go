package database

import (
	"database/sql"
	"log"

	"github.com/akmal4410/gestapo/internal/config"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(database *config.Database) (*Storage, error) {
	log.Print(database.DBSource)
	db, err := sql.Open(database.DBDriver, database.DBSource)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	gormDB, err := gorm.Open(postgres.Open(database.DBSource), &gorm.Config{})
	if err != nil {
		log.Fatal("Error executing gorm  - ", err)
		return nil, err
	}

	AutoMigrateTables(gormDB)
	return &Storage{DB: db}, nil
}
