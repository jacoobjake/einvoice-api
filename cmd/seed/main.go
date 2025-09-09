package main

import (
	"log"

	"github.com/jacoobjake/einvoice-api/config"
	"github.com/jacoobjake/einvoice-api/internal/database/seeders"
	"github.com/stephenafamo/bob"
)

func main() {
	cfg := config.Load()
	dbCfg := cfg.DBConfig
	db, err := bob.Open(dbCfg.Driver, dbCfg.ConnectionString())

	if err != nil {
		panic(err)
	}

	defer db.Close()

	log.Println("Seeding users...")
	if err := seeders.SeedUsers(&db); err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}
	log.Println("Seeding completed.")
}
