package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jacoobjake/einvoice-api/config"
	"github.com/jacoobjake/einvoice-api/internal/routes"
	"github.com/jacoobjake/einvoice-api/pkg/redisclient"
	_ "github.com/lib/pq"
	"github.com/stephenafamo/bob"
)

func main() {
	// Initialize Gin with default middleware (logger and recovery)
	r := gin.Default()

	// Initialize DB
	cfg := config.Load()
	db := initDB(cfg)

	defer db.Close()

	// Initialize redis
	rdb := redisclient.NewRedisClient(cfg.RedisConfig)

	// Pass db to routes if needed (example: api.RegisterRoutes(apiGroup, db))
	routes.RegisterRoutes(r, db, cfg, rdb)

	// Example: Register routes from other modules
	// invoice.RegisterRoutes(apiGroup, db)
	// user.RegisterRoutes(apiGroup, db)

	// Start server
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func initDB(cfg *config.Config) bob.DB {
	dbCfg := cfg.DBConfig
	db, err := bob.Open(dbCfg.Driver, dbCfg.ConnectionString())

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Optionally, you can set connection pool settings here
	// db.SetMaxOpenConns(25)
	// db.SetMaxIdleConns(25)
	// db.SetConnMaxLifetime(5 * time.Minute)

	return db
}
