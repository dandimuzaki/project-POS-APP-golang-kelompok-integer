package main

import (
	"log"
	"project-POS-APP-golang-integer/cmd"
	"project-POS-APP-golang-integer/internal/data"
	"project-POS-APP-golang-integer/internal/data/repository"
	"project-POS-APP-golang-integer/internal/wire"
	"project-POS-APP-golang-integer/pkg/database"
	"project-POS-APP-golang-integer/pkg/utils"
)

func main() {
	config, err := utils.ReadConfiguration()
	if err != nil {
		log.Fatalf("failed to read file config: %v", err)
	}

	db, err := database.InitDB(config.DB)
	if err != nil {
		log.Fatalf("failed to connect to postgres database: %v", err)
	}

	logger, err := utils.InitLogger(config.PathLogging, config.Debug)

	// migration
	err = data.AutoMigrate(db)
	if err != nil {
		log.Println(err)
	}

	// seeder
	err = data.SeedAll(db)
	if err != nil {
		log.Println(err)
	}

	repo := repository.NewRepository(db, logger)
	route := wire.Wiring(repo, logger, config)
	cmd.APiserver(route)
}
