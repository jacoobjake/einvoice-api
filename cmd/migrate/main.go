package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jacoobjake/einvoice-api/config"
	"github.com/jacoobjake/einvoice-api/internal/database/migrations"
)

func main() {
	// Parse CLI flags (e.g. DSN and direction)
	cfg := config.Load()
	dbCfg := cfg.DBConfig
	connectionString := dbCfg.ConnectionString()
	fmt.Println(connectionString)
	direction := flag.String("direction", "up", "Migration direction: up, down or force")
	forceVersion := flag.Int("version", 0, "Force DB version")
	flag.Parse()

	switch *direction {
	case "up":
		if err := migrations.RunMigrations(connectionString); err != nil {
			log.Fatal(err)
		}
	case "down":
		if err := migrations.RollbackMigrations(connectionString); err != nil {
			log.Fatal(err)
		}
	case "force":
		if err := migrations.ForceDBVersion(connectionString, *forceVersion); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Unknown direction: %s", *direction)
	}
}
